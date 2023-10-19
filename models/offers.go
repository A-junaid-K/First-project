package models

import "time"

type Category_Offer struct {
	ID          uint `gorm:"primaryKey;unique"`
	Category_Id uint `gorm:"not null"`
	Offer       bool `gorm:"not null"`
	Percentage  uint
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
