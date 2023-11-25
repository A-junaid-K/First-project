package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func Checkout(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	//get cart data
	var cartdata []models.Cart
	err := database.DB.Where("user_id=?", userid).Find(&cartdata).Error
	if err != nil {
		c.HTML(400, "checkout.html", gin.H{"error": "Please check your cart"})
		return
	}

	var userr models.User
	var adr []models.Address

	database.DB.Where("user_id=?", userid).First(&userr)
	database.DB.Where("user_id=?", userid).Find(&adr)

	//Add total amount
	var totalprice uint
	err = database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", userid).Scan(&totalprice).Error
	if err != nil {
		c.HTML(400, "checkout.html", gin.H{"error": "Failed to find the total price", "message": "please check your cart"})
		return
	}

	c.HTML(200, "checkout.html", gin.H{
		"Users":      userr,
		"Addresses":  adr,
		"Carts":      cartdata,
		"totalprice": totalprice,
	})
}

func PostCheckout(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	// Recieve user details from Front-end
	name := c.Request.FormValue("name")
	email := c.Request.FormValue("email")
	phone := c.Request.FormValue("number")

	log.Println("name : ", name)
	log.Println("phone : ", phone)

	err := namevalidator(name)
	if err != nil {
		c.HTML(http.StatusBadRequest, "checkout.html", gin.H{"error": err})
	}
	err = emailvalidator(email)
	if err != nil {
		c.HTML(http.StatusBadRequest, "checkout.html", gin.H{"error": err})
	}
	err = numbervalidator(phone)
	if err != nil {
		c.HTML(http.StatusBadRequest, "checkout.html", gin.H{"error": err})
	}

	//get the user
	var dtuser models.User
	err = database.DB.First(&dtuser, userid).Error
	if err != nil {
		c.HTML(404, "checkout.html", gin.H{"error": "This user not found"})
		return
	}

	//Fetching cart data and Count
	var cartdata []models.Cart
	database.DB.Where("user_id=?", userid).Find(&cartdata)

	for _, v := range cartdata {
		if v.Category_Offer != 0 {

			discount := v.Total_Price * v.Category_Offer / 100
			newprice := v.Total_Price - discount

			database.DB.Model(&models.Cart{}).Where("user_id=? AND id=?", userid, v.ID).Update("total_price", newprice)
			log.Println("Offer price updated")
		}
	}

	// recieving the address
	adrid, _ := strconv.Atoi(c.PostForm("userchosenaddress"))
	var postadr models.Address
	err = database.DB.Where("address_id=?", adrid).First(&postadr).Error
	if err != nil {
		fmt.Println("errr in cod")
	}

	//user chosen payment method
	paymentMethod := c.PostForm("payment")
	cod := "cash-on-delivery"
	razorpay := "razorpay"
	wallet := "wallet"

	//creating contact details
	result := database.DB.Where("user_id=?", userid).Create(&models.Contactdetails{
		Name:           name,
		Email:          email,
		Address_ID:     uint(adrid),
		Payment_Method: paymentMethod,
		User_ID:        userid,
		Phone:          phone,
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
	} else if paymentMethod == wallet {
		c.Redirect(303, "/user/payment-wallet")
	}
}

var totalprice uint
var productId int
var qty int

func BuyNow(c *gin.Context) {

	// Retrieving user details
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	var userr models.User
	var adr []models.Address

	database.DB.Where("user_id=?", userid).First(&userr)
	database.DB.Where("user_id=?", userid).Find(&adr)

	// Retreive product details
	productId, _ = strconv.Atoi(c.Param("product_id"))

	// Get the size & quantity
	pdsize := c.Query("size")
	pdquantity, _ := strconv.Atoi(c.Query("qty"))
	log.Println("pdsize : ", pdsize)
	log.Println("pdquantity : ", pdquantity)

	qty = pdquantity

	log.Println("PRODFSD ID : ", productId)

	var product models.Product
	database.DB.First(&product, productId)

	totalprice = product.Price * uint(pdquantity)

	log.Println("total price : ", totalprice)

	c.HTML(200, "singleCheckout.html", gin.H{
		"Users":      userr,
		"Addresses":  adr,
		"Product":    product,
		"totalprice": totalprice,
	})

}

func PostBuyCheckout(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	// Recieve user details from Front-end
	name := c.Request.FormValue("name")
	email := c.Request.FormValue("email")
	phone := c.Request.FormValue("number")

	log.Println("name : ", name)
	log.Println("phone : ", phone)

	err := namevalidator(name)
	if err != nil {
		c.HTML(http.StatusBadRequest, "singleCheckout.html", gin.H{"error": err})
	}
	err = emailvalidator(email)
	if err != nil {
		c.HTML(http.StatusBadRequest, "singleCheckout.html", gin.H{"error": err})
	}
	err = numbervalidator(phone)
	if err != nil {
		c.HTML(http.StatusBadRequest, "singleCheckout.html", gin.H{"error": err})
	}

	// recieving the address
	adrid, _ := strconv.Atoi(c.PostForm("userchosenaddress"))
	var postadr models.Address
	err = database.DB.Where("address_id=?", adrid).First(&postadr).Error
	if err != nil {
		fmt.Println("errr in cod")
	}

	totalprice := totalprice

	var product models.Product
	database.DB.First(&product, productId)

	if product.Percentage != 0 {
		discount := totalprice * product.Percentage / 100
		newprice := totalprice - discount
		log.Println("newprice : ", newprice)
	}

	//user chosen payment method
	paymentMethod := c.PostForm("payment")
	cod := "cash-on-delivery"
	razorpay := "razorpay"
	wallet := "wallet"

	//creating contact details
	result := database.DB.Where("user_id=?", userid).Create(&models.Contactdetails{
		Name:           name,
		Email:          email,
		Address_ID:     uint(adrid),
		Payment_Method: paymentMethod,
		User_ID:        userid,
		Phone:          phone,
	})
	if result.Error != nil {
		c.HTML(400, "checkout.html", gin.H{"error": "Failed to creat contact details"})
		return
	}

	//Redirecting to chosen payment method
	if paymentMethod == cod {
		c.Redirect(303, "/user/payment-single-cod")
	} else if paymentMethod == razorpay {
		c.Redirect(303, "/user/payment-single-razorpay")
	} else if paymentMethod == wallet {
		c.Redirect(303, "/user/payment-single-wallet")
	}
}
