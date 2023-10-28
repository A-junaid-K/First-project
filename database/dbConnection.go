package database

import (
	"log"
	"os"

	"github.com/first_project/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connetdb() {
	db_URL := os.Getenv("DSN")
	db, err := gorm.Open(postgres.Open(db_URL), &gorm.Config{})
	if err != nil {
		log.Panic(err)
		return
	}
	DB = db
}
func SyncDatabase() {

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Product{})
	DB.AutoMigrate(&models.Brand{})
	DB.AutoMigrate(&models.Category{})
	DB.AutoMigrate(&models.Cart{})
	DB.AutoMigrate(&models.Coupon{})
	DB.AutoMigrate(&models.Address{})
	DB.AutoMigrate(&models.Contactdetails{})
	DB.AutoMigrate(&models.Order{})
	DB.AutoMigrate(&models.OrderItem{})
	DB.AutoMigrate(&models.Coupon{})
	DB.AutoMigrate(&models.Category_Offer{})
	DB.AutoMigrate(&models.Payment{})
	DB.AutoMigrate(&models.RazorPay{})
	DB.AutoMigrate(&models.Wishlist{})

}
