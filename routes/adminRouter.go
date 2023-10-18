package routes

import (
	"github.com/first_project/controllers"
	"github.com/first_project/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {

	// router := r.Group("/admin")
	// {
	r.GET("/admin-login", controllers.Adminlogin)
	r.POST("/admin-login", controllers.PostAdminlogin)

	// user
	r.GET("/users-list", controllers.Listusers)
	r.GET("/block/:user_id", controllers.Blockuser)
	r.GET("/unblock/:user_id", controllers.Unblockuser)

	// products
	r.GET("/add-product", controllers.Addproducts)
	r.POST("/add-product", controllers.PostAddproducts)
	r.GET("/admin-products-list", middleware.AdminAuthentication, controllers.AdminListproducts)
	r.GET("/edit-product/:id", controllers.Editproduct)
	r.POST("/edit-product/:id", controllers.PostEditproduct)

	// Brand
	r.POST("/admin-products-list-addbrand", controllers.AddBrand)
	r.GET("/products-list", controllers.ListBrand)

	// Add Category
	r.POST("/admin-products-list-addcategory", controllers.AddCategory)

	// Category Offers

	// Coupen
	//add coupoun
	r.POST("/addcoupon/:coupon_code", middleware.AdminAuthentication, controllers.AddCoupon)
	// r.GET("/listcoupons", middleware.AdminAuthentication, controllers.ListCoupons)
	// r.PUT("/cancelcoupon/:coupon_id", middleware.AdminAuthentication, controllers.CancelCoupon)
	

	// r.POST("/userslist", controllers.Postlistusers)
	// r.Static("/assets", "./assets")
	// r.StaticFS("/more_static", http.Dir("my_file_system"))
	// r.StaticFile("/favicon.ico", "./resources/favicon.ico")

}
