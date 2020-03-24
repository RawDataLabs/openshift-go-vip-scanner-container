package store

import (
	clair "github.com/rawc0der/VIP_Scanner/services/clair"
)

// Repository is a data access layer.
type Repository interface {
	Exists(name string) (bool, error)
	Create(podName string, imageName string) (*Image, error)
}

type Image struct {
	PodName         string   `json:"pod_name"`
	Image           string   `json:"image"`
	Vulnerabilities []string `json:"vulnerabilities"`
	Scanned         bool     `json:"scanned"`
}

// MemStore is a memroy storage for users.
type ImageStore struct {
	Images     []Image
	Repository *Repository
	Scanner    clair.Scannable
}

func NewImageStore() *ImageStore {
	var images = []Image{}
	return &ImageStore{images, new(Repository), clair.NewScanner()}
}

// Create creates user in the database for a form.
func (s *ImageStore) Create(podName string, dockerImage string) (*Image, error) {
	image := NewImage(podName, dockerImage)
	s.Images = append(s.Images, *image)
	return &*image, nil
}

// Exists checks if a email exists in the database.
func (s *ImageStore) Exists(image string) (bool, error) {
	for _, i := range s.Images {
		if i.Image == image {
			return true, nil
		}
	}
	return false, nil
}

func NewImage(podName string, dockerImage string) *Image {
	return &Image{podName, dockerImage, []string{}, false}
}

// func NewRepository() *Repository {
// 	return &Repository{}
// }
