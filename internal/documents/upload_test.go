package documents

import (
	"bytes"
	"testing"
)

func makeHeader(name, contentType string, size int64) *multipartFileHeader {
	return &multipartFileHeader{
		Filename:    name,
		ContentType: contentType,
		Size:        size,
	}
}

// Acceptance: "Upload validation rejects unsafe or oversized content."
func TestValidateUpload_Oversized_Rejected(t *testing.T) {
	content := make([]byte, MaxFileSize+1)
	header := makeHeader("big.pdf", "application/pdf", int64(len(content)))
	_, err := ValidateUpload(header, content, "quotation", "123")
	if err == nil {
		t.Error("expected error for oversized file")
	}
}

func TestValidateUpload_Empty_Rejected(t *testing.T) {
	header := makeHeader("empty.pdf", "application/pdf", 0)
	_, err := ValidateUpload(header, nil, "quotation", "123")
	if err == nil {
		t.Error("expected error for empty file")
	}
}

func TestValidateUpload_InvalidExtension_Rejected(t *testing.T) {
	pdfContent := []byte("\x25\x50\x44\x46\x2D\x31\x2E\x34")
	header := makeHeader("malware.exe", "application/octet-stream", int64(len(pdfContent)))
	_, err := ValidateUpload(header, pdfContent, "quotation", "123")
	if err == nil {
		t.Error("expected error for .exe extension")
	}
}

func TestValidateUpload_InvalidMIME_Rejected(t *testing.T) {
	pdfContent := []byte("\x25\x50\x44\x46\x2D\x31\x2E\x34")
	header := makeHeader("doc.pdf", "application/octet-stream", int64(len(pdfContent)))
	_, err := ValidateUpload(header, pdfContent, "quotation", "123")
	if err == nil {
		t.Error("expected error for invalid MIME type")
	}
}

func TestValidateUpload_ContentMismatch_Rejected(t *testing.T) {
	// PNG content but declared as PDF
	pngContent := []byte("\x89\x50\x4E\x47\x0D\x0A\x1A\x0A")
	header := makeHeader("fake.pdf", "application/pdf", int64(len(pngContent)))
	_, err := ValidateUpload(header, pngContent, "quotation", "123")
	if err == nil {
		t.Error("expected content mismatch error (PNG declared as PDF)")
	}
}

func TestValidateUpload_PathTraversal_Rejected(t *testing.T) {
	pdfContent := []byte("\x25\x50\x44\x46\x2D\x31\x2E\x34")
	header := makeHeader("../../../etc/passwd.pdf", "application/pdf", int64(len(pdfContent)))
	_, err := ValidateUpload(header, pdfContent, "quotation", "123")
	if err == nil {
		t.Error("expected error for path traversal in filename")
	}
}

func TestValidateUpload_EntityPathTraversal_Rejected(t *testing.T) {
	pdfContent := []byte("\x25\x50\x44\x46\x2D\x31\x2E\x34")
	header := makeHeader("doc.pdf", "application/pdf", int64(len(pdfContent)))
	_, err := ValidateUpload(header, pdfContent, "quotation", "../../../etc")
	if err == nil {
		t.Error("expected error for path traversal in entity ID")
	}
}

func TestValidateUpload_ValidPDF_Succeeds(t *testing.T) {
	pdfContent := []byte("\x25\x50\x44\x46\x2D\x31\x2E\x34\x0A")
	header := makeHeader("invoice.pdf", "application/pdf", int64(len(pdfContent)))
	vf, err := ValidateUpload(header, pdfContent, "quotation", "quot-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vf.DocumentName != "invoice.pdf" {
		t.Errorf("document name = %s", vf.DocumentName)
	}
	if vf.ContentType != "application/pdf" {
		t.Errorf("content type = %s", vf.ContentType)
	}
	if vf.StorageKey != "quotation/quot-001/invoice.pdf" {
		t.Errorf("storage key = %s, want quotation/quot-001/invoice.pdf", vf.StorageKey)
	}
}

func TestValidateUpload_ValidPNG_Succeeds(t *testing.T) {
	pngContent := []byte("\x89\x50\x4E\x47\x0D\x0A\x1A\x0A\x00\x00")
	header := makeHeader("photo.png", "image/png", int64(len(pngContent)))
	vf, err := ValidateUpload(header, pngContent, "disbursement", "disb-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vf.StorageKey != "disbursement/disb-001/photo.png" {
		t.Errorf("storage key = %s", vf.StorageKey)
	}
}

// Acceptance: storage round-trip
func TestMemStorage_StoreRetrieve(t *testing.T) {
	s := NewMemStorage()
	ctx := t.Context()
	key := "quotation/quot-001/invoice.pdf"
	content := []byte("\x25\x50\x44\x46")
	if err := s.Store(ctx, key, content, "application/pdf"); err != nil {
		t.Fatalf("Store: %v", err)
	}
	if !s.Exists(ctx, key) {
		t.Error("file should exist after store")
	}
	got, ct, err := s.Retrieve(ctx, key)
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Error("content mismatch")
	}
	if ct != "application/pdf" {
		t.Errorf("content type = %s", ct)
	}
	if err := s.Delete(ctx, key); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if s.Exists(ctx, key) {
		t.Error("file should not exist after delete")
	}
}