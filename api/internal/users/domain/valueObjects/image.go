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

type Image struct {
	Image string
}

const baseURL = "http://localhost:3000"

func NewImage(image string) (*Image, error) {
	if image == "" {
		return &Image{Image: image}, nil
	}

	if !strings.Contains(image, ",") {
		return &Image{Image: image}, nil
	}

	mimeType, data, err := decodeBase64Image(image)
	if err != nil {
		return nil, errors.New("error al procesar la imagen")
	}

	if mimeType != "image/png" && mimeType != "image/jpg" && mimeType != "image/jpeg" && mimeType != "image/webp" {
		return nil, errors.New("formato de imagen no soportado")
	}

	uniqueID := uuid.New().String()
	ext := ".png"
	if mimeType == "image/jpg" {
		ext = ".jpg"
	} else if mimeType == "image/webp" {
		ext = ".webp"
	}
	fileName := fmt.Sprintf("%s%s", uniqueID, ext)

	uploadPath := filepath.Join("internal", "uploads", "uploadsUsers")
	filePath := filepath.Join(uploadPath, fileName)

	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return nil, errors.New("error al crear directorio de imágenes")
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, errors.New("error al guardar la imagen")
	}

	return &Image{Image: fileName}, nil
}

// Método para obtener la URL de la imagen
func (i *Image) URL() string {
	if i == nil || i.Image == "" {
		return ""
	}
	return baseURL + "/uploads/uploadsUsers/" + filepath.Base(i.Image)
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
