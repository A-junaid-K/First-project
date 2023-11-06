package routes

import (
	"github.com/first_project/controllers"
	"github.com/first_project/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {

	r.GET("/admin-login", controllers.Adminlogin)
	r.POST("/admin-login", controllers.PostAdminlogin)

	// user
	r.GET("/users-list", middleware.AdminAuthentication, controllers.Listusers)
	r.GET("/block/:user_id", middleware.AdminAuthentication, controllers.Blockuser)
	r.GET("/unblock/:user_id", middleware.AdminAuthentication, controllers.Unblockuser)

	// products
	r.GET("/add-product", middleware.AdminAuthentication, controllers.Addproducts)
	r.POST("/add-product", middleware.AdminAuthentication, controllers.PostAddproducts)
	r.GET("/admin-products-list", middleware.AdminAuthentication, controllers.AdminListproducts)
	r.GET("/edit-product/:id", middleware.AdminAuthentication, controllers.Editproduct)
	r.POST("/edit-product/:id", middleware.AdminAuthentication, controllers.PostEditproduct)
	r.GET("/delete-product/:prdctid", middleware.AdminAuthentication, controllers.Deleteproduct)

	// Brand
	r.GET("/admin-brands", middleware.AdminAuthentication, controllers.DisplayBrands)
	r.POST("/admin-addbrand", middleware.AdminAuthentication, controllers.AddBrand)

	// Add Category
	r.GET("/admin-category", middleware.AdminAuthentication, controllers.DisplayCategory)
	r.POST("/admin-addcategory", middleware.AdminAuthentication, controllers.AddCategory)
	r.GET("/admin-listcategory/:category_id", middleware.AdminAuthentication, controllers.UnlistCategory)
	r.GET("/admin-unlistcategory/:category_id", middleware.AdminAuthentication, controllers.ListCategory)

	//add coupoun
	r.GET("/admin-coupon", middleware.AdminAuthentication, controllers.Coupon)
	r.POST("/addcoupon/:coupon_code", middleware.AdminAuthentication, controllers.PostAddCoupon)

	//offer
	r.GET("/admin-offer", middleware.AdminAuthentication, controllers.Offer)
	r.POST("/admin/add-offer", middleware.AdminAuthentication, controllers.PostAddOffer)
	r.GET("/admin/remove-offer", middleware.AdminAuthentication, controllers.RemoveOffer)

	r.GET("/generate-sales-report", middleware.AdminAuthentication, controllers.Getsalesreport)
	r.POST("/generate-sales-reports", middleware.AdminAuthentication, controllers.SalesReport)

}
