package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

// func RazorPay(c *gin.Context){

// }

func RazorPay(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	//Add total amount
	var totalprice uint
	err := database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", userid).Scan(&totalprice).Error
	if err != nil {
		c.HTML(400, "checkout.html", gin.H{"error": "Failed to find the total price", "message": "please check your cart"})
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
		fmt.Println("failed to get razor client : ", err)
		c.HTML(400, "checkout.html", gin.H{"error": err})
		return
	}
	value := body["id"]
	c.HTML(http.StatusOK, "checkout.html", gin.H{
		"userid":     userid,
		"totalprice": totalprice,
		"paymentid":  value,
	})

	c.Redirect(303, "/user/payment-razorpay-success")

}
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
		c.HTML(400, "checkout.html", gin.H{"Error": err.Error()})
		return
	}

	//searching for database all cart data
	var cartdata []models.Cart
	err = database.DB.Where("user_id=?", userid).Find(&cartdata).Error
	if err != nil {
		c.HTML(400, "checkout.html", gin.H{"error": "Please check your cart"})
		return
	}

	//getting total price of cart
	var totalprice uint
	err = database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", userid).Scan(&totalprice).Error
	if err != nil {
		c.HTML(400, "checkout.html", gin.H{"error": "Failed to find total price", "message": "cart is empty"})
		return
	}

	//checking stock level
	var product models.Product
	for _, v := range cartdata {
		database.DB.First(&product, v.Product_ID)
		if product.Stock-int(v.Quantity) < 0 {
			c.HTML(400, "checkout.html", gin.H{
				"error": "Please check quantity",
			})
			return
		}
	}

	var adrid int
	err = database.DB.Model(&models.Contactdetails{}).Select("address_id").Where("user_id=?", userid).Scan(&adrid).Error
	if err != nil {
		fmt.Println("failed to fetch address id from checkout page")
		c.HTML(400, "checkout.html", gin.H{"error": "Failed to find address,choose different id"})
		return
	}

	var order models.Order
	order.Address_ID = uint(adrid)

	database.DB.Create(&models.Payment{
		Payment_Type:   "RAZOR PAY",
		Total_Amount:   totalprice,
		Payment_Status: "Completed",
		User_ID:        userid,
		Date:           time.Now(),
	})
	var payment models.Payment
	database.DB.Last(&payment)
	var address models.Address
	err = database.DB.Where("user_id=? AND address_id=?", userid, order.Address_ID).First(&address).Error
	if err != nil {
		c.HTML(400, "checkout.html", gin.H{"error": "Failed to find address,choose different id"})
		return
	}

	err = database.DB.Create(&models.Order{
		User_ID:     userid,
		Address_ID:  order.Address_ID,
		Total_Price: totalprice,
		Payment_ID:  payment.Payment_ID,
		Status:      "processing",
	}).Error
	if err != nil {
		fmt.Println("failed to creat order")
		c.HTML(400, "checkout.html", gin.H{"error": err.Error()})
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
		fmt.Println("failed to range cartdata")
		c.HTML(400, "checkout.html", gin.H{"error": err.Error()})
		return
	}

	//reducing the stock count in the database
	var products models.Product
	for _, v := range cartdata {
		database.DB.First(&products, v.Product_ID)
		database.DB.Model(&models.Product{}).Where("id=?", v.Product_ID).Update("stock", product.Stock-int(v.Quantity))
	}

	//deleting the checked out cart
	err = database.DB.Delete(&models.Cart{}, "user_id=?", userid).Error
	if err != nil {
		c.HTML(400, "checkout.html", gin.H{"error": "failed to delete used cart" + err.Error()})
		return
	}
	// //giving success message
	// c.JSON(200, gin.H{
	// 	"message": "successfully ordered your cart",
	// })


	var carts models.Cart
	var addresses models.Address
	var contactdetailss models.Contactdetails
	// var orders models.Order

	database.DB.Where("user_id=",userid).Last(&contactdetailss)
	database.DB.Where("address_id=?",adrid).Last(&addresses)
	database.DB.Find(&carts)

	rzrpaydata := struct {
		Carts      models.Cart
		Addresses  models.Address
		Contactdetails models.Contactdetails
	}{
		Carts:      carts,
		Addresses:  addresses,
		Contactdetails: contactdetailss,
	}


	c.HTML(200,"razorpay.html",rzrpaydata)

}
