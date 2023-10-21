package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

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
	err := database.DB.First(&dtuser, userid).Error
	if err != nil {
		c.HTML(404, "payment.html", gin.H{"error": "This user not found"})
		return
	}

	//get cart data
	var cartdata models.Cart
	err = database.DB.Where("user_id=?", userid).Find(&cartdata).Error
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": "Please check your cart"})
		return
	}

	//creating contact details
	result := database.DB.Where("user_id=?", userid).Create(&models.Contactdetails{
		Name:    name,
		Email:   email,
		User_ID: userid,
	})
	if result.Error != nil {
		c.HTML(400, "payment.html", gin.H{"error": "Failed to creat contact details"})
		return
	}

	//Add total amount
	var totalprice uint
	err = database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", userid).Scan(&totalprice).Error
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": "Failed to find the total price", "message": "please check your cart"})
		return
	}

	var address []models.Address
	database.DB.Where("user_id=?", userid).Find(&address)
	c.HTML(200, "payment.html", address)

	paymentMethod := c.PostForm("payment")
	cod := "cash-on-delivery"
	razorpay := "razorpay"

	if paymentMethod == cod {
		c.Redirect(303, "/user/checkout-cod")
	} else if paymentMethod == razorpay {
		c.Redirect(303, "/user/checkout-razorpay")
	} else {
		c.HTML(400, "payment.html", gin.H{"error": "Select any payment method"})
		return
	}

}

//-----------------------------------------Razor pay-------------------------------//

func Razorpay(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id
	//get the user
	var dtuser models.User
	err := database.DB.First(&dtuser, userid).Error
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": "This user didn't find"})
		return
	}

	//Add total amount
	var totalprice uint
	err = database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", userid).Scan(&totalprice).Error
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": "Failed to find the total price", "message": "please check your cart"})
		return
	}

	client := razorpay.NewClient(os.Getenv("RAZOR_kEY"), os.Getenv("RAZOR_SECRET"))
	data := map[string]interface{}{
		"amount":   totalprice,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": err})
		return
	}

	value := body["id"]
	c.HTML(http.StatusOK, "payment.html", gin.H{
		"userid":     userid,
		"totalprice": totalprice,
		"paymentid":  value,
	})

	var address []models.Address
	database.DB.Where("user_id=?", userid).Find(&address)
	c.HTML(200, "payment.html", address)

	c.Redirect(303, "/user/checkout-razorpay-success")

}

// -------------Razorpay Success------------------------//
func RazorpaySuccess(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	orderid := c.Query("order_id")
	paymentid := c.Query("payment_id")
	signature := c.Query("signature")
	totalamount := c.Query("total")

	err := database.DB.Create(&models.RazorPay{
		User_id:          uint(userid),
		RazorPayment_id:  paymentid,
		Signature:        signature,
		RazorPayOrder_id: orderid,
		AmountPaid:       totalamount,
	}).Error
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"Error": err.Error()})
		return
	}

	//searching for database all cart data
	var cartdata []models.Cart
	err = database.DB.Where("user_id=?", userid).Find(&cartdata).Error
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": "Please check your cart"})
		return
	}

	//getting total price of cart
	var totalprice uint
	err = database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", userid).Scan(&totalprice).Error
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": "Failed to find total price", "message": "cart is empty"})
		return
	}

	//checking stock level
	var product models.Product
	for _, v := range cartdata {
		database.DB.First(&product, v.Product_ID)
		if product.Stock-int(v.Quantity) < 0 {
			c.HTML(400, "payment.html", gin.H{
				"error": "Please check quantity",
			})
			return
		}
	}

	database.DB.Create(&models.Payment{
		Payment_Type:   "Razor pay",
		Total_Amount:   totalprice,
		Payment_Status: "Completed",
		User_ID:        userid,
		Date:           time.Now(),
	})
	var order models.Order
	var payment models.Payment
	database.DB.Last(&payment)
	var address models.Address
	err = database.DB.Where("user_id=? AND address_id=?", userid, order.Address_ID).First(&address).Error
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": "Failed to find address,choose different id"})
		return
	}

	err = database.DB.Create(&models.Order{
		User_ID:     userid,
		Address_ID:  order.Address_ID,
		Total_Price: totalprice,
		Payment_ID:  payment.Payment_ID,
		Status:      "Processing",
	}).Error
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": err.Error()})
		return
	}

	var cartbrand struct {
		Brand_Name    string
		Category_Name string
	}

	var order1 models.Order
	database.DB.Last(&order1)
	for _, cartdata := range cartdata {
		database.DB.Table("products").Select("brand_name,category_name").Where("id=?", cartdata.Product_ID).Scan(&cartbrand)
		err = database.DB.Create(&models.OrderItem{
			Order_ID:    order1.Order_ID,
			User_ID:     userid,
			Product_ID:  uint(cartdata.Product_ID),
			Address_ID:  order.Address_ID,
			Brand:       cartbrand.Brand_Name,
			Category:    cartbrand.Category_Name,
			Quantity:    cartdata.Quantity,
			Price:       uint(cartdata.Price),
			Total_Price: cartdata.Total_Price,
			Discount:    cartdata.Category_Offer + cartdata.Coupon_Discount,
			Cart_ID:     cartdata.ID,
			Status:      "processing",
			Created_at:  time.Now(),
		}).Error
		if err != nil {
			break
		}
	}
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": err.Error()})
		return
	}

	//reducing the stock count in database
	var products models.Product
	for _, v := range cartdata {
		database.DB.First(&products, v.Product_ID)
		database.DB.Model(&models.Product{}).Where("id=?", v.Product_ID).Update("stock", product.Stock-int(v.Quantity))
	}
	//deleting the checked out cart
	err = database.DB.Delete(&models.Cart{}, "user_id=?", userid).Error
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": "failed to delete used cart" + err.Error()})
		return
	}
	//giving success message
	c.HTML(200, "payment.html", gin.H{"message": "successfully ordered your cart"})

	var addresss []models.Address
	database.DB.Where("user_id=?", userid).Find(&addresss)
	c.HTML(200, "payment.html", address)

	c.Redirect(303, "/user/checkout-success")

}

