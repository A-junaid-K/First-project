package models

import "time"

type Order struct {
	Order_ID     uint `gorm:"primaryKey;unique"`
	User_ID      uint `gorm:"not null"`
	Address_ID   uint `gorm:"not null"`
	Total_Price  uint `gorm:"not null"`
	Payment_Type string
	Payment_ID   uint   `gorm:"not null"`
	Status       string `gorm:"not null"`
	Date         time.Time
}

type OrderItem struct {
	Order_ItemID uint   `gorm:"primaryKey;unique"`
	User_ID      uint   `gorm:"not null"`
	Order_ID     uint   `gorm:"not null"`
	Product_ID   uint   `gorm:"not null"`
	Address_ID   uint   `gorm:"not null"`
	Category     string `gorm:"not null"`
	Brand        string `gorm:"not null"`
	Quantity     uint   `gorm:"not null"`
	Price        uint   `gorm:"not null"`
	Total_Price  uint   `gorm:"not null"`
	Discount     uint
	Cart_ID      uint   `gorm:"not null"`
	Status       string `gorm:"not null"`
	Created_at   time.Time
}

type Payment struct {
	Payment_ID     uint   `gorm:"primaryKey;unique"`
	Payment_Type   string `gorm:"not null"`
	Total_Amount   uint   `gorm:"not null"`
	Payment_Status string `gorm:"not null"`
	User_ID        uint   `gorm:"not null"`
	Date           time.Time
}

type RazorPay struct {
	User_id          uint
	RazorPayment_id  string `gorm:"primaryKey"`
	RazorPayOrder_id string
	Signature        string
	AmountPaid       string
}
