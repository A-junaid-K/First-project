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
		Order_ID     uint
		Payment_ID   uint
		CTime        time.Time
	}

	var order []orderedItems

	database.DB.Table("order_items").
		Select("products.image,products.name,order_items.price,order_items.total_price,order_items.quantity,order_items.created_at,orders.status,orders.payment_type,order_items.order_item_id,orders.order_id,orders.payment_id").
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
	if orderItem.Status == "returned" {
		log.Println("Order already returned")
		c.HTML(http.StatusBadRequest, "return.html", gin.H{
			"error": "Order already returned",
		})
		return
	}

	//Check it Delivered
	var orders models.Order
	db.Where("order_id=?", orderItem.Order_ID).First(&orders)

	if orders.Status != "delivered" {
		log.Println("This order cannot return")
		c.HTML(http.StatusBadRequest, "return.html", gin.H{
			"error": "This Order cannot return",
		})
		return
	}

	// Checking the Returning Time period
	twoDaysAgo := time.Now().Add(-48 * time.Hour) // Subtract 48 hours for a two-day period
	if orderItem.Created_at.Before(twoDaysAgo) {
		log.Println("Cannot return. Returning time expired")
		c.HTML(400, "return.html", gin.H{"error": "Returning time expired.You can only return a product within two days"})
		return
	}

	// changing the order status in database
	db.Model(&models.Order{}).Where("order_id=?", orderItem.Order_ID).Updates(map[string]interface{}{
		"status": "returned",
		"date":   time.Now(),
	})

	err = db.Model(&models.OrderItem{}).Where("order_item_Id=?", orderItem.Order_ItemID).Updates(map[string]interface{}{
		"status":     "returned",
		"created_at": time.Now(),
	}).Error

	if err != nil {
		log.Println("failed to return order in table : ", err)
		return
	}
	log.Println("db status updated : returned")

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

	// Refund
	var order models.Order
	db.Where("order_id=?", orderItem.Order_ID).First(&order)

	// Calculate the time 7 days ago
	weekAgo := time.Now().Add(-7 * 24 * time.Hour)

	// Check if the product was returned at least 7 days ago
	if orderItem.Created_at.Before(weekAgo) {
		if order.Payment_Type != "COD" {
			var cancellingUser models.User
			db.Where("user_id=?", orderItem.User_ID).First(&cancellingUser)

			// Calculate the new wallet balance after the refund
			newWalletBalance := cancellingUser.Wallet + int(orderItem.Total_Price)

			// Update the user's wallet in the database
			err = db.Table("users").Where("user_id", orderItem.User_ID).Update("wallet", newWalletBalance).Error
			if err != nil {
				log.Println("Failed to update wallet in db : ", err)
				return
			}

			log.Println("Wallet updated for refund : ", orderItem.Total_Price)
		}
	} else {
		log.Println("Refund period has not elapsed yet")
	}

	log.Println("Successfully Returned")

	c.Redirect(303, "/user/orders")

}

func Reason(c *gin.Context) {

	orderItemId, _ := strconv.Atoi(c.Param("orderitem_id"))

	log.Println("order item id : ", orderItemId)

	c.HTML(200, "return.html", gin.H{
		"orderitemid": orderItemId,
	})

}
