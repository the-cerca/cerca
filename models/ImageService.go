package models

import (
	"database/sql"
	"time"
)

type Image struct {
	ID        int        `json:"id,omitempty"`
	FileName  string     `json:"file_name,omitempty"`
	Path      string     `json:"path,omitempty"`
	Alt       string     `json:"alt,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type ImageManager struct {
	DB *sql.DB ``
}

// Create a function that allows adding the path of the images to the database
func (im *ImageManager) CreateToDatabase(filename, path, alt_text string) (*Image, error) {
	// Execute the INSERT statement and retrieve the last inserted ID
	var image Image
	err := im.DB.QueryRow(`INSERT INTO image_services (filename, path, alt_text) VALUES ($1,$2,$3)
	 RETURNING id, filename, path, alt_text, created_at`, filename, path, alt_text).Scan(&image.ID, &image.FileName, &image.Path, &image.Alt, &image.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

