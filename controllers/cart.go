package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func AddtoCart(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_id

	//get the product id
	product_id, _ := strconv.Atoi(c.Param("id"))

	// Get the size & quantity
	pdsize := c.PostForm("size")
	log.Println("pdsize : ", pdsize)
	pdquantity, _ := strconv.Atoi(c.PostForm("quantity"))
	log.Println("pdquantity : ", pdquantity)

	var product models.Product

	err := database.DB.First(&product, product_id).Error
	if err != nil {
		fmt.Println(product_id)
		fmt.Println("product not found ", err)
		c.HTML(404, "cart2.html", gin.H{"error": "Product not found"})
		return
	}

	totalprice := product.Price * uint(pdquantity)

	//creating cart
	var dtcart models.Cart
	err = database.DB.Where("product_id=? AND user_id=?", product_id, userId).First(&dtcart).Error
	if err != nil {
		err = database.DB.Create(&models.Cart{
			Product_ID:     product_id,
			Name:           product.Name,
			Description:    product.Description,
			Quantity:       pdquantity,
			Stock:          int(product.Stock),
			Price:          int(product.Price),
			Size:           pdsize,
			Total_Price:    totalprice,
			Category_Name:  product.Category_Name,
			Brand_Name:     product.Brand_Name,
			Category_Offer: product.Percentage,
			User_ID:        userId,
			Image:          product.Image,
		}).Error
	} else {
		fmt.Println("already exist in cart")
		c.Redirect(303, "/user/cart")
		return
	}
	if err != nil {
		c.HTML(400, "cart2.html", gin.H{
			"error": "Failed to fetch cart database",
		})
		return
	}
	fmt.Println("added to cart. \nError : ", err)

	//searching for database all cart data
	var cartdata []models.Cart
	err = database.DB.Where("user_id=?", userId).Find(&cartdata).Error
	if err != nil {
		c.HTML(400, "cod.html", gin.H{"error": "Please check your cart"})
		return
	}

	//checking stock level
	var level int

	database.DB.First(&product, product_id)
	level = int(product.Stock) - pdquantity

	if level < 0 {
		log.Println("error : please check quantity : ", err)
		log.Println("stock : ", product.Stock)
		log.Println("quantity : ", pdquantity)
		log.Println("levell : ", level)
		err = database.DB.Where("product_id=?", product.ID).Delete(&models.Cart{}).Error
		if err != nil {
			log.Println("checking stock level in cart : ", err)
		}
	}

	// Redirection
	if level >= 0 {
		c.Redirect(303, "/user/cart")
	} else {
		c.Redirect(303, fmt.Sprintf("/user/product-details/%v", product_id))

	}

}
func ListCart(c *gin.Context) {
	user, _ := c.Get("user")
	user_id := user.(models.User).User_id
	var usercart []models.Cart
	database.DB.Where("user_id=?", user_id).Find(&usercart)

	c.HTML(200, "cart2.html", usercart)
}
func RemoveFromCart(c *gin.Context) {
	productId, _ := strconv.Atoi(c.Param("productid"))
	user, _ := c.Get("user")
	user_id := user.(models.User).User_id
	database.DB.Where("user_id=? AND product_id=?", user_id, productId).Delete(&models.Cart{})
	fmt.Println("removed from cart")
	c.Redirect(303, "/user/cart")
}
