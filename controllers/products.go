package controllers

import (
	// "html/template"

	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Addproducts(c *gin.Context) {
	c.HTML(http.StatusOK, "addProduct.html", nil)
}

func PostAddproducts(c *gin.Context) {
	var err error

	name := c.Request.FormValue("productName")
	description := c.Request.FormValue("productDescription")
	stock, _ := strconv.Atoi(c.Request.FormValue("productStock"))
	price, _ := strconv.Atoi(c.Request.FormValue("productPrice"))
	category_name := c.Request.FormValue("categoryName")
	brand_name := c.Request.FormValue("brandName")

	//Get the image file
	file, err := c.FormFile("productImage")
	if err != nil {
		c.HTML(http.StatusBadRequest, "addProduct.html", gin.H{
			"error": "Failed to upload image",
		})
		return
	}
	//Save the image file
	err = c.SaveUploadedFile(file, "./static/images/"+file.Filename)
	if err != nil {
		c.HTML(http.StatusBadRequest, "addProduct.html", gin.H{
			"error": "Failed to save image",
		})
		return
	}
	if err != nil {
		c.HTML(http.StatusBadRequest, "addProduct.html", gin.H{
			"error": "Failed to add product",
		})
		return
	}
	var dtproduct models.Product
	database.DB.Where("name=?", name).First(&dtproduct)

	if name == dtproduct.Name {
		c.HTML(http.StatusBadRequest, "addProduct.html", gin.H{
			"error": "This product already exist",
		})
		return
	}
	//----------------

	//adding category
	var dtcategory models.Category
	database.DB.Table("categories").Where("name=?", category_name).Scan(&dtcategory)
	if dtcategory.Name != category_name {
		log.Println("category is not exist")
		c.HTML(http.StatusBadRequest, "addProduct.html", gin.H{
			"error": "This category is not exist",
		})
		return
	}
	//----------------

	result := database.DB.Create(&models.Product{
		Name:          name,
		Description:   description,
		Stock:         uint(stock),
		Price:         uint(price),
		Category_Name: category_name,
		Brand_Name:    brand_name,
		Image:         file.Filename,
	})
	if result.Error != nil {
		log.Println("Failed to add product", err)
		c.HTML(http.StatusBadRequest, "addProduct.html", gin.H{
			"error": "Failed to add product",
		})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin-products-list")
}

func AdminListproducts(c *gin.Context) {
	data := DtTables()
	c.HTML(200, "adminProductlist.html", data)
}

func Listproducts(c *gin.Context) {
	var exp []models.Category_Offer
	database.DB.Find(&exp)
	currentdate := time.Now()

	for _, offer := range exp {
		if currentdate.After(offer.Expiry_date) {
			log.Println("offer has expired")
		}
	}

	data := DtTables()

	c.HTML(200, "productsList2.html", data)
}

func ProductDetails(c *gin.Context) {
	productiD := c.Param("id")

	var products models.Product
	database.DB.Table("products").Where("id=?", productiD).First(&products)

	//checking the offer
	if products.Offer_Name != "" {
		offerprice := products.Price * products.Percentage / 100
		products.Price -= offerprice
	}

	c.HTML(http.StatusOK, "productDetails2.html", products)
}

func Editproduct(c *gin.Context) {
	iD, _ := strconv.Atoi(c.Param("id"))
	var editproduct models.Product
	database.DB.Where("id=?", iD).First(&editproduct)
	c.HTML(200, "editproduct.html", editproduct)
}
func PostEditproduct(c *gin.Context) {
	var err error

	name := c.Request.FormValue("productName")
	description := c.Request.FormValue("productDescription")
	stock, _ := strconv.Atoi(c.Request.FormValue("productStock"))
	price, _ := strconv.Atoi(c.Request.FormValue("productPrice"))
	category_name := c.Request.FormValue("categoryName")
	brand_name := c.Request.FormValue("brandName")

	//Get the image file
	file, err := c.FormFile("productImage")
	if err != nil {
		// c.HTML(http.StatusBadRequest, "editproduct.html", gin.H{
		// 	"error": "Failed to upload image",
		// })
		// return
		file = nil
	}
	if file != nil {
		//Save the image file
		err = c.SaveUploadedFile(file, "./static/images/"+file.Filename)
		if err != nil {
			c.HTML(http.StatusBadRequest, "editproduct.html", gin.H{
				"error": "Failed to save the edited image",
			})
			return
		}
	}
	//checking category
	var dtcategory models.Category
	database.DB.Table("categories").Where("name=?", category_name).Scan(&dtcategory)
	if dtcategory.Name != category_name {
		log.Println("category is not exist")
		c.HTML(http.StatusBadRequest, "editproduct.html", gin.H{
			"error": "This category is not exist",
		})
		return
	}

	iD, _ := strconv.Atoi(c.Param("id"))
	var result *gorm.DB
	if file != nil {
		result = database.DB.Model(&models.Product{}).Where("id=?", iD).Updates(map[string]interface{}{
			"name":          name,
			"description":   description,
			"stock":         stock,
			"price":         price,
			"category_name": category_name,
			"brand_name":    brand_name,
			"image":         file.Filename,
		})
	} else {
		result = database.DB.Model(&models.Product{}).Where("id=?", iD).Updates(map[string]interface{}{
			"name":          name,
			"description":   description,
			"stock":         stock,
			"price":         price,
			"category_name": category_name,
			"brand_name":    brand_name,
		})
	}
	if result.Error != nil {
		log.Println("Failed to edit product")
		c.HTML(http.StatusBadRequest, "editproduct.html", gin.H{
			"error": "Failed to edit product",
		})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin-products-list")
}
func Deleteproduct(c *gin.Context) {
	idstr := c.Param("prdctid")
	id, _ := strconv.Atoi(idstr)

	res := database.DB.Where("id", id).Delete(&models.Product{})
	if res.RowsAffected == 0 {
		fmt.Println("Failed to find the product")
		c.HTML(400, "adminProductlist.html", gin.H{"error": "Failed to find product"})
		return
	}
	c.Redirect(303, "/admin-products-list")
}
func DtTables() interface{} {

	var products []models.Product
	var categories []models.Category
	var carts []models.Cart
	var addresses []models.Address
	var users []models.User
	var brands []models.Brand
	// var offers []models.Category_Offer

	database.DB.Find(&products)
	database.DB.Find(&categories)
	database.DB.Find(&carts)
	database.DB.Find(&addresses)
	database.DB.Find(&users)
	database.DB.Find(&brands)
	// database.DB.Find(&offers)

	data := struct {
		Products   []models.Product
		Categories []models.Category
		Carts      []models.Cart
		Addresses  []models.Address
		Users      []models.User
		Brands     []models.Brand
	}{
		Products:   products,
		Categories: categories,
		Carts:      carts,
		Addresses:  addresses,
		Users:      users,
		Brands:     brands,
	}

	return data
}
