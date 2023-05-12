package models

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrUserDoNotExiste       = errors.New("user do not exist")
	ErrCantBuyYourOwnService = errors.New("you can't buy your own service")
)

type Delivery struct {
	ID          string     `json:"id,omitempty"`
	ServiceId   int        `json:"service_id,omitempty"`
	UserID      User       `json:"user_id,omitempty"`
	FreelanceID User       `json:"freelance_id,omitempty"`
	Price       float64    `json:"price,omitempty"`
	Start       *time.Time `json:"start,omitempty"`
	End         *time.Time `json:"end,omitempty"`
	Complete    bool       `json:"complete,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type DeliveryManager struct {
	DB *sql.DB
}

func (dm *DeliveryManager) CreateDelivery(service int, userID, freelanceID string, price float64, start, end time.Time) (*Delivery, error) {
	var user, freelance bool
	if err := dm.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)`, userID).Scan(&user); err != nil {
		return nil, ErrUserDoNotExiste
	}

	if err := dm.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)`, freelanceID).Scan(&freelance); err != nil {
		return nil, ErrUserDoNotExiste
	}
	if userID == freelanceID {
		return nil, ErrCantBuyYourOwnService
	}
	var d Delivery
	if err := dm.DB.QueryRow(`insert into deliveries (service_id, user_id,freelance_id, price, start, end_date ) 
	values ($1,$2,$3,$4,$5,$6) returning id, service_id, userID, freelanceID, price, start,end_date,complete created_at, updated_at
	 `, service, userID, freelanceID, price, start, end).Scan(
		d.ID, d.ServiceId, d.UserID, d.FreelanceID, d.Price, d.Start, d.End,d.Complete, d.CreatedAt, d.UpdatedAt); err != nil {
		return nil, err
	}
	return &d, nil
}
func (dm *DeliveryManager) DeleteDelivery(id string) error {
	if _, err := dm.DB.Exec(`delete from deliveries where id=$1`, id); err != nil {
		return err
	}
	return nil
}
func (dm *DeliveryManager)CompleteDelivery(id string)(*Delivery, error)  {
	var deliveries Delivery
	if err := dm.DB.QueryRow(`UPDATE deliveries SET complete=$1 where id=$2 RETURNING *`, true).Scan(&deliveries.ID, &deliveries.ServiceId, &deliveries.UserID, deliveries.FreelanceID, &deliveries.Price, &deliveries.Start, &deliveries.End, &deliveries.Complete, &deliveries.CreatedAt, &deliveries.UpdatedAt); err != nil {
		return nil, err 
	}
	return &deliveries, nil 
}
