package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func Coupon(c *gin.Context) {
	var coupon []models.Coupon
	database.DB.Find(&coupon)
	c.HTML(200, "coupon.html", coupon)
}
func PostAddCoupon(c *gin.Context) {
	layout := "2006-02-01"

	coupon_code := c.Request.FormValue("couponcode")
	coupon_type := c.PostForm("type")

	raw_expiry_date := c.PostForm("expiryDate")
	expiry_date, _ := time.Parse(layout, raw_expiry_date)
	discount, _ := strconv.Atoi(c.Request.FormValue("discount"))

	var coupon1 models.Coupon
	database.DB.Where("coupon_code=?", coupon_code).First(&coupon1)

	if coupon_code == coupon1.Coupon_Code {
		log.Println("This coupon code already exist in database")
		c.HTML(400, "coupon.html", gin.H{
			"error": "This coupon code already exist in database",
		})
		return
	}

	if len(coupon_code) < 5 || len(coupon_code) > 10 {
		log.Println("Coupon code must be lenght between 5 to 10")
		c.HTML(400, "coupon.html", gin.H{
			"error": "Coupon code must be lenght between 5 to 10",
		})
		return
	}

	database.DB.Create(&models.Coupon{
		Coupon_Code:   coupon_code,
		Starting_Time: time.Now(),
		Ending_Time:   time.Now().Add(time.Hour * 24 * time.Duration(expiry_date.Day())),
		Type:          coupon_type,
		Value:         uint(discount),
		// Max_Discount:  coupon.Max_Discount,
		// Min_Discount:  coupon.Min_Discount,
	})
	c.Redirect(303, "/admin-coupon")
}
func CancelCoupon(c *gin.Context) {
	cid := c.Query("coupon_id")
	log.Println("COUPON ID		: ", cid)
	var cou models.Coupon
	err := database.DB.First(&cou, cid).Error

	if err != nil {
		log.Println("error : Failed to find coupon please try different id")
		return
	}

	err = database.DB.Model(&models.Coupon{}).Where("coupon_id=?", cid).Updates(map[string]interface{}{"cancel": true, "ending_time": time.Now()}).Error
	if err != nil {
		log.Println("failed to cancel coupon : ", err)
		return
	}

	c.Redirect(303, "/admin-coupon")
}
func ApproveCoupon(c *gin.Context) {
	cid := c.Query("coupon_id")

	var cou models.Coupon
	err := database.DB.First(&cou, cid).Error

	if err != nil {
		log.Println("error : Failed to find coupon please try different id")
		return
	}

	err = database.DB.Model(&models.Coupon{}).Where("coupon_id=?", cid).Updates(map[string]interface{}{
		"cancel":      false,
		"ending_time": time.Now().Add(time.Hour * 24 * time.Duration(cou.Ending_Time.Day())),
	}).Error
	if err != nil {
		log.Println("failed to approve coupon : ", err)
		return
	}

	c.Redirect(303, "/admin-coupon")
}
func RemoveCoupon(c *gin.Context) {
	remove_coupon := c.Query("coupon_id")
	err := database.DB.Where("coupon_id=?", remove_coupon).Delete(&models.Coupon{}).Error
	if err != nil {
		log.Println("Failed to remove coupon : ", err)
		return
	}
	c.Redirect(303, "/admin-coupon")
}

func ApplyCoupon(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_id

	coupon_code := c.Request.FormValue("coupon_code")
	log.Println("coupon  : ", coupon_code)
	if coupon_code == "" {
		c.Next()
		return
	}

	//find the coupon
	var coupon1 models.Coupon
	row := database.DB.Where("coupon_code=?", coupon_code).First(&coupon1).RowsAffected
	if row == 0 {
		log.Println("Failed to find coupon")
		c.HTML(400, "checkout.html", gin.H{"error": "Failed to find coupon"})
		return
	}

	//checking coupon expired or not
	if time.Now().Unix() > (coupon1.Ending_Time).Unix() {
		log.Println("Coupon expired")
		c.HTML(400, "checkout.html", gin.H{"error": "Coupon expired"})
		return
	}

	//checking the coupon already applied or not
	var cart []models.Cart
	row = database.DB.Where("user_id=? AND coupon_applied = true", userId).Find(&cart).RowsAffected
	if row >= 1 {
		log.Println("coupon already applied")
		c.HTML(400, "checkout.html", gin.H{"error": "coupon already applied"})
		return
	}

	//getting the cart data
	var cart1 []models.Cart
	err := database.DB.Where("user_id=?", userId).Find(&cart1).Error
	if err != nil {
		log.Println("cart is empty")
		return
	}

	//checking coupon valid or not
	if coupon1.Cancel {
		log.Println("Coupon is not valid")
		c.HTML(400, "checkout.html", gin.H{"error": "Coupon is not valid"})
		return
	}

	if coupon1.Type == "Percentage" {

		for _, v := range cart1 {
			discount := (v.Total_Price * coupon1.Value / 100)
			err := database.DB.Model(&models.Cart{}).Where("user_id=? AND id=?", userId, v.ID).Updates(map[string]interface{}{"total_price": v.Total_Price - discount, "coupon_discount": discount, "coupon_applied": true}).Error
			if err != nil {
				log.Println(err)
				return
			}
		}
		log.Println("Coupon applied successfully")
		c.Next()
	} else {
		var cartitems int64
		database.DB.Model(&models.Cart{}).Where("user_id=?", userId).Count(&cartitems)
		log.Println("cart lentght  ", len(cart1))
		for _, v := range cart1 {
			discount := coupon1.Value / uint(cartitems)
			err := database.DB.Model(&models.Cart{}).Where("user_id=? AND id=?", userId, v.ID).Updates(map[string]interface{}{"total_price": v.Total_Price - discount, "coupon_discount": discount, "coupon_applied": true}).Error
			if err != nil {
				log.Println(err)
				return
			}
		}
		log.Println("coupon applied successfully")
		c.Next()
	}
}

