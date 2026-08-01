package documents

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/observability"
	"github.com/google/uuid"
)

var (
	ErrDocumentNotFound = errors.New("document not found")
	ErrUploadRejected   = errors.New("upload rejected")
)

// DocumentRecord represents a stored supporting document.
type DocumentRecord struct {
	ID            string
	EntityType    string
	EntityID      string
	DocumentName  string
	StorageKey    string
	ContentType   string
	FileSizeBytes int64
	UploadedBy    *domain.UserID
	CreatedAt     time.Time
}

// DocumentService handles upload validation, storage, retrieval, and
// authorization for supporting documents.
type DocumentService struct {
	storage Storage
	records map[string]*DocumentRecord // in-memory metadata (pgx adapter later)
	audit   func(ctx context.Context, e *domain.AuditEvent) error
	log     *slog.Logger
}

// NewDocumentService creates a document service with the given storage backend.
func NewDocumentService(storage Storage, auditFunc func(ctx context.Context, e *domain.AuditEvent) error, log *slog.Logger) *DocumentService {
	return &DocumentService{
		storage: storage,
		records: make(map[string]*DocumentRecord),
		audit:   auditFunc,
		log:     log,
	}
}

// Upload validates, stores, and records a supporting document. The caller
// must have already verified the user is authenticated and authorized to
// access the parent entity.
func (s *DocumentService) Upload(ctx context.Context, actor domain.UserID, entityType, entityID string, header *multipartFileHeader, content []byte) (*DocumentRecord, error) {
	corrID := observability.CorrelationFrom(ctx)

	vf, err := ValidateUpload(header, content, entityType, entityID)
	if err != nil {
		s.log.InfoContext(ctx, "upload rejected", "entity_type", entityType, "entity_id", entityID, "error", err)
		return nil, fmt.Errorf("%w: %v", ErrUploadRejected, err)
	}

	// Store the file.
	if err := s.storage.Store(ctx, vf.StorageKey, vf.Content, vf.ContentType); err != nil {
		return nil, fmt.Errorf("store file: %w", err)
	}

	// Record metadata.
	record := &DocumentRecord{
		ID:            uuid.NewString(),
		EntityType:    entityType,
		EntityID:      entityID,
		DocumentName:  vf.DocumentName,
		StorageKey:    vf.StorageKey,
		ContentType:   vf.ContentType,
		FileSizeBytes: vf.FileSize,
		UploadedBy:    &actor,
		CreatedAt:     time.Now().UTC(),
	}
	s.records[record.ID] = record

	// Audit.
	_ = s.audit(ctx, &domain.AuditEvent{
		CorrelationID: corrID,
		ActorUserID:   &actor,
		ActorRole:     "administrator",
		Action:        "document_uploaded",
		EntityType:    "supporting_document",
		EntityID:      record.ID,
		Details: map[string]any{
			"file_name":    vf.DocumentName,
			"content_type": vf.ContentType,
			"file_size":    vf.FileSize,
			"parent_type":  entityType,
			"parent_id":    entityID,
		},
	})

	return record, nil
}

// Download retrieves a document by ID. The caller must verify that the
// authenticated user is authorized to access the parent entity — this method
// does not perform authorization checks itself. The route handler handles
// that via RequireAuth and entity ownership verification.
func (s *DocumentService) Download(ctx context.Context, documentID string) (content []byte, contentType string, documentName string, err error) {
	record, ok := s.records[documentID]
	if !ok {
		return nil, "", "", ErrDocumentNotFound
	}
	content, contentType, err = s.storage.Retrieve(ctx, record.StorageKey)
	if err != nil {
		return nil, "", "", fmt.Errorf("retrieve: %w", err)
	}
	return content, contentType, record.DocumentName, nil
}

// ListByEntity returns documents for a given entity.
func (s *DocumentService) ListByEntity(ctx context.Context, entityType, entityID string) []*DocumentRecord {
	var out []*DocumentRecord
	for _, r := range s.records {
		if r.EntityType == entityType && r.EntityID == entityID {
			out = append(out, r)
		}
	}
	return out
}

// Delete removes a document and its stored file.
func (s *DocumentService) Delete(ctx context.Context, documentID string, actor domain.UserID) error {
	record, ok := s.records[documentID]
	if !ok {
		return ErrDocumentNotFound
	}
	if err := s.storage.Delete(ctx, record.StorageKey); err != nil {
		return fmt.Errorf("delete file: %w", err)
	}
	delete(s.records, documentID)
	_ = s.audit(ctx, &domain.AuditEvent{
		CorrelationID: observability.CorrelationFrom(ctx),
		ActorUserID:   &actor,
		ActorRole:     "administrator",
		Action:        "document_deleted",
		EntityType:    "supporting_document",
		EntityID:      documentID,
		Details:       map[string]any{"file_name": record.DocumentName},
	})
	return nil
}

// ServeFile writes document content to an HTTP response with appropriate
// headers. Content is always served as an attachment with nosniff and
// no-store to prevent caching and MIME-type confusion.
func ServeFile(w http.ResponseWriter, content []byte, contentType, fileName string) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "private, no-store")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}