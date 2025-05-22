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

const (
	baseURLEnv       = "BASE_URL"
	uploadsPathEnv   = "UPLOADS_DIR"
	defaultUploads   = "internal/uploads"
	subfolderUploads = "uploadsArchives"
	prodEnv          = "PROD"
)

var (
	baseURL        string
	baseUploadPath string
)

func init() {
	// Leer URL base de la API
	baseURL = os.Getenv(baseURLEnv)
	prod := os.Getenv(prodEnv)
	if prod == "false" {
		baseURL = "http://localhost:3000"
	} else {
		baseURL = "https://api-4618.onrender.com"
	}

	// Leer path de subida de archivos
	baseUploadPath = os.Getenv(uploadsPathEnv)
	if baseUploadPath == "" {
		baseUploadPath = defaultUploads
	}
}

// Archive representa un archivo subido
// Guarda solo el nombre de archivo, la ruta física y la URL se calculan dinámicamente

type Archive struct {
	Archive string
}

// NewArchive procesa un string Base64, valida el mime y guarda el archivo
// retorna un Archive con el nombre de archivo generado
func NewArchive(archive string) (*Archive, error) {
	if archive == "" || !strings.Contains(archive, ",") {
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

	// Generar un nombre único
	uniqueID := uuid.New().String()
	fileName := fmt.Sprintf("%s%s", uniqueID, ext)

	// Construir rutas según variables de entorno
	uploadPath := filepath.Join(baseUploadPath, subfolderUploads)
	filePath := filepath.Join(uploadPath, fileName)

	// Crear directorios si no existen
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("error al crear directorio de archivos: %w", err)
	}

	// Escribir el archivo en disco persistente
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, fmt.Errorf("error al guardar el archivo: %w", err)
	}

	return &Archive{Archive: fileName}, nil
}

// URL devuelve la ruta pública completa donde se sirve el archivo
func (a *Archive) URL() string {
	if a == nil || a.Archive == "" {
		return ""
	}
	// Ej: https://mi-api.com/uploads/uploadsArchives/<fileName>
	return fmt.Sprintf("%s/uploads/%s/%s", baseURL, subfolderUploads, filepath.Base(a.Archive))
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

// Extensiones permitidas para archivos de propuesta
var allowedExtensions = map[string]string{
	"application/pdf":    ".pdf",
	"application/msword": ".doc",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
	"application/vnd.ms-excel": ".xls",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
}
