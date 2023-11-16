package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func Userorder(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	type orderedItems struct {
		Image        string
		Name         string
		Price        uint
		Total_Price  uint
		Quantity     uint
		Created_at   time.Time
		Status       string
		Payment_Type string
		Order_ItemID uint
	}

	var order []orderedItems

	database.DB.Table("order_items").
		Select("products.image,products.name,order_items.price,order_items.total_price,order_items.quantity,order_items.created_at,orders.status,orders.payment_type,order_items.order_item_id").
		Joins("INNER JOIN products ON products.id=order_items.product_id").
		Joins("INNER JOIN orders ON orders.order_id=order_items.order_id").
		Where("order_items.user_id=?", userid).Scan(&order)

	c.HTML(200, "userorder.html", order)

}

func CancelOrder(c *gin.Context) {
	db := database.DB

	orderItemId, _ := strconv.Atoi(c.Param("orderitem_id"))
	var orderItem models.OrderItem

	// Ordered item data from table
	err := db.First(&orderItem, orderItemId).Error
	if err != nil {
		log.Println("Order id does not exist : ", err)
		return
	}

	//checking it already cancelled or not
	if orderItem.Status == "cancelled" {
		log.Println("Order already cancelled")
		c.HTML(http.StatusBadRequest, "userorder.html", gin.H{
			"error": "Order already cancelled",
		})
		return
	}

	// changing the order status in database
	db.Model(&models.Order{}).Where("order_id=?", orderItem.Order_ID).Update("status", "cancelled")
	err = db.Model(&models.OrderItem{}).Where("order_item_Id=?", orderItem.Order_ItemID).Update("status", "cancelled").Error
	if err != nil {
		log.Println("failed to cancel order in table : ", err)
		return
	}
	log.Println("db status updated : cancelled")

	// Update the Inventory
	var orderedProduct models.Product
	db.Where("id=?", orderItem.Product_ID).First(&orderedProduct)

	stock := orderedProduct.Stock + orderItem.Quantity

	err = db.Table("products").Where("id=?", orderItem.Product_ID).Update("stock", stock).Error
	if err != nil {
		log.Println("Failed to update stock in db : ", err)
		return
	}
	log.Println("updated stock : ", stock)

	//Refund
	var order models.Order
	db.Where("order_id=?", orderItem.Order_ID).First(&order)

	if order.Payment_Type != "COD" {
		var cancellingUser models.User
		db.Where("user_id=?", orderItem.User_ID).First(&cancellingUser)
		wallet := cancellingUser.Wallet + int(orderItem.Total_Price)
		err = db.Table("users").Where("user_id", orderItem.User_ID).Update("wallet", wallet).Error
		if err != nil {
			log.Println("Failed to update wallet in db : ", err)
			return
		}
		log.Println("waller updated : ", orderItem.Total_Price)

	}
	c.Redirect(303, "/user/orders")
}

func ReturnOrder(c *gin.Context) {

	db := database.DB

	orderItemId, _ := strconv.Atoi(c.Param("orderitem_id"))
	var orderItem models.OrderItem

	// Ordered item data from table
	err := db.First(&orderItem, orderItemId).Error
	if err != nil {
		log.Println("Order id does not exist : ", err)
		return
	}

	//checking it already cancelled or not
	if orderItem.Status == "cancelled" {
		log.Println("Order already cancelled")
		c.HTML(http.StatusBadRequest, "userorder.html", gin.H{
			"error": "Order already cancelled",
		})
		return
	}

	// changing the order status in database
	db.Model(&models.Order{}).Where("order_id=?", orderItem.Order_ID).Update("status", "cancelled")
	err = db.Model(&models.OrderItem{}).Where("order_item_Id=?", orderItem.Order_ItemID).Update("status", "cancelled").Error
	if err != nil {
		log.Println("failed to cancel order in table : ", err)
		return
	}
	log.Println("db status updated : cancelled")

	// Update the Inventory
	var orderedProduct models.Product
	db.Where("id=?", orderItem.Product_ID).First(&orderedProduct)

	stock := orderedProduct.Stock + orderItem.Quantity

	err = db.Table("products").Where("id=?", orderItem.Product_ID).Update("stock", stock).Error
	if err != nil {
		log.Println("Failed to update stock in db : ", err)
		return
	}
	log.Println("updated stock : ", stock)

	//Refund
	var order models.Order
	db.Where("order_id=?", orderItem.Order_ID).First(&order)

	if order.Payment_Type != "COD" {
		var cancellingUser models.User
		db.Where("user_id=?", orderItem.User_ID).First(&cancellingUser)
		wallet := cancellingUser.Wallet + int(orderItem.Total_Price)
		err = db.Table("users").Where("user_id", orderItem.User_ID).Update("wallet", wallet).Error
		if err != nil {
			log.Println("Failed to update wallet in db : ", err)
			return
		}
		log.Println("waller updated : ", orderItem.Total_Price)

	}
	c.Redirect(303, "/user/orders")

}
