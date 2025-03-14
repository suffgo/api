package valueobjects

import (
	"path/filepath"
)

type Image struct {
	Image string
}

func NewImage(image string) (*Image, error) {
	return &Image{
		Image: image,
	}, nil
}

const baseURL = "http://localhost:3000"

func (i *Image) URL() string {
	if i == nil || i.Image == "" {
		return ""
	}
	return baseURL + "/uploads/" + filepath.Base(i.Image)
}
