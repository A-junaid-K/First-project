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
	log.Println("category id  = ", id)
	var category models.Category
	result := database.DB.First(&category, id)
	if result.Error != nil {
		log.Println(result.Error)
		c.HTML(http.StatusBadRequest, "category.html", gin.H{"error": result.Error})
		return
	}
	//Unlisting the Category
	result = database.DB.Model(&models.Category{}).Where("category_id=?", id).Update("unlist", true)
	if result.Error != nil {
		log.Println("Failed to Unlist Category")
		c.HTML(http.StatusBadRequest, "category.html", gin.H{"error": "Failed to Unlist Category"})
		return
	}
	log.Println("Unlisted in category")

	//Unlisting the Category in products
	result = database.DB.Model(&models.Product{}).Where("category_name=?", category.Name).Update("unlist", true)
	if result.Error != nil {
		log.Println("Failed to Unlist Category in products")
		c.HTML(http.StatusBadRequest, "category.html", gin.H{"error": "Failed to Unlist Category in products"})
		return
	}
	log.Println("unlisted in products")

	c.Redirect(http.StatusSeeOther, "/admin-category")
}

func ListCategory(c *gin.Context) {
	//get the id
	strid := c.Param("category_id")
	log.Println("category string id = ", strid)
	id, err := strconv.Atoi(strid)
	log.Println("category id = ", strid)
	if err != nil {
		log.Println("category list error : ", err)
		return
	}

	var category models.Category
	result := database.DB.First(&category, id)
	if result.Error != nil {
		log.Println("failed to fetch category : ", result.Error)
		c.HTML(http.StatusBadRequest, "category.html", gin.H{"error": result.Error})
		return
	}
	//listing the Category
	result = database.DB.Model(&models.Category{}).Where("category_id=?", id).Update("unlist", false)
	if result.Error != nil {
		log.Println("Failed to list Category")
		c.HTML(http.StatusBadRequest, "category.html", gin.H{"error": "Failed to list Category"})
		return
	}
	log.Println("listed in category")
	//listing the Category in products
	result = database.DB.Model(&models.Product{}).Where("category_name=?", category.Name).Update("unlist", false)
	if result.Error != nil {
		log.Println("Failed to list Category products")
		c.HTML(http.StatusBadRequest, "category.html", gin.H{"error": "Failed to list Category  products"})
		return
	}
	log.Println("listed in products")

	c.Redirect(http.StatusSeeOther, "/admin-category")
}

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

// -----------------------------------Brand------------------------------//
func DisplayBrands(c *gin.Context) {
	brands := DtTables()

	c.HTML(200, "brand.html", brands)
}
func AddBrand(c *gin.Context) {
	brandName := c.Request.FormValue("Brand_name")
	log.Println("brand name : ", brandName)
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
	c.Redirect(303, "/admin-brands")
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

func RemoveBrand(c *gin.Context) {

	//get the id
	id, _ := strconv.Atoi(c.Param("brand_id"))
	log.Println("brand  id  = ", id)
	var brand models.Brand
	result := database.DB.First(&brand, id)
	if result.Error != nil {
		log.Println(result.Error)
		c.HTML(http.StatusBadRequest, "brand.html", gin.H{"error": result.Error})
		return
	}
	//Removing the Brand
	result = database.DB.Where("brand_id=?", id).Delete(&models.Brand{})
	if result.Error != nil {
		log.Println("Failed to Remove brand")
		c.HTML(http.StatusBadRequest, "brand.html", gin.H{"error": "Failed to Remove brand"})
		return
	}
	log.Println("Removed in Brands")

	//Removing brands in products
	result = database.DB.Where("brand_name=?", brand.Brand_Name).Delete(&models.Product{})
	if result.Error != nil {
		log.Println("Failed to Remove brands in products")
		c.HTML(http.StatusBadRequest, "brand.html", gin.H{"error": "Failed to Remove brands in products"})
		return
	}
	log.Println("Removed brands in products")

	c.Redirect(http.StatusSeeOther, "/admin-brands")

}
