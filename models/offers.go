package models

import "time"

type Category_Offer struct {
	ID            uint   `gorm:"primaryKey;unique"`
	Offer_Name    string `gorm:"not null"`
	Category_Name string `gorm:"not null"`
	Offer         bool   `gorm:"not null"`
	Percentage    uint
	Starting_Time time.Time `gorm:"not null"`
	Expiry_date   time.Time `gorm:"not null"`
}

type Coupon struct {
	CouponId      int       `gorm:"primaryKey;unique"`
	Coupon_Code   string    `gorm:"not null"`
	Starting_Time time.Time `gorm:"not null"`
	Ending_Time   time.Time `gorm:"not null"`
	Value         uint      `gorm:"not null"`
	Type          string    `gorm:"not null"`
	Max_Discount  uint      `gorm:"not null"`
	Min_Discount  uint      `gorm:"not null"`
	Cancel        bool      `gorm:"default:false"`
}
