package models

import (
	"database/sql"
	"strings"
	"time"
)

type Category struct {
	ID          int        `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

type CategoryManager struct {
	DB *sql.DB
}

func (cm *CategoryManager)SearchCategory(name string)(*Category, error)  {
	var category Category
	name = strings.ToLower(name)
	err := cm.DB.QueryRow(`SELECT * FROM categories WHERE name=$1`,name).Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt)
	if err != nil {
		
		return nil, err
	}
	return &category, nil
}