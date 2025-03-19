package models

import (
	"time"
)

type Item struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ItemNumber  string    `json:"item_number" gorm:"unique"`
	Name        string    `json:"name"`
	Department  string    `json:"department"`
	Status      string    `json:"status" gorm:"default:'in_stock'"`  // in_stock, confirmed, signed
	Location    string    `json:"location"`
	Description string    `json:"description"`
	Orders      []Order   `json:"orders,omitempty" gorm:"foreignKey:ItemID"`
}

type Order struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	OrderNumber string    `json:"order_number" gorm:"unique"`
	ItemID      uint      `json:"item_id"`
	Item        Item      `json:"item" gorm:"foreignKey:ItemID"`
	Status      string    `json:"status" gorm:"default:'pending'"` // pending, signed
	Recipient   string    `json:"recipient"`
	Department  string    `json:"department"`
	Remarks     string    `json:"remarks"`
} 