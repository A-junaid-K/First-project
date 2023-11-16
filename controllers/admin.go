package controllers

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/first_project/database"
	"github.com/first_project/helpers"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func Adminlogin(c *gin.Context) {
	c.HTML(http.StatusOK, "adminLogin.html", nil)
}
func PostAdminlogin(c *gin.Context) {
	type adminDetails struct {
		Email    string
		Password string
	}
	var admin adminDetails
	admin.Email = c.Request.FormValue("email")
	admin.Password = c.Request.FormValue("password")

	Email := os.Getenv("ADMIN_NAME")
	Password := os.Getenv("ADMIN_PASSWORD")
	// ---------cheking email & password----------
	if admin.Email != Email || admin.Password != Password {
		//log.Println("Unauthorized access invalid username or password")
		c.HTML(http.StatusUnauthorized, "adminLogin.html", gin.H{
			"error": "Unauthorized access invalid username or password",
		})
		return
	}
	//-----------generating token------------
	tokenstring, err := helpers.GenerateJWTToken(Email, "admin", 0)
	if err != nil {
		log.Println("Failed to generate jwt", err)
		return
	}

	//----------set token into browser-------
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt_admin", tokenstring, 3600*24*30, "", "", true, false)
	c.Redirect(303, "/admin-dashboard")
}
func Listusers(c *gin.Context) {
	type user struct {
		User_id   uint
		Name      string
		Email     string
		Phone     string
		IsBlocked bool
	}
	var users []user
	result := database.DB.Table("users").Select("user_id,name,email,phone,is_blocked").Scan(&users)
	if result.Error != nil {
		log.Println(result.Error)
	}
	c.HTML(http.StatusOK, "users-list.html", users)
}
func Blockuser(c *gin.Context) {
	//get the id
	id := c.Param("user_id")
	initid, _ := strconv.Atoi(id)
	var user models.User
	result := database.DB.First(&user, initid)
	if result.Error != nil {
		log.Println(result.Error)
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{
			"error": result.Error,
		})
		return
	}
	if user.IsBlocked {
		log.Println("This user is already blocked")
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{
			"error": "This user is already blocked",
		})
		return
	}
	//blocking the user
	result = database.DB.Model(&models.User{}).Where("user_id=?", initid).Update("is_blocked", true)
	if result.Error != nil {
		log.Println("Failed to blocking user")
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{
			"error": "Failed to blocking user",
		})
		return
	}
	c.Redirect(http.StatusSeeOther, "/users-list")
}
func Unblockuser(c *gin.Context) {

	//get the id
	id := c.Param("user_id")
	initid, _ := strconv.Atoi(id)
	var user models.User
	result := database.DB.First(&user, initid)
	if result.Error != nil {
		log.Println(result.Error)
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{
			"error": result.Error,
		})
		return
	}
	//checking user blocked or not
	if !user.IsBlocked {
		log.Println("This user is already unblocked")
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{
			"error": "This user is already unblocked",
		})
		return
	}
	//Unblocking the user
	result = database.DB.Model(&models.User{}).Where("user_id=?", initid).Update("is_blocked", false)
	if result.Error != nil {
		log.Println("Failed to unblock user")
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{
			"error": "Failed to unblock user",
		})
		return
	}
	// c.HTML(http.StatusOK, "users-list.html", gin.H{
	// 	"message": "Successfully unblocked" + user.Name,
	// })
	log.Println("Successfully unblocked " + user.Name)
	c.Redirect(http.StatusSeeOther, "/users-list")
}
func AdminDashboard(c *gin.Context) {
	db := database.DB

	// Revenue
	var revenue uint
	db.Table("payments").Select("SUM(total_amount)").Scan(&revenue)
	log.Println("revenue : ", revenue)

	// Pending Orders
	var pending_orders int64
	db.Model(&models.Order{}).Where("status=?", "pending").Count(&pending_orders)
	log.Println("PENDING ORDER : ", pending_orders)

	// Total products
	var products int64
	db.Model(&models.Product{}).Where("deleted_at IS NULL").Count(&products)
	log.Println("TOTAL PRODUCTS : ", products)

	// Users
	var users []models.User
	db.Find(&users)
	// log.Println("users : ", users)

	// Latest Orders
	var orders models.Order
	db.Last(&orders)

	var latest_orders []models.Order

	log.Println("order.userid : ", orders.Order_ID)

	for i := orders.Order_ID; i > orders.Order_ID-5; i-- {
		var order models.Order
		db.Where("order_id=?", i).Last(&order)
		if order.Order_ID != 0 {
			latest_orders = append(latest_orders, order)
		}
	}

	c.HTML(200, "adminDashboard.html", gin.H{
		"revenue":       revenue,
		"pendingorders": pending_orders,
		"totalproducts": products,
		"users":         users,
		"orders":        latest_orders,
	})

}

//--------------------------order---------------------//

func Order(c *gin.Context) {
	var order []models.Order
	database.DB.Find(&order)
	c.HTML(200, "order.html", order)
}
func PostOrder(c *gin.Context) {
	order_id, _ := strconv.Atoi(c.Param("order_id"))
	log.Println("order_id : ", order_id)
	status := c.PostForm("status")

	log.Println("status : ", status)

	//Update  Status in Orders
	err := database.DB.Table("orders").Select("status").Where("order_id=?", order_id).Updates(map[string]interface{}{
		"status": status,
		"date":   time.Now(),
	}).Error
	if err != nil {
		log.Println("failed to update order status : ", err)
	}

	//Update Status in Order_Items
	err = database.DB.Table("order_items").Select("status").Where("order_id=?", order_id).Updates(map[string]interface{}{
		"status": status,
		"created_at":   time.Now(),
	}).Error
	if err != nil {
		log.Println("failed to update order items status : ", err)
	}

	c.Redirect(303, "/admin-order")
}