//--------------------------------Success---------------------//

func Success(c *gin.Context) {

	pid, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	c.HTML(200, "success.html", gin.H{
		"paymentid": pid,
	})
}

//--------------------------COD------------------------------//

func Cod(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	//searching for database all cart data
	var cartdata []models.Cart
	err := database.DB.Where("user_id=?", userid).Find(&cartdata).Error
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": "Please check your cart"})
		return
	}

	//getting total price of cart
	var totalprice uint
	err = database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", userid).Scan(&totalprice).Error
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": "Failed to find total price", "message": "cart is empty"})
		return
	}

	//checking stock level
	var product models.Product
	for _, v := range cartdata {
		database.DB.First(&product, v.Product_ID)
		if product.Stock-int(v.Quantity) < 0 {
			c.HTML(400, "payment.html", gin.H{
				"error": "Please check quantity",
			})
			return
		}
	}

	database.DB.Create(&models.Payment{
		Payment_Type:   "COD",
		Total_Amount:   totalprice,
		Payment_Status: "Pending",
		User_ID:        userid,
		Date:           time.Now(),
	})

	var order models.Order
	var payment models.Payment
	database.DB.Last(&payment)
	var address models.Address
	err = database.DB.Where("user_id=? AND address_id=?", userid, order.Address_ID).First(&address).Error
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": "Failed to find address,choose different id"})
		return
	}

	err = database.DB.Create(&models.Order{
		User_ID:     userid,
		Address_ID:  order.Address_ID,
		Total_Price: totalprice,
		Payment_ID:  payment.Payment_ID,
		Status:      "Processing",
	}).Error
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": err.Error()})
		return
	}

	var cartbrand struct {
		Brand_Name    string
		Category_Name string
	}
	var order1 models.Order
	database.DB.Last(&order1)

	for _, cartdata := range cartdata {
		database.DB.Table("products").Select("brands.brand_name,categories.category_name").
			Joins("INNER JOIN brands ON brands.brand_id=products.brand_id").
			Joins("INNER JOIN categories ON categories.category_id=products.category_id").
			Where("id=?", cartdata.Product_ID).Scan(&cartbrand)
		err = database.DB.Create(&models.OrderItem{
			Order_ID:    order1.Order_ID,
			User_ID:     userid,
			Product_ID:  uint(cartdata.Product_ID),
			Address_ID:  order.Address_ID,
			Brand:       cartbrand.Brand_Name,
			Category:    cartbrand.Category_Name,
			Quantity:    cartdata.Quantity,
			Price:       uint(cartdata.Price),
			Total_Price: cartdata.Total_Price,
			Discount:    cartdata.Category_Offer + cartdata.Coupon_Discount,
			Cart_ID:     cartdata.ID,
			Status:      "processing",
			Created_at:  time.Now(),
		}).Error
		if err != nil {
			break
		}
	}
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": err.Error()})
		return
	}

	//reducing the stock count in database
	var products models.Product
	for _, v := range cartdata {
		database.DB.First(&products, v.Product_ID)
		database.DB.Model(&models.Product{}).Where("id=?", v.Product_ID).Update("stock", product.Stock-int(v.Quantity))
	}

	//deleting the checked out cart
	err = database.DB.Delete(&models.Cart{}, "user_id=?", userid).Error
	if err != nil {
		c.HTML(400, "payment.html", gin.H{"error": "failed to delete used cart" + err.Error()})
		return
	}

	//giving success message
	c.HTML(200, "payment.html", gin.H{"message": "successfully ordered your cart"})

	var addresss []models.Address
	database.DB.Where("user_id=?", userid).Find(&addresss)
	c.HTML(200, "payment.html", address)

	c.Redirect(303, "/user/checkout-success")

}
