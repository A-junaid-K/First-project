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

	//creating cart
	var dtcart models.Cart
	err = database.DB.Where("product_id=? AND user_id=?", product_id, userId).First(&dtcart).Error
	if err != nil {
		err = database.DB.Create(&models.Cart{
			Product_ID:    product_id,
			Name:          product.Name,
			Description:   product.Description,
			Quantity:      uint(pdquantity),
			Stock:         int(product.Stock),
			Price:         int(product.Price),
			Size:          pdsize,
			Total_Price:   uint(product.Price),
			Category_Name: product.Category_Name,
			Brand_Name:    product.Brand_Name,
			User_ID:       userId,
			Image:         product.Image,
		}).Error
	} else {
		fmt.Println("already exist in cart")
		c.Redirect(303, "/user/cart")
		// c.HTML(http.StatusBadRequest, "cart2.html", gin.H{
		// 	"error": "the product already exist in cart",
		// })
		return
	}
	if err != nil {
		fmt.Println("error @ 51")
		c.HTML(400, "cart2.html", gin.H{
			"error": "Failed to fetch cart database",
		})
		return
	}
	fmt.Println("added to cart. \nError : ", err)

	// c.HTML(200, "cart2.html", dtcart)
	c.Redirect(303, "/user/cart")
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
