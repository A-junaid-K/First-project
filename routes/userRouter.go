package routes

import (
	"github.com/first_project/controllers"
	"github.com/first_project/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*.html")
	// r.LoadHTMLGlob("nest-backend/*.html")
	r.Static("/static", "./static")
	// r.Static("/nest-backend/assets", "./assets")
	r.GET("/", middleware.UserAuthentication, controllers.Home)
	router := r.Group("/user")
	{
		//   User
		router.GET("/signup", controllers.SignUp)
		router.POST("/signup", controllers.PostSignUp)
		router.GET("/varifyotp", controllers.VarifyOtp)
		router.POST("/varifyotp", controllers.PostVarifyOtp)
		router.GET("/login", controllers.Login)
		router.POST("/login", controllers.Postlogin)
		router.GET("/log-out", middleware.UserAuthentication, controllers.Logout)

		router.GET("/about", controllers.About)
		router.GET("/gallery", controllers.Gallery)
		router.GET("/testimonial", controllers.Testimonial)
		router.GET("/contac", controllers.Contact)
		router.GET("/news", controllers.News)

		// User details
		router.GET("/user-details", middleware.UserAuthentication, controllers.ListUserDetails)

		//edit address
		router.GET("/add-address", middleware.UserAuthentication, controllers.AddAddress)
		router.POST("/add-address", middleware.UserAuthentication, controllers.PostAddAddress)
		router.GET("/edit-address/:adrid", middleware.UserAuthentication, controllers.EditAddress)
		router.POST("/edit-address/:adrid", middleware.UserAuthentication, controllers.PostEditAddress)

		//edit profile
		router.GET("/edit-profile", middleware.UserAuthentication, controllers.Editprofile)
		router.POST("/edit-profile", middleware.UserAuthentication, controllers.PostEditprofile)

		//Product //category //brand
		router.GET("/products-list", controllers.Listproducts)
		router.GET("/product-details/:id", controllers.ProductDetails)
		router.POST("/product-details", controllers.ProductDetails)

		//Cart
		router.GET("/cart/:id", middleware.UserAuthentication, controllers.AddtoCart)
		router.GET("/cart", middleware.UserAuthentication, controllers.ListCart)
		router.GET("/remove-from-cart/:productid", middleware.UserAuthentication, controllers.RemoveFromCart)

		//Wishlist
		router.GET("/wishlist/:id", middleware.UserAuthentication, controllers.AddToWishlist)
		router.GET("/wishlist", middleware.UserAuthentication, controllers.ListWishlist)

		// Checkout
		router.GET("/checkout", middleware.UserAuthentication, controllers.Checkout)
		router.POST("/checkout", middleware.UserAuthentication, controllers.PostCheckout)

		// Payment
		router.GET("/payment-cod", middleware.UserAuthentication, controllers.Cod)
		// router.GET("/payment-razorpay", middleware.UserAuthentication, controllers.RazorPay)
		router.GET("/payment-razorpay", middleware.UserAuthentication, controllers.RazorPay)
		router.GET("/payment-razorpay-success", middleware.UserAuthentication, controllers.RazorpaySuccess)
		router.GET("/payment-success", middleware.UserAuthentication, controllers.Success)
	}

}
