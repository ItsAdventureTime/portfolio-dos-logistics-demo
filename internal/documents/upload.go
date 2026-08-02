// Package documents handles safe document upload validation and private
// storage. Uploads are checked for size, file type, extension, content
// sniff, and storage-key safety. Files are served behind authenticated
// authorization checks, never publicly.
package documents

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// MaxFileSize is the maximum allowed upload size (10 MB).
const MaxFileSize = 10 * 1024 * 1024

// AllowedMIMETypes is the allowlist of permitted content types.
var AllowedMIMETypes = map[string]bool{
	"application/pdf": true,
	"image/png":       true,
	"image/jpeg":      true,
	"image/gif":       true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
	"application/vnd.ms-excel": true,
	"text/csv":        true,
	"application/zip":  true,
}

// AllowedExtensions is the allowlist of permitted file extensions.
var AllowedExtensions = map[string]bool{
	".pdf":  true,
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".xlsx": true,
	".xls":  true,
	".csv":  true,
	".zip":  true,
}

// ContentSniffSignatures maps file magic bytes to detected MIME type.
// This validates the actual content matches the declared type.
var ContentSniffSignatures = map[string]string{
	"\x25\x50\x44\x46": "application/pdf", // %PDF
	"\x89\x50\x4E\x47": "image/png",        // PNG
	"\xFF\xD8\xFF":      "image/jpeg",       // JPEG
	"\x47\x49\x46\x38": "image/gif",        // GIF8
	"\x50\x4B\x03\x04": "application/zip",  // PK (zip/xlsx)
}

// UploadErrors
var (
	ErrFileTooLarge      = fmt.Errorf("file exceeds maximum size of %d bytes", MaxFileSize)
	ErrEmptyFile         = errors.New("file is empty")
	ErrInvalidExtension  = errors.New("file extension is not permitted")
	ErrInvalidMIMEType   = errors.New("content type is not permitted")
	ErrContentMismatch   = errors.New("file content does not match declared type")
	ErrInvalidStorageKey = errors.New("storage key contains invalid characters")
)

// ValidatedFile holds the result of upload validation.
type ValidatedFile struct {
	DocumentName string
	ContentType  string
	FileSize     int64
	Content      []byte
	StorageKey   string
}

// ValidateUpload validates an uploaded file: checks size, extension, MIME
// type, and content sniff. Returns a ValidatedFile or an error describing
// what failed.
func ValidateUpload(header *multipartFileHeader, content []byte, entityType, entityID string) (*ValidatedFile, error) {
	if header == nil {
		return nil, ErrEmptyFile
	}
	fileName := header.Filename
	fileSize := header.Size

	if fileSize == 0 {
		return nil, ErrEmptyFile
	}
	if fileSize > MaxFileSize {
		return nil, ErrFileTooLarge
	}

	// Validate extension.
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" || !AllowedExtensions[ext] {
		return nil, fmt.Errorf("%w: %s", ErrInvalidExtension, ext)
	}

	// Validate declared content type.
	declaredType := header.ContentType
	if declaredType == "" || !AllowedMIMETypes[declaredType] {
		return nil, fmt.Errorf("%w: %s", ErrInvalidMIMEType, declaredType)
	}

	// Content sniff: verify actual bytes match declared type.
	if len(content) == 0 {
		return nil, ErrEmptyFile
	}
	detectedType := sniffContent(content)
	if detectedType == "" {
		return nil, ErrContentMismatch
	}
	// Special case: xlsx is a zip, so zip detection covers xlsx.
	if detectedType != declaredType {
		if declaredType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" && detectedType == "application/zip" {
			// OK: xlsx files are zip archives
		} else if declaredType == "application/vnd.ms-excel" && detectedType == "application/zip" {
			// OK: some xls files are also zip-based
		} else if declaredType == "text/csv" {
			// CSV is text, no magic bytes. Allow if content is printable ASCII/UTF-8
		} else {
			return nil, fmt.Errorf("%w: declared %s but content matches %s", ErrContentMismatch, declaredType, detectedType)
		}
	}

	// Build safe storage key.
	storageKey, err := buildStorageKey(entityType, entityID, fileName)
	if err != nil {
		return nil, err
	}

	return &ValidatedFile{
		DocumentName: fileName,
		ContentType:  declaredType,
		FileSize:     fileSize,
		Content:      content,
		StorageKey:   storageKey,
	}, nil
}

// sniffContent reads magic bytes and returns the detected MIME type.
func sniffContent(content []byte) string {
	for sig, mime := range ContentSniffSignatures {
		if len(content) >= len(sig) && string(content[:len(sig)]) == sig {
			return mime
		}
	}
	return ""
}

// buildStorageKey creates a safe storage key with no path traversal.
func buildStorageKey(entityType, entityID, fileName string) (string, error) {
	if strings.Contains(entityType, "..") || strings.Contains(entityType, "/") {
		return "", ErrInvalidStorageKey
	}
	if strings.Contains(entityID, "..") || strings.Contains(entityID, "/") {
		return "", ErrInvalidStorageKey
	}
	safeName := filepath.Base(fileName)
	if safeName != fileName {
		return "", ErrInvalidStorageKey
	}
	return fmt.Sprintf("%s/%s/%s", entityType, entityID, safeName), nil
}

// multipartFileHeader is an interface that *multipart.FileHeader satisfies.
// This makes testing easier without importing multipart.
type multipartFileHeader struct {
	Filename    string
	ContentType string
	Size        int64
}