//-------------------------------------------------OFFER-----------------------------------------------//

func Offer(c *gin.Context) {
	dtcategory := DtTables()
	c.HTML(200, "offer.html", dtcategory)
}
func PostAddOffer(c *gin.Context) {
	layout := "2006-01-02"
	offer_name := c.Request.FormValue("offer_name")

	raw_starting_date := c.PostForm("startingDate")
	starting_date, err := time.Parse(layout, raw_starting_date)
	if err != nil {
		log.Println("failed to get starting date : ", err)
	}
	raw_expiry_date := c.PostForm("expiryDate")
	expiry_date, errr := time.Parse(layout, raw_expiry_date)
	if errr != nil {
		log.Println("failed to get starting date : ", err)
	}

	log.Println("starting time : ", starting_date, "\n ending time : ", expiry_date)

	percentage, err := strconv.Atoi(c.Request.FormValue("percentage"))

	if err != nil {
		log.Println("failed form value : ", err)
	}

	// recieving the chosen category
	category := c.PostForm("chosencategory")

	var dtoffer models.Category_Offer
	err = database.DB.Where("offer_name=?", offer_name).First(&dtoffer).Error

	if err != nil {
		log.Println("database err  : ", err)
	}

	log.Println("dtoffer : ", dtoffer)
	if offer_name == dtoffer.Offer_Name {
		log.Println("This offer already exist : ")
		c.HTML(http.StatusBadRequest, "offer.html", gin.H{
			"error": "This offer already exist",
		})
		return
	}

	// Adding offer
	result := database.DB.Create(&models.Category_Offer{
		Offer_Name:    offer_name,
		Category_Name: category,
		Starting_Time: starting_date,
		Expiry_date:   expiry_date,
		Percentage:    uint(percentage),
		Offer:         true,
	})
	if result.Error != nil {
		log.Println("Failed to add offer : ", result.Error)
		c.HTML(http.StatusBadRequest, "offer.html", gin.H{
			"error": "Failed to add offer",
		})
		return
	}

	database.DB.Table("products").Where("category_name=?", category).Updates(map[string]interface{}{
		"offer_name": offer_name,
		"percentage": percentage,
	})
	err = database.DB.Table("categories").Where("name=?", category).Updates(map[string]interface{}{
		"offer_name": offer_name,
		"percentage": percentage,
	}).Error
	if err != nil {
		log.Println("errrrrrr : ", err)
	}

	err = database.DB.Table("carts").Where("category_name=?", category).Update("category_offer", percentage).Error
	if err != nil {
		log.Println("errrrrrr : ", err)
	}

	c.Redirect(http.StatusSeeOther, "/admin-offer")
}
func RemoveOffer(c *gin.Context) {
	db := database.DB
	delete_offer := c.Query("offer_name")
	// Find and update the rows where 'offer_name' is same
	result := db.Model(&models.Category{}).Where("offer_name = ?", delete_offer).Updates(map[string]interface{}{
		"offer_name": nil,
		"percentage": 0,
	})
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to delete data"})
		return
	}

	// Find and update the rows where 'offer_name' is same
	result2 := db.Model(&models.Product{}).Where("offer_name = ?", delete_offer).Updates(map[string]interface{}{
		"offer_name": nil,
		"percentage": 0,
	})
	if result2.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to delete data"})
		return
	}
	db.Where("offer_name=?", delete_offer).Delete(&models.Category_Offer{})

	c.Redirect(303, "/admin-category")
}
