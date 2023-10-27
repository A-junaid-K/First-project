package controllers

import (
	"fmt"
	"strconv"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func Checkout(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id
	var address []models.Address
	database.DB.Where("user_id=?", userid).Find(&address)
	c.HTML(200, "checkout.html", address)
}
func PostCheckout(c *gin.Context) {
	name := c.Request.FormValue("name")
	email := c.Request.FormValue("email")
	user, _ := c.Get("user")
	userid := user.(models.User).User_id
	//get the user
	var dtuser models.User
	err := database.DB.First(&dtuser, userid).Error
	if err != nil {
		c.HTML(404, "checkout.html", gin.H{"error": "This user not found"})
		return
	}

	//get cart data
	var cartdata models.Cart
	err = database.DB.Where("user_id=?", userid).Find(&cartdata).Error
	if err != nil {
		c.HTML(400, "checkout.html", gin.H{"error": "Please check your cart"})
		return
	}

	//Add total amount
	var totalprice uint
	err = database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", userid).Scan(&totalprice).Error
	if err != nil {
		c.HTML(400, "checkout.html", gin.H{"error": "Failed to find the total price", "message": "please check your cart"})
		return
	}

	// recieving the address
	adrid, _ := strconv.Atoi(c.PostForm("userchosenaddress"))
	fmt.Println("address getting from checkout : ", adrid)
	var postadr models.Address
	err = database.DB.Where("address_id=?", adrid).First(&postadr).Error
	if err != nil {
		fmt.Println("errr in cod")
	}

	//user chosen payment method
	paymentMethod := c.PostForm("payment")
	cod := "cash-on-delivery"
	razorpay := "razorpay"

	//creating contact details
	result := database.DB.Where("user_id=?", userid).Create(&models.Contactdetails{
		Name:           name,
		Email:          email,
		Address_ID:     uint(adrid),
		Payment_Method: paymentMethod,
		User_ID:        userid,
	})
	if result.Error != nil {
		c.HTML(400, "checkout.html", gin.H{"error": "Failed to creat contact details"})
		return
	}

	//Redirecting to chosen payment method
	if paymentMethod == cod {
		// c.Redirect(303, "/user/payment-success")
		c.Redirect(303, "/user/payment-cod")
	} else if paymentMethod == razorpay {
		c.Redirect(303, "/user/payment-razorpay")
	} else {
		c.HTML(400, "checkout.html", gin.H{"error": "Select any payment method"})
		return
	}

}
