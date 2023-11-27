package controllers

import (
	"fmt"
	"log"
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

	db := database.DB

	//Add total amount
	var totalprice uint

	row := db.Table("carts").Where("user_id=?", userid).Select("SUM(total_price)").Row()
	err := row.Scan(&totalprice)
	if err != nil {
		c.HTML(400, "app.html", gin.H{"error": "Failed to find the total price", "message": "please check your cart"})
		return
	}

	client := razorpay.NewClient(os.Getenv("RAZOR_kEY"), os.Getenv("RAZOR_SECRET"))
	data := map[string]interface{}{
		"amount":   totalprice * 100,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		fmt.Println("failed to get razor client : ", err)
		c.HTML(400, "app.html", gin.H{"error": err})
		return
	}
	value := body["id"]

	var contactdetails models.Contactdetails
	db.Where("user_id=?", userid).Last(&contactdetails)

	var address models.Address
	db.Where("address_id=?", contactdetails.Address_ID).Last(&address)

	c.HTML(http.StatusOK, "app.html", gin.H{
		"userid":     userid,
		"totalprice": totalprice,
		"paymentid":  value,
	})

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
		log.Println("Error : ", err)
		return
	}

	//searching for database all cart data
	var cartdata []models.Cart
	err = database.DB.Where("user_id=?", userid).Find(&cartdata).Error
	if err != nil {
		log.Println("error: Please check your cart : ", err)
		return
	}

	//getting total price of cart
	var totalprice uint
	err = database.DB.Table("carts").Select("SUM(total_price * quantity)").Where("user_id=?", userid).Scan(&totalprice).Error
	if err != nil {
		log.Println("error: Failed to find total price, message: cart is empty : ", err)
		return
	}
	var product models.Product

	//checking stock level
	for _, v := range cartdata {
		database.DB.First(&product, v.Product_ID)
		if int(product.Stock)-v.Quantity < 0 {
			log.Println("error: Please check quantity : ", err)
			return
		}
	}

	var adrid int
	err = database.DB.Model(&models.Contactdetails{}).Select("address_id").Where("user_id=?", userid).Scan(&adrid).Error
	if err != nil {
		log.Println("failed to fetch address id from checkout page : ", err)
		return
	}

	var order models.Order
	order.Address_ID = uint(adrid)

	err = database.DB.Create(&models.Payment{
		Payment_Type:   "RAZOR PAY",
		Total_Amount:   totalprice,
		Payment_Status: "Completed",
		User_ID:        userid,
		Date:           time.Now().Add(5*time.Hour + 30*time.Minute),
	}).Error
	if err != nil {
		log.Println("Failed to creat payment : ", err)
		return
	}

	var payment models.Payment
	database.DB.Last(&payment)
	var address models.Address
	err = database.DB.Where("user_id=? AND address_id=?", userid, order.Address_ID).First(&address).Error
	if err != nil {
		log.Println("error: Failed to find address : ", err)
		return
	}

	err = database.DB.Create(&models.Order{
		User_ID:      userid,
		Address_ID:   order.Address_ID,
		Total_Price:  totalprice,
		Payment_ID:   payment.Payment_ID,
		Status:       "processing",
		Date:         time.Now().Add(5*time.Hour + 30*time.Minute),
		Payment_Type: "RAZOR PAY",
	}).Error
	if err != nil {
		log.Println("failed to creat order : ", err)
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
			Quantity:    uint(cartdata.Quantity),
			Price:       uint(cartdata.Price),
			Total_Price: totalprice,
			Discount:    cartdata.Category_Offer + cartdata.Coupon_Discount,
			Cart_ID:     cartdata.ID,
			Status:      "processing",
			Created_at:  time.Now().Add(5*time.Hour + 30*time.Minute),
		}).Error
		if err != nil {
			log.Println("Failed to creat OrderItem : ", err)
			break
		}
	}
	if err != nil {
		log.Println("failed to range cartdata : ", err)
		return
	}

	//reducing the stock count in the database
	var products models.Product
	for _, v := range cartdata {
		database.DB.First(&products, v.Product_ID)
		database.DB.Model(&models.Product{}).Where("id=?", v.Product_ID).Update("stock", product.Stock-uint(v.Quantity))
	}

	//deleting the checked out cart
	err = database.DB.Delete(&models.Cart{}, "user_id=?", userid).Error
	if err != nil {
		log.Println("failed to delete used cart : ", err)
		return
	}

	// c.Redirect(303, "/user/payment-success")

	log.Println("message : successfully ordered your cart")

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"payment_id": payment.Payment_ID,
		"Message":    "Order Placed Successfully",
		"notice":     "Item removed from cart",
	})

}

