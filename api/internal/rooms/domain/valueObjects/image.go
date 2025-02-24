package valueobjects

type Image struct {
	image string
}

func NewImage(image string) (*Image, error) {

	return &Image{image: image}, nil
}

func (p *Image) Path() string {
	return p.image
}
