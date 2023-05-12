package models

import (
	"database/sql"
	"strings"
)

type Services struct {
	ID          int      `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	StaringAt   float64  `json:"staring_at,omitempty"`
	Description string   `json:"description,omitempty"`
	Category    Category `json:"category,omitempty"`
	Image       Image    `json:"image,omitempty"`
}
type NewService struct {
	Name      string   `json:"name,omitempty"`
	StaringAt int      `json:"staring_at,omitempty"`
	Image     []string `json:"image,omitempty"`
}

type ServiceManager struct {
	DB *sql.DB
}

func (sm *ServiceManager) CreateService(userId, name, description string, startingAt float64, imageID, categoryID int) (*Services, error) {
	name = strings.ToLower(name)
	description = strings.ToLower(description)
	var id int
	if err := sm.DB.QueryRow(`INSERT INTO services (user_id, name, starting_at, description, image_service_id, category_id)
	 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		userId, name, startingAt, description, imageID, categoryID).Scan(&id); err != nil {
		return nil, err
	}
	var service Services
	if err := sm.DB.QueryRow(`SELECT services.id, services.name, services.starting_at, services.description, categories.name, image_services.filename, 
	image_services.path FROM services INNER JOIN categories ON services.category_id = categories.id
	 INNER JOIN image_services ON services.image_service_id = image_services.id
	 WHERE services.id = $1`, id).Scan(&service.ID, &service.Name, &service.StaringAt,
		&service.Description, &service.Category.Name, &service.Image.FileName,
		&service.Image.Path); err != nil {
		return nil, err
	}
	return &service, nil
}
func (sm *ServiceManager) ListOfServices() ([]Services, error) {
	var list []Services
	row, err := sm.DB.Query(`SELECT * FROM services`)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var service Services
		if err := row.Scan(&service.ID, &service.Name, &service.StaringAt, &service.Description, &service.Image.ID, &service.Category.ID); err != nil {
			return nil, err
		}
		list = append(list, service)
	}
	if err := row.Err(); err != nil {
		return nil, err
	}
	return list, nil
}
func (sm *ServiceManager) DeleteService(id int) error {
	if _, err := sm.DB.Exec(`delete from services where id=$1`, id); err != nil {
		return err
	}
	return nil
}
func (sm *ServiceManager) UpdateService(id, name,description string, starting_at float64, categorieId,imageServiceId int  )(*Services, error) {
	var idService string
	if err := sm.DB.QueryRow("UPDATE services SET services.name=$1, services.starting_at=$2, services.description=$3, services.category_id=$4, services.image_service_id=$5 WHERE id=$6 RETURNING id", name, starting_at, categorieId,imageServiceId, id).Scan(&idService);  err != nil {
		return nil, err 
	}
	var service Services
	if err := sm.DB.QueryRow(`SELECT services.id, services.name, services.starting_at, services.description, categories.name, image_services.filename, 
	image_services.path FROM services INNER JOIN categories ON services.category_id = categories.id
	 INNER JOIN image_services ON services.image_service_id = image_services.id
	 WHERE services.id = $1`, idService).Scan(&service.ID, &service.Name, &service.StaringAt,
		&service.Description, &service.Category.Name, &service.Image.FileName,
		&service.Image.Path); err != nil {
		return nil, err
	}

	return &service, nil 
}
