package models

import (
	"database/sql"
	"strings"
	"time"
)

type Category struct {
	ID          int        `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"-"`
	CreatedAt   *time.Time `json:"-"`
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

func (cm *CategoryManager)GetAllCategories() ([]Category, error) {
	var categories []Category
  rows, err := cm.DB.Query(`SELECT * FROM categories`)
  if err!= nil {
    return nil, err
  }
  for rows.Next() {
    var category Category
    err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt)
    if err!= nil {
      return nil, err
    }
    categories = append(categories, category)
  }
  return categories, nil

}