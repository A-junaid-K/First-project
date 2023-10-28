package controllers

import (
	"log"
	"net/http"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

// Add category
func AddCategory(c *gin.Context) {
	categoryName := c.Request.FormValue("category_name")
	var ctr models.Category
	database.DB.Where("category_name=?", categoryName).First(&ctr)
	if ctr.Name == categoryName {
		log.Println("the category exist")
		c.HTML(http.StatusUnauthorized, "adminProductlist.html", gin.H{
			"error": categoryName + " is exist",
		})
		return
	}
	err := database.DB.Create(&models.Category{Name: categoryName}).Error
	if err != nil {
		log.Println("Failed to add category")
		return
	}
	c.Redirect(303, "/admin-products-list")
}
func AddBrand(c *gin.Context) {
	brandName := c.Request.FormValue("brand_name")
	var dtbrand models.Brand
	database.DB.Where("brand_name=?", brandName).First(&dtbrand)
	if brandName == dtbrand.Brand_Name {
		log.Println("the brand still exist")
		c.HTML(http.StatusBadRequest, "adminProductlist.html", gin.H{
			"error": brandName + " already exsit",
		})
		return
	}
	err := database.DB.Create(&models.Brand{Brand_Name: brandName}).Error
	if err != nil {
		log.Println("Failed to add brand")
		return
	}
	c.Redirect(303, "/admin-products-list")
}

// func ListBrand(c *gin.Context) {
// 	var Brand []models.Brand
// 	result := database.DB.Raw("select * from brands").Scan(&Brand)
// 	if result.Error != nil {
// 		log.Println("error @ list brand")
// 		c.HTML(http.StatusBadRequest, "adminProductlist.html", gin.H{
// 			"error": "Failed to list brand",
// 		})
// 		return
// 	}

// 	c.HTML(http.StatusAccepted, "adminProductlist.html", Brand)
// }
