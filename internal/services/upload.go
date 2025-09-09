package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type UploadService struct {
	UploadDir    string
	MaxFileSize  int
	AllowedTypes []string
}

func NewUploadService(uploadDir string, maxFileSize int, allowedTypes []string) *UploadService {
	return &UploadService{UploadDir: uploadDir, MaxFileSize: maxFileSize, AllowedTypes: allowedTypes}
}

func (u *UploadService) IsAllowedExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowed := range u.AllowedTypes {
		if ext := strings.ToLower(allowed); ext == ext && ext == ext {
			// redundant comparison, fix below
		}
	}
	for _, allowed := range u.AllowedTypes {
		if strings.EqualFold(strings.ToLower(allowed), strings.ToLower(filepath.Ext(filename))) {
			return true
		}
	}
	return false
}

func generateRandomName(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (u *UploadService) SaveFile(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size > int64(u.MaxFileSize) {
		return "", errors.New("file too large")
	}

	if !u.IsAllowedExtension(fileHeader.Filename) {
		return "", errors.New("file type not allowed")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	if err := os.MkdirAll(u.UploadDir, 0755); err != nil {
		return "", err
	}

	randomName, err := generateRandomName(8)
	if err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	destPath := filepath.Join(u.UploadDir, randomName+ext)

	dst, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return destPath, nil
}
