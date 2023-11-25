package controllers

import (
	"log"
	"strconv"
	"time"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func Wallet(c *gin.Context) {
	//Initializing DB
	db := database.DB

	//find user
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	// Get the wallet from front-end
	wallet, _ := strconv.Atoi(c.Request.FormValue("wallet"))
	log.Println("wallet  : ", wallet)
	if wallet == 0 {
		c.Next()
		return
	}

	// Get the cart data
	var cartdata []models.Cart
	db.Where("user_id=?", userid).Find(&cartdata)

	// Get the user data
	var userwallet models.User
	db.Where("user_id=?", userid).Find(&userwallet)

	//getting total price of cart
	var totalprice uint
	err := db.Table("carts").Select("SUM(total_price)").Where("user_id=?", userid).Scan(&totalprice).Error
	if err != nil {
		log.Println("Failed to find total price : ", err)
		c.HTML(400, "checkout.html", gin.H{"error": "Failed to find total price"})
		return
	}

	log.Println("totalprice : ", totalprice)

	// Validate user entered Wallet amount
	if wallet >= int(totalprice) {

		// Validate Wallet balance
		if userwallet.Wallet < int(totalprice) {
			log.Println("Insufficient Funds")
			c.HTML(400, "checkout.html", gin.H{"error": "Sorry, you don't have enough money in your wallet"})
			return
		}
		log.Println("pay full with wallet")
		c.Redirect(303, "/user/payment-wallet")
		return
	} else {

		//Validte Wallet balance
		if userwallet.Wallet < wallet {
			log.Println("Insufficient Funds")
			c.HTML(400, "checkout.html", gin.H{"error": "Sorry, you don't have enough money in your wallet"})
			return
		}

		walletprice := totalprice - uint(wallet)
		var cartitems int64
		database.DB.Model(&models.Cart{}).Where("user_id=?", userid).Count(&cartitems)
		for _, v := range cartdata {
			newprice := walletprice / uint(cartitems)
			err := database.DB.Model(&models.Cart{}).Where("user_id=? AND id=?", userid, v.ID).Updates(map[string]interface{}{"total_price": newprice}).Error
			if err != nil {
				log.Println(err)
				return
			}
		}

		userwallet.Wallet -= wallet
		if err := database.DB.Save(&userwallet).Error; err != nil {
			log.Println("Failed to update user wallet: ", err)
			c.HTML(500, "wallet.html", gin.H{"error": "Failed to update user wallet"})
			return
		}

		log.Println("Wallet amount apllied Successfully")
		c.Next()
	}

}

//----------------------Payment with Wallet--------------//

func PaywithWallet(c *gin.Context) {
	//Initializing db
	db := database.DB

	// Find user
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	var wallet models.User
	db.First(&wallet, userid)

	// Retrieve cart data
	var cartdata models.Cart
	db.Find(&cartdata, userid)

	// Fetch Total price from cart
	var totalprice uint
	err := db.Table("carts").Select("SUM(total_price)").Where("user_id=?", userid).Scan(&totalprice).Error
	if err != nil {
		log.Println("Failed to find total price : ", err)
		c.HTML(400, "checkout.html", gin.H{"error": "Failed to find total price"})
		return
	}

	// Validate Wallet balance
	if wallet.Wallet < int(totalprice) {
		log.Println("Insufficient Funds")
		c.HTML(400, "checkout.html", gin.H{"error": "Sorry, you don't have enough money in your wallet"})
		return
	}

	// Fetch the payment from database
	var payment models.Payment
	database.DB.Last(&payment)

	c.HTML(200, "wallet.html", gin.H{
		"userid":     userid,
		"paymentid":  payment.Payment_ID + 1,
		"totalprice": totalprice,
	})
}

//--------------Wallet Success----------------//

func WalletSuccess(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	//Retrieve cart data from DB
	var cartdata []models.Cart
	err := database.DB.Where("user_id=?", userid).Find(&cartdata).Error
	if err != nil {
		log.Println("Please check your cart : ", err)
		c.HTML(400, "wallet.html", gin.H{"error": "Please check your cart"})
		return
	}

	//getting total price of cart
	var totalprice uint
	err = database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", userid).Scan(&totalprice).Error
	if err != nil {
		log.Println("Failed to find total price")
		c.HTML(400, "wallet.html", gin.H{"error": "Failed to find total price"})
		return
	}

	var product models.Product

	//checking stock level
	for _, v := range cartdata {
		database.DB.First(&product, v.Product_ID)
		level := int(product.Stock) - v.Quantity
		if int(level) < 0 {
			log.Println("error : please check quantity : ", err)
			c.HTML(400, "wallet.html", gin.H{
				"error": "Please check quantity",
			})
			return
		}
	}

	// Retrieve address id from DB
	var adrid int
	err = database.DB.Model(&models.Contactdetails{}).Select("address_id").Where("user_id=?", userid).Scan(&adrid).Error
	if err != nil {
		log.Println("failed to fetch address id from checkout page")
		c.HTML(400, "wallet.html", gin.H{"error": "Failed to find address,choose different id"})
		return
	}

	var order models.Order
	order.Address_ID = uint(adrid)

	var payment models.Payment
	database.DB.Last(&payment)
	var address models.Address
	err = database.DB.Where("user_id=? AND address_id=?", userid, order.Address_ID).Last(&address).Error
	if err != nil {
		c.HTML(400, "wallet.html", gin.H{"error": "Failed to find address,choose different id"})
		return
	}

	err = database.DB.Create(&models.Order{
		User_ID:      userid,
		Address_ID:   order.Address_ID,
		Total_Price:  totalprice,
		Payment_ID:   payment.Payment_ID,
		Status:       "Processing",
		Payment_Type: "Wallet",
		Date:         time.Now(),
	}).Error
	if err != nil {
		log.Println("failed to create order")
		c.HTML(500, "wallet.html", gin.H{"error": err.Error()})
		return
	}

	//creating Wallet
	database.DB.Create(&models.Payment{
		Payment_Type:   "Wallet",
		Total_Amount:   totalprice,
		Payment_Status: "Completed",
		User_ID:        userid,
		Date:           time.Now(),
	})

	var order1 models.Order
	database.DB.Last(&order1)

	for _, cartdata := range cartdata {

		err = database.DB.Create(&models.OrderItem{
			Order_ID:    order1.Order_ID,
			User_ID:     userid,
			Product_ID:  uint(cartdata.Product_ID),
			Address_ID:  order.Address_ID,
			Brand:       cartdata.Brand_Name,
			Category:    cartdata.Category_Name,
			Quantity:    uint(cartdata.Quantity),
			Price:       uint(cartdata.Price),
			Total_Price: totalprice,
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
		log.Println("Failed to create Order Item", err)
		c.HTML(400, "wallet.html", gin.H{"error": err.Error()})
		return
	}

	// Deduct the total price from the user's wallet
	var userwallet models.User
	database.DB.Where("user_id=?", userid).Find(&userwallet)

	userwallet.Wallet -= int(totalprice)
	if err := database.DB.Save(&userwallet).Error; err != nil {
		log.Println("Failed to update user wallet: ", err)
		c.HTML(500, "wallet.html", gin.H{"error": "Failed to update user wallet"})
		return
	}

	//reducing the stock count in database
	var products models.Product
	for _, v := range cartdata {
		database.DB.First(&products, v.Product_ID)
		err = database.DB.Model(&models.Product{}).Where("id=?", v.Product_ID).Update("stock", products.Stock-uint(v.Quantity)).Error
		if err != nil {
			log.Println("failed to update stock in database : ", err)
		}
	}

	//deleting the checked out cart
	err = database.DB.Delete(&models.Cart{}, "user_id=?", userid).Error
	if err != nil {
		log.Println("Failed to delete checked out cart data")
		c.HTML(400, "wallet.html", gin.H{"error": "failed to delete used cart" + err.Error()})
		return
	}

	c.Redirect(303, "/user/payment-success")
}

//----------------------------------------------//

func SinglePaywithWallet(c *gin.Context) {

	//Initializing db
	db := database.DB

	// Find user
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	var wallet models.User
	db.First(&wallet, userid)

	// Fetch Total price from cart
	totalprice := totalprice

	// Validate Wallet balance
	if wallet.Wallet < int(totalprice) {
		log.Println("Insufficient Funds")
		c.HTML(400, "singleCheckout.html", gin.H{"error": "Sorry, you don't have enough money in your wallet"})
		return
	}

	// Fetch the payment from database
	var payment models.Payment
	database.DB.Last(&payment)

	c.HTML(200, "singleWallet.html", gin.H{
		"userid":     userid,
		"paymentid":  payment.Payment_ID + 1,
		"totalprice": totalprice,
	})
}

func SingleWalletSuccess(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	totalprice := totalprice

	var product models.Product
	database.DB.First(&product, productId)

	//checking stock level
	level := int(product.Stock) - qty
	if int(level) < 0 {
		log.Println("error : please check quantity")
		c.HTML(400, "wallet.html", gin.H{
			"error": "Please check quantity",
		})
		return
	}

	// Retrieve address id from DB
	var adrid int
	err := database.DB.Model(&models.Contactdetails{}).Select("address_id").Where("user_id=?", userid).Scan(&adrid).Error
	if err != nil {
		log.Println("failed to fetch address id from checkout page")
		c.HTML(400, "wallet.html", gin.H{"error": "Failed to find address,choose different id"})
		return
	}

	var order models.Order
	order.Address_ID = uint(adrid)

	var payment models.Payment
	database.DB.Last(&payment)
	var address models.Address
	err = database.DB.Where("user_id=? AND address_id=?", userid, order.Address_ID).Last(&address).Error
	if err != nil {
		c.HTML(400, "wallet.html", gin.H{"error": "Failed to find address,choose different id"})
		return
	}

	err = database.DB.Create(&models.Order{
		User_ID:      userid,
		Address_ID:   order.Address_ID,
		Total_Price:  totalprice,
		Payment_ID:   payment.Payment_ID,
		Status:       "Processing",
		Payment_Type: "Wallet",
		Date:         time.Now(),
	}).Error
	if err != nil {
		log.Println("failed to create order")
		c.HTML(500, "wallet.html", gin.H{"error": err.Error()})
		return
	}

	//creating Wallet
	database.DB.Create(&models.Payment{
		Payment_Type:   "Wallet",
		Total_Amount:   totalprice,
		Payment_Status: "Completed",
		User_ID:        userid,
		Date:           time.Now(),
	})

	var order1 models.Order
	database.DB.Last(&order1)

	err = database.DB.Create(&models.OrderItem{
		Order_ID:    order1.Order_ID,
		User_ID:     userid,
		Product_ID:  product.ID,
		Address_ID:  order.Address_ID,
		Brand:       product.Brand_Name,
		Category:    product.Category_Name,
		Quantity:    uint(qty),
		Price:       product.Price,
		Total_Price: totalprice,
		Discount:    0,
		Status:      "processing",
		Created_at:  time.Now(),
	}).Error

	if err != nil {
		log.Println("Failed to create Order Item", err)
		c.HTML(400, "singleWallet.html", gin.H{"error": err.Error()})
		return
	}

	// Deduct the total price from the user's wallet
	var userwallet models.User
	database.DB.Where("user_id=?", userid).Find(&userwallet)

	userwallet.Wallet -= int(totalprice)
	if err := database.DB.Save(&userwallet).Error; err != nil {
		log.Println("Failed to update user wallet: ", err)
		c.HTML(500, "singleWallet.html", gin.H{"error": "Failed to update user wallet"})
		return
	}

	//reducing the stock count in database

	err = database.DB.Model(&models.Product{}).Where("id=?", product.ID).Update("stock", product.Stock-uint(qty)).Error
	if err != nil {
		log.Println("failed to update stock in database : ", err)
	}

	c.Redirect(303, "/user/payment-success")
}
