package controllers

import (
	"fmt"
	"strconv"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func AddToWishlist(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_id

	product_id, _ := strconv.Atoi(c.Param("id"))

	var product models.Product

	err := database.DB.First(&product, product_id).Error
	if err != nil {
		fmt.Println("product not found ", err)
		c.HTML(404, "wishlist.html", gin.H{"error": "Product not found"})
		return
	}

	//creating cart
	var dtwishlist models.Cart
	err = database.DB.Where("product_id=? AND user_id=?", product_id, userId).First(&dtwishlist).Error
	if err != nil {
		err = database.DB.Create(&models.Wishlist{
			Product_ID:  product_id,
			Name:        product.Name,
			Description: product.Description,
			Quantity:    1,
			Stock:       product.Stock,
			Price:       product.Price,
			// Size:          size,
			Total_Price:   uint(product.Price),
			Category_Name: product.Category_Name,
			Brand_Name:    product.Brand_Name,
			User_ID:       userId,
			Image:         product.Image,
		}).Error
	} else {
		fmt.Println("already exist")
		c.Redirect(303, "/user/wishlist")
		return
	}
	if err != nil {
		fmt.Println("Failed to fetch wishlist database")
		c.HTML(400, "wishlist.html", gin.H{"error": "Failed to fetch wishlist database"})
		return
	}
	fmt.Println("added to cart. \nError : ", err)

	// c.HTML(200, "cart2.html", dtcart)
	c.Redirect(303, "/user/wishlist")
}
func ListWishlist(c *gin.Context) {
	user, _ := c.Get("user")
	user_id := user.(models.User).User_id
	var userwishlist []models.Wishlist
	database.DB.Where("user_id=?", user_id).Find(&userwishlist)

	c.HTML(200, "wishlist.html", userwishlist)
}
