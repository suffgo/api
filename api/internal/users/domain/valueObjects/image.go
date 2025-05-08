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
	// Puedes también leer esto de env si lo necesitas
	baseURLEnv       = "BASE_URL"
	uploadsPathEnv   = "UPLOADS_DIR"
	defaultUploads   = "internal/uploads"
	subfolderUploads = "uploadsUsers"
)

var (
	// path absoluto donde montar tu disco persistente
	baseUploadPath string
	// URL base de tu API (ej: https://mi-app.onrender.com)
	baseURL string
)

func init() {
	// 1) Path en disco
	baseUploadPath = os.Getenv(uploadsPathEnv)
	if baseUploadPath == "" {
		baseUploadPath = defaultUploads
	}

	// 2) URL pública
	baseURL = os.Getenv(baseURLEnv)
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
}

// Image representa el VO de una imagen subida
type Image struct {
	Image string
}

// NewImage procesa el Base64, genera un UUID y guarda el archivo bajo:
//   <baseUploadPath>/<subfolderUploads>/<uuid>.<ext>
func NewImage(image string) (*Image, error) {
	if image == "" || !strings.Contains(image, ",") {
		return &Image{Image: image}, nil
	}

	mimeType, data, err := decodeBase64Image(image)
	if err != nil {
		return nil, errors.New("error al procesar la imagen")
	}

	// valida formatos
	switch mimeType {
	case "image/png", "image/jpg", "image/jpeg", "image/webp":
	default:
		return nil, errors.New("formato de imagen no soportado")
	}

	// genera nombre único
	uniqueID := uuid.New().String()
	ext := ".png"
	if mimeType == "image/jpg" {
		ext = ".jpg"
	} else if mimeType == "image/webp" {
		ext = ".webp"
	}
	fileName := uniqueID + ext

	// construye rutas según env
	uploadPath := filepath.Join(baseUploadPath, subfolderUploads)
	filePath := filepath.Join(uploadPath, fileName)

	// crea jerarquía si hace falta
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("error al crear directorio de imágenes: %w", err)
	}

	// escribe el archivo
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, fmt.Errorf("error al guardar la imagen: %w", err)
	}

	return &Image{Image: fileName}, nil
}

// URL devuelve la ruta pública donde estará servida la imagen
func (i *Image) URL() string {
	if i == nil || i.Image == "" {
		return ""
	}
	// ejemplo: https://mi-app.onrender.com/uploads/uploadsUsers/<fileName>
	return fmt.Sprintf("%s/uploads/%s/%s", baseURL, subfolderUploads, filepath.Base(i.Image))
}

func decodeBase64Image(base64Image string) (string, []byte, error) {
	parts := strings.Split(base64Image, ",")
	if len(parts) != 2 {
		return "", nil, errors.New("formato de imagen Base64 inválido")
	}

	mimeType := strings.TrimPrefix(strings.Split(parts[0], ";")[0], "data:")
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", nil, errors.New("error al decodificar Base64")
	}

	return mimeType, data, nil
}
