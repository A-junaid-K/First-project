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
		log.Println("error @ 17")
		c.HTML(http.StatusUnauthorized, "adminProductlist.html", gin.H{
			"error": categoryName + " is exist",
		})
		return
	}
	database.DB.Create(&models.Category{Name: categoryName})
	// if result.Error != nil {
	// 	log.Println("error @ 25")
	// 	c.HTML(400, "adminProductlist.html", gin.H{
	// 		"error": "Failed to Add Category",
	// 	})
	// 	return
	// }
	c.HTML(202, "adminProductlist.html", gin.H{"message": "Successfully Added " + categoryName})
}

// List category
func Listcatagory(c *gin.Context) {
	var category models.Category
	result := database.DB.Raw("SELECT * FROM categorys").Scan(&category)
	if result.Error != nil {
		c.HTML(400, "adminProductlist.html", gin.H{
			"error": "Failed to list category",
		})
		return
	}

}

//Offers

