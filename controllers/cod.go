package controllers

import (
	"fmt"
	"time"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func Cod(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	//searching for database all cart data
	var cartdata []models.Cart
	err := database.DB.Where("user_id=?", userid).Find(&cartdata).Error
	if err != nil {
		c.HTML(400, "cod.html", gin.H{"error": "Please check your cart"})
		return
	}

	//getting total price of cart
	var totalprice uint
	err = database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", userid).Scan(&totalprice).Error
	if err != nil {
		if err != nil {
			c.HTML(400, "cod.html", gin.H{"error": "Failed to find total price", "message": "cart is empty"})
			return
		}
		return
	}

	var product models.Product

	// //checking stock level
	// for _, v := range cartdata {
	// 	database.DB.First(&product, v.Product_ID)
	// 	if product.Stock-v.Quantity < 0 {
	// 		c.HTML(400, "cod.html", gin.H{
	// 			"error": "Please check quantity",
	// 		})
	// 		return
	// 	}
	// }

	//creating COD(not completed)
	database.DB.Create(&models.Payment{
		Payment_Type:   "COD",
		Total_Amount:   totalprice,
		Payment_Status: "Pending",
		User_ID:        userid,
		Date:           time.Now(),
	})

	var adrid int
	// err = database.DB.Table("contactdetails").Where("user_id=?", userid).Last(&order).Error
	err = database.DB.Model(&models.Contactdetails{}).Select("address_id").Where("user_id=?", userid).Scan(&adrid).Error
	if err != nil {
		fmt.Println("failed to fetch address id from checkout page")
		c.HTML(400, "cod.html", gin.H{"error": "Failed to find address,choose different id"})
		return
	}

	var order models.Order
	order.Address_ID = uint(adrid)

	var payment models.Payment
	database.DB.Last(&payment)
	var address models.Address
	err = database.DB.Where("user_id=? AND address_id=?", userid, order.Address_ID).Last(&address).Error
	if err != nil {
		c.HTML(400, "cod.html", gin.H{"error": "Failed to find address,choose different id"})
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
		fmt.Println("failed to create order")
		c.HTML(400, "cod.html", gin.H{"error": err.Error()})
		return
	}

	var cartbrand struct {
		Brand_Name    string
		Category_Name string
	}

	var order1 models.Order
	database.DB.Last(&order1)

	//have to check it out

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
		fmt.Println(err)
		c.HTML(400, "cod.html", gin.H{"error": err.Error()})
		return
	}

	//reducing the stock count in database
	var products models.Product
	for _, v := range cartdata {
		database.DB.First(&products, v.Product_ID)
		database.DB.Model(&models.Product{}).Where("id=?", v.Product_ID).Update("stock", product.Stock-v.Quantity)
	}

	//deleting the checked out cart
	err = database.DB.Delete(&models.Cart{}, "user_id=?", userid).Error
	if err != nil {
		fmt.Println("Failed to delete checked out cart data")
		c.HTML(400, "cod.html", gin.H{"error": "failed to delete used cart" + err.Error()})
		return
	}

	c.Redirect(303, "/user/payment-success")

	//-----------------------------success------------------------------------------

	// pid, err := strconv.Atoi(c.Query("id"))
	// if err != nil {
	// 	c.HTML(400, "success.html", gin.H{
	// 		"error": "Error in string conversion",
	// 	})
	// 	return
	// }

	// c.HTML(200, "success.html", gin.H{
	// 	"paymentid": pid,
	// })

}
