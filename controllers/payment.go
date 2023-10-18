package controllers

import (
	"os"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

func Checkout(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id
	var address []models.Address
	database.DB.Where("user_id=?", userid).Find(&address)
	c.HTML(200, "payment.html", address)
}
func PostCheckout(c *gin.Context) {
	name := c.Request.FormValue("name")
	email := c.Request.FormValue("email")

	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	//get the user
	var dtuser models.User
	err := database.DB.First(&dtuser,userid).Error
	if err != nil {c.HTML(404, "payment.html", gin.H{"error": "This user not found",});return}

	//get cart data
	var cartdata models.Cart
	err = database.DB.Where("user_id=?", userid).Find(&cartdata).Error
	if err != nil {c.HTML(400, "payment.html", gin.H{"error": "Please check your cart",});return}

	//creating contact details
	result := database.DB.Where("user_id=?", userid).Create(&models.Contactdetails{
		Name:    name,
		Email:   email,
		User_ID: userid,
	})
	if result.Error != nil {c.HTML(400, "payment.html", gin.H{"error": "Failed to creat contact details",});return}
	
	//Add total amount
	var totalprice uint
	err= database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?",userid).Scan(&totalprice).Error
	if err != nil {c.HTML(400, "payment.html", gin.H{"error": "Failed to find the total price","message":"please check your cart",});return}

	paymentMethod := c.PostForm("payment")
	cod := "cash-on-delivery"
	razorpay := "razorpay"

	if paymentMethod == cod {
		Cod(c)
	} else if paymentMethod == razorpay {
		Razorpay(c)
	} else {
		c.HTML(400, "payment.html", gin.H{
			"error": "Select any payment method",
		})
		return
	}

}
func Razorpay(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id
	//get the user
	var dtuser models.User
	err := database.DB.First(&dtuser, userid).Error
	if err != nil {c.HTML(400,"razorpay.html" ,gin.H{"error": "This user didn't find",});return}

	client := razorpay.NewClient(os.Getenv("RAZOR_kEY"),os.Getenv("RAZOR_SECRET"))
	
}
func Cod(c *gin.Context) {

}
