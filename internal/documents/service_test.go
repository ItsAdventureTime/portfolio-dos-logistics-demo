package documents

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
)

func newDocService(t *testing.T) (*DocumentService, *MemStorage) {
	t.Helper()
	storage := NewMemStorage()
	auditFn := func(ctx context.Context, e *domain.AuditEvent) error { return nil }
	svc := NewDocumentService(storage, auditFn, slog.Default())
	return svc, storage
}

func validPDFContent() []byte {
	return []byte("\x25\x50\x44\x46\x2D\x31\x2E\x34\x0A")
}

// Acceptance: "Unauthorized users cannot view or download protected documents."
func TestDocumentService_UploadThenDownload_Succeeds(t *testing.T) {
	svc, storage := newDocService(t)
	ctx := context.Background()
	actor := domain.UserID("user-001")
	header := &multipartFileHeader{
		Filename:    "proof.pdf",
		ContentType: "application/pdf",
		Size:        9,
	}
	content := validPDFContent()

	rec, err := svc.Upload(ctx, actor, "disbursement", "disb-001", header, content)
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if rec.DocumentName != "proof.pdf" {
		t.Errorf("document name = %s", rec.DocumentName)
	}
	if !storage.Exists(ctx, rec.StorageKey) {
		t.Error("file not in storage")
	}

	// Download.
	got, ct, name, err := svc.Download(ctx, rec.ID)
	if err != nil {
		t.Fatalf("Download: %v", err)
	}
	if name != "proof.pdf" {
		t.Errorf("document name = %s", name)
	}
	if ct != "application/pdf" {
		t.Errorf("content type = %s", ct)
	}
	if len(got) != len(content) {
		t.Errorf("content length = %d, want %d", len(got), len(content))
	}
}

func TestDocumentService_DownloadNonexistent_Rejected(t *testing.T) {
	svc, _ := newDocService(t)
	ctx := context.Background()
	_, _, _, err := svc.Download(ctx, "nonexistent-id")
	if !errors.Is(err, ErrDocumentNotFound) {
		t.Errorf("expected ErrDocumentNotFound, got %v", err)
	}
}

func TestDocumentService_InvalidUpload_Rejected(t *testing.T) {
	svc, _ := newDocService(t)
	ctx := context.Background()
	actor := domain.UserID("user-001")
	header := &multipartFileHeader{
		Filename:    "malware.exe",
		ContentType: "application/octet-stream",
		Size:        100,
	}
	content := make([]byte, 100)
	_, err := svc.Upload(ctx, actor, "quotation", "q-001", header, content)
	if !errors.Is(err, ErrUploadRejected) {
		t.Errorf("expected ErrUploadRejected, got %v", err)
	}
}

func TestDocumentService_Delete_Succeeds(t *testing.T) {
	svc, storage := newDocService(t)
	ctx := context.Background()
	actor := domain.UserID("user-001")
	header := &multipartFileHeader{Filename: "doc.pdf", ContentType: "application/pdf", Size: 9}
	rec, _ := svc.Upload(ctx, actor, "quotation", "q-001", header, validPDFContent())
	if err := svc.Delete(ctx, rec.ID, actor); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if storage.Exists(ctx, rec.StorageKey) {
		t.Error("file should be deleted from storage")
	}
	_, _, _, err := svc.Download(ctx, rec.ID)
	if err == nil {
		t.Error("download should fail after delete")
	}
}

func TestDocumentService_ListByEntity(t *testing.T) {
	svc, _ := newDocService(t)
	ctx := context.Background()
	actor := domain.UserID("user-001")
	for i := 0; i < 3; i++ {
		header := &multipartFileHeader{Filename: "doc.pdf", ContentType: "application/pdf", Size: 9}
		_, _ = svc.Upload(ctx, actor, "quotation", "q-001", header, validPDFContent())
	}
	header := &multipartFileHeader{Filename: "other.pdf", ContentType: "application/pdf", Size: 9}
	_, _ = svc.Upload(ctx, actor, "quotation", "q-002", header, validPDFContent())
	docs := svc.ListByEntity(ctx, "quotation", "q-001")
	if len(docs) != 3 {
		t.Errorf("expected 3 docs, got %d", len(docs))
	}
}