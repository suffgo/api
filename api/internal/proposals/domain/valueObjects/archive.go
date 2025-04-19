package valueobjects

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type Archive struct {
	Archive string
}

const fileBaseURL = "http://localhost:3000"

// Extensiones permitidas
var allowedExtensions = map[string]string{
	"application/pdf":    ".pdf",
	"application/msword": ".doc",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
	"application/vnd.ms-excel": ".xls",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
}

func NewArchive(archive string) (*Archive, error) {
	if archive == "" {
		return &Archive{Archive: archive}, nil
	}

	if !strings.Contains(archive, ",") {
		return &Archive{Archive: archive}, nil
	}

	mimeType, data, err := decodeBase64File(archive)
	if err != nil {
		return nil, errors.New("error al procesar el archivo")
	}

	ext, valid := allowedExtensions[mimeType]
	if !valid {
		return nil, errors.New("formato de archivo no soportado")
	}

	uniqueID := uuid.New().String()
	fileName := fmt.Sprintf("%s%s", uniqueID, ext)

	uploadPath := filepath.Join("internal", "uploads", "proposalFiles")
	filePath := filepath.Join(uploadPath, fileName)

	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return nil, errors.New("error al crear directorio de archivos")
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, errors.New("error al guardar el archivo")
	}

	return &Archive{Archive: fileName}, nil
}

// Método para obtener la URL del archivo
func (a *Archive) URL() string {
	if a == nil || a.Archive == "" {
		return ""
	}
	return fileBaseURL + "/uploads/proposalFiles" + filepath.Base(a.Archive)
}

func decodeBase64File(base64File string) (string, []byte, error) {
	parts := strings.Split(base64File, ",")
	if len(parts) != 2 {
		return "", nil, errors.New("formato de archivo Base64 inválido")
	}

	mimeType := strings.TrimPrefix(strings.Split(parts[0], ";")[0], "data:")
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", nil, errors.New("error al decodificar Base64")
	}

	return mimeType, data, nil
}
