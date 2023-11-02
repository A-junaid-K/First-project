package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

// Category
func DisplayCategory(c *gin.Context) {
	category := DtTables()

	c.HTML(200, "category.html", category)
}
func AddCategory(c *gin.Context) {
	categoryName := c.Request.FormValue("category_name")
	var ctr models.Category
	database.DB.Where("category_name=?", categoryName).First(&ctr)
	if ctr.Name == categoryName {
		log.Println("the category exist")
		c.HTML(http.StatusUnauthorized, "category.html", gin.H{
			"error": categoryName + " is exist",
		})
		return
	}
	err := database.DB.Create(&models.Category{Name: categoryName}).Error
	if err != nil {
		log.Println("Failed to add category")
		return
	}
	c.Redirect(303, "/admin-category")
}
func UnlistCategory(c *gin.Context) {
	//get the id
	id, _ := strconv.Atoi(c.Param("category_id"))
	var category models.Category
	result := database.DB.First(&category, id)
	if result.Error != nil {
		log.Println(result.Error)
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{"error": result.Error})
		return
	}
	//Unlisting the Category
	result = database.DB.Model(&models.Category{}).Where("category_id=?", id).Update("unlist", true)
	if result.Error != nil {
		log.Println("Failed to Unlist Category")
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{"error": "Failed to Unlist Category"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/category")
}
func ListCategory(c *gin.Context) {
	//get the id
	id, _ := strconv.Atoi(c.Param("category_id"))
	var category models.Category
	result := database.DB.First(&category, id)
	if result.Error != nil {
		log.Println(result.Error)
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{"error": result.Error})
		return
	}
	//Unlisting the Category
	result = database.DB.Model(&models.Category{}).Where("category_id=?", id).Update("unlist", false)
	if result.Error != nil {
		log.Println("Failed to list Category")
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{"error": "Failed to list Category"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/category")
}
func FilterCategory(c *gin.Context) {
	// id,_ := strconv.Atoi(c.Param("category_id"))
	catagory := c.Query("catagory_name")

	var categorydb models.Category
	err := database.DB.Where("catagory_name=?", catagory).First(&categorydb).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "This catagory not exist",
		})
		return
	}
	type product struct {
		ID            uint
		Product_Name  string
		Price         uint
		Brand_Name    string
		Catagory_Name string
		Stock         uint
	}
	var products []product
	err = database.DB.Table("products").Select("products.id,products.product_name,products.price,brands.brand_name,catagories.catagory_name,products.stock").
		Joins("INNER JOIN brands ON brands.brand_id=products.brand_id").Joins("INNER JOIN catagories ON catagories.catagory_id=products.catagory_id").
		Where("products.catagory_id=?", categorydb.Category_id).Scan(&products).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "database error",
		})
		return
	}
	c.JSON(200, gin.H{
		"products": products,
	})
}

// -----------------------------------Brand------------------------------//
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
