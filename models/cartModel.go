package models

type Cart struct {
	ID              uint   `gorm:"primaryKey;unique"`
	Name            string `gorm:"not null"`
	Description     string `gorm:"not null"`
	Stock           int    `gorm:"not null"`
	Price           int    `gorm:"not null"`
	Size            string `gorm:"not null"`
	Category_Name   string `gorm:"not null"`
	Brand_Name      string `gorm:"not null"`
	Product_ID      int    `gorm:"not null"`
	Quantity        int   `gorm:"not null"`
	Total_Price     uint   `gorm:"not null"`
	User_ID         uint   `gorm:"not null"`
	Coupon_Applied  bool   `gorm:"default:false"`
	Coupon_Discount uint
	Category_Offer  uint
	Image           string `gorm:"not null"`
}

type Wishlist struct {
	ID              uint   `gorm:"primaryKey;unique"`
	Name            string `gorm:"not null"`
	Description     string `gorm:"not null"`
	Stock           int    `gorm:"not null"`
	Price           int    `gorm:"not null"`
	Size            string `gorm:"not null"`
	Category_Name   string `gorm:"not null"`
	Brand_Name      string `gorm:"not null"`
	Product_ID      int    `gorm:"not null"`
	Quantity        uint   `gorm:"not null"`
	Total_Price     uint   `gorm:"not null"`
	User_ID         uint   `gorm:"not null"`
	Coupon_Applied  bool   `gorm:"default:false"`
	Coupon_Discount uint
	Category_Offer  uint
	Image           string `gorm:"not null"`
}
