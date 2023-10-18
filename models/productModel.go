package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Description   string `gorm:"not null"`
	Stock         int    `gorm:"not null"`
	Price         int    `gorm:"not null"`
	Size          int
	Category_Name string `gorm:"not null"`
	Brand_Name    string `gorm:"not null"`
	Image         string `gorm:"not null"`
}
type Brand struct {
	Brand_id   uint   `gorm:"primaryKey;unique"`
	Brand_Name string `gorm:"not null"`
}
type Category struct {
	Category_id uint   `gorm:"primaryKey;unique"`
	Name        string `gorm:"not null"`
	Unlist      bool   `gorm:"default:false"`
}
// type Image struct {
// 	Id         uint   `gorm:"primaryKey;unique"`
// 	Product_id uint   `gorm:"not null"`
// 	Image      string `gorm:"not null"`
// }