//-------------------------------------------Instand Purchase-----------------------------------------//

func SingleRazorpay(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	totalprice := totalprice

	client := razorpay.NewClient(os.Getenv("RAZOR_kEY"), os.Getenv("RAZOR_SECRET"))
	data := map[string]interface{}{
		"amount":   totalprice * 100,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		fmt.Println("failed to get razor client : ", err)
		c.HTML(400, "singleApp.html", gin.H{"error": err})
		return
	}
	value := body["id"]

	c.HTML(http.StatusOK, "singleApp.html", gin.H{
		"userid":     userid,
		"totalprice": totalprice,
		"paymentid":  value,
	})
}

func SingleRazorpaySuccess(c *gin.Context) {

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
		log.Println("Error : ", err)
		return
	}

	//searching for database all cart data
	// var cartdata []models.Cart
	// err = database.DB.Where("user_id=?", userid).Find(&cartdata).Error
	// if err != nil {
	// 	log.Println("error: Please check your cart : ", err)
	// 	return
	// }

	// //getting total price of cart
	// var totalprice uint
	// err = database.DB.Table("carts").Select("SUM(total_price * quantity)").Where("user_id=?", userid).Scan(&totalprice).Error
	// if err != nil {
	// 	log.Println("error: Failed to find total price, message: cart is empty : ", err)
	// 	return
	// }

	totalprice := totalprice

	var product models.Product

	//checking stock level
	database.DB.First(&product, productId)

	if int(product.Stock)-qty < 0 {
		log.Println("error: Please check quantity : ")
		return
	}

	var adrid int
	err = database.DB.Model(&models.Contactdetails{}).Select("address_id").Where("user_id=?", userid).Scan(&adrid).Error
	if err != nil {
		log.Println("failed to fetch address id from checkout page in single razorpay success: ", err)
		return
	}

	var order models.Order
	order.Address_ID = uint(adrid)

	err = database.DB.Create(&models.Payment{
		Payment_Type:   "RAZOR PAY",
		Total_Amount:   totalprice,
		Payment_Status: "Completed",
		User_ID:        userid,
		Date:           time.Now(),
	}).Error
	if err != nil {
		log.Println("Failed to creat payment in single razorpay success: ", err)
		return
	}

	var payment models.Payment
	database.DB.Last(&payment)
	var address models.Address
	err = database.DB.Where("user_id=? AND address_id=?", userid, order.Address_ID).First(&address).Error
	if err != nil {
		log.Println("error: Failed to find address in single razorpay success: ", err)
		return
	}

	err = database.DB.Create(&models.Order{
		User_ID:      userid,
		Address_ID:   order.Address_ID,
		Total_Price:  totalprice,
		Payment_ID:   payment.Payment_ID,
		Status:       "processing",
		Date:         time.Now(),
		Payment_Type: "RAZOR PAY",
	}).Error
	if err != nil {
		log.Println("failed to creat order : ", err)
		return
	}

	var cartbrand struct {
		Brand_Name    string
		Category_Name string
	}

	var order1 models.Order
	database.DB.Last(&order1)

	database.DB.Table("products").Select("brand_name,category_name").Where("id=?", product.ID).Scan(&cartbrand)
	err = database.DB.Create(&models.OrderItem{
		Order_ID:    order1.Order_ID,
		User_ID:     userid,
		Product_ID:  uint(product.ID),
		Address_ID:  order.Address_ID,
		Brand:       cartbrand.Brand_Name,
		Category:    cartbrand.Category_Name,
		Quantity:    uint(qty),
		Price:       product.Price,
		Total_Price: totalprice,
		Discount:    0,
		Status:      "processing",
		Created_at:  time.Now(),
	}).Error

	if err != nil {
		log.Println("failed to range cartdata : ", err)
		return
	}

	//reducing the stock count in the database
	database.DB.Model(&models.Product{}).Where("id=?", product.ID).Update("stock", product.Stock-uint(qty))

	c.Redirect(303, "/user/payment-success")

}
