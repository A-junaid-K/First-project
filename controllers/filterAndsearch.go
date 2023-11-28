package controllers

import (
	"log"
	"strings"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func FilterCategory(c *gin.Context) {
	filtered_category := c.Query("category_name")
	var filterproduct []models.Product
	database.DB.Table("products").Where("category_name=?", filtered_category).Find(&filterproduct)
	if filterproduct == nil {
		log.Println("error : No products in this catagory")
		c.HTML(400, "productsList2.html", gin.H{
			"error": "No products in this catagory",
		})
		return
	}

	var categories []models.Category
	var brands []models.Brand

	database.DB.Find(&categories)
	database.DB.Find(&brands)

	data := struct {
		Products   []models.Product
		Categories []models.Category
		Brands     []models.Brand
	}{
		Products:   filterproduct,
		Categories: categories,
		Brands:     brands,
	}

	c.HTML(200, "productsList2.html", data)
}

func FilterBrand(c *gin.Context) {
	filtered_brand := c.Query("brand_name")
	var filterproduct []models.Product
	database.DB.Table("products").Where("brand_name=?", filtered_brand).Find(&filterproduct)
	if filterproduct == nil {
		log.Println("error : No products in this brand")
		c.HTML(400, "productsList2.html", gin.H{
			"error": "No products in this brand",
		})
		return
	}

	var categories []models.Category
	var brands []models.Brand

	database.DB.Find(&categories)
	database.DB.Find(&brands)

	data := struct {
		Products   []models.Product
		Categories []models.Category
		Brands     []models.Brand
	}{
		Products:   filterproduct,
		Categories: categories,
		Brands:     brands,
	}

	c.HTML(200, "productsList2.html", data)
}

func MultipleFilter(c *gin.Context) {

}

func Search(c *gin.Context) {

	name := c.Query("name")

	log.Println("name : ", name)

	var products []models.Product

	err := database.DB.Find(&products).Error
	if err != nil {
		log.Println("failed to search product : ", err)
		c.HTML(400, "productList2.html", gin.H{
			"error": "No values in database",
		})
		return
	}

	for _, v := range products {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(name)) {
			products = append(products, v)
			log.Println("containing")
		} else {
			log.Println("no product there")
		}
	}

	var categories []models.Category
	var brands []models.Brand

	database.DB.Find(&categories)
	database.DB.Find(&brands)

	data := struct {
		Products   []models.Product
		Categories []models.Category
		Brands     []models.Brand
	}{
		Products:   products,
		Categories: categories,
		Brands:     brands,
	}

	c.HTML(200, "productsList2.html", data)
}
