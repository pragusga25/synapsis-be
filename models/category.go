package models

import "time"

type Category struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `json:",omitempty"`
	UpdatedAt time.Time `json:",omitempty"`
	Name      string    `gorm:"type:varchar(120);not null;unique" json:",omitempty"`
}
