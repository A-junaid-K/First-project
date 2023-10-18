package controllers

import (
	"log"
	"net/http"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)


func AddBrand(c *gin.Context) {
	brandName := c.Request.FormValue("brand_name")
	var dtbrand models.Brand
	database.DB.Where("brand_name=?", brandName).First(&dtbrand)
	if brandName == dtbrand.Brand_Name {
		log.Println("error @ add brand")
		c.HTML(http.StatusBadRequest, "adminProductlist.html", gin.H{
			"error": brandName + " already exsit",
		})
		return
	}
	database.DB.Create(&models.Brand{Brand_Name: brandName})
	c.HTML(http.StatusAccepted, "adminProductlist.html", gin.H{
		"message": "Successfully Added " + brandName,
	})
}
func ListBrand(c *gin.Context){
	var Brand []models.Brand
	result := database.DB.Raw("select * from brands").Scan(&Brand)
	if result.Error != nil {
		log.Println("error @ list brand")
		c.HTML(http.StatusBadRequest, "adminProductlist.html", gin.H{
			"error": "Failed to list brand",
		})
		return
	}

	c.HTML(http.StatusAccepted, "adminProductlist.html", Brand)
}
