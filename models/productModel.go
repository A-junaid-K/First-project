package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Description   string `gorm:"not null"`
	Stock         uint   `gorm:"not null"`
	Price         uint   `gorm:"not null"`
	Size          uint
	Category_Name string `gorm:"not null"`
	Brand_Name    string `gorm:"not null"`
	Image         string `gorm:"not null"`
	Offer_Name    string
	Percentage    uint
	Unlist      bool   `gorm:"default:false"`
}

type Brand struct {
	Brand_id   uint   `gorm:"primaryKey;unique"`
	Brand_Name string `gorm:"not null"`
}

type Category struct {
	Category_id uint   `gorm:"primaryKey;unique"`
	Name        string `gorm:"not null"`
	Unlist      bool   `gorm:"default:false"`
	Offer_Name  string
	Percentage  uint
}
