package models

import (
	"time"
)

type User struct {
	User_id    uint   `gorm:"primaryKey;unique"`
	Name       string `validate:"required, min=2,max=50"`
	Password   string `validate:"required, min=6"`
	Email      string `gorm:"not null"`
	Phone      string
	IsBlocked  bool
	Otp        string
	Validate   bool
	User_type  string `validate:"required, eq=ADMIN|eq=USER"`
	Created_at time.Time
	Wallet     int
}

type Address struct {
	Address_ID    uint   `gorm:"primaryKey;unique"`
	Building_Name string `gorm:"not null"`
	City          string `gorm:"not null"`
	State         string `gorm:"not null"`
	Landmark      string `gorm:"not null"`
	Zip_code      string `gorm:"not null"`
	User_ID       uint   `gorm:"not null"`
	Primary       bool
}
type Contactdetails struct {
	Contactdetails_id int    `gorm:"primaryKey"`
	Name              string `validate:"required, min=2,max=50"`
	Phone             string
	Email             string `gorm:"not null"`
	Address_ID        uint   `gorm:"not null"`
	Payment_Method    string
	User_ID           uint `gorm:"not null"`
}
