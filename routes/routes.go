package routes

import (
	address "backend/controllers/address_controller"
	cart "backend/controllers/cart_controller"
	checkout "backend/controllers/checkout_controller"
	homepage "backend/controllers/homepage_controller"
	products "backend/controllers/products_controller"
	token "backend/controllers/token_controller"
	users "backend/controllers/users_controller"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {

	// USERS API
	auth := router.Group("/api/auth")
	{
		auth.POST("/signup", users.Register)
		auth.POST("/login", users.Login)
		auth.GET("/refresh", token.RefreshToken)
		auth.POST("/logout", middleware.AuthMiddleware(), users.SingleSessionLogout)
		auth.POST("/logout-all-devices", middleware.AuthMiddleware(), users.AllSessionLogout)
	}

	protected := router.Group("/api/users")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/all_users", users.AllUsers)
		protected.GET("/auth_user", users.AuthUser)
		protected.POST("/change_role", users.ChangeRole)
	}

	// PRODUCTS API
	productsapi := router.Group("/api/products")
	{
		productsapi.GET("/single_category_products", products.SingleCategoryProduct)
		productsapi.GET("/category_products/:category", products.GetCategoryProducts)
		productsapi.GET("/get_products", products.GetProducts)
		productsapi.POST("/create_product", middleware.AuthMiddleware(), products.UploadProduct)
		productsapi.PUT("/update_products", middleware.AuthMiddleware(), products.UpdateProduct)
	}

	// HOMEPAGE API
	homepageapi := router.Group("/api")
	{
		homepageapi.GET("/homepage", middleware.OptionalAuthMiddleware(), homepage.GetHomePage)
	}

	// CART ITEMS API
	cartapi := router.Group("/api")
	{
		cartapi.GET("/cart_products", middleware.AuthMiddleware(), cart.GetCartItems)
		cartapi.POST("/add_to_cart/:productId", middleware.AuthMiddleware(), cart.AddToCart)
		cartapi.PUT("/reduce_cart_item/:productId", middleware.AuthMiddleware(), cart.ReduceCartItem)
		cartapi.DELETE("/delete_product/:productId", middleware.AuthMiddleware(), cart.DeleteCartItem)
	}

	// ADDRESS API
	addressapi := router.Group("/api")
	addressapi.Use(middleware.AuthMiddleware())
	{
		addressapi.POST("/create_address", address.CreateAddress)
		addressapi.GET("/get_addresses", address.GetAddresses)
		addressapi.PATCH("/update_address/:addressId", address.UpdateAddress)
		addressapi.GET("/update_address/default/:addressId", address.UpdateDefaultAddress)
		addressapi.DELETE("/delete_address/:addressId", address.DeleteAddress)
	}

	// CHECKOUT API
	checkoutapi := router.Group("/api")
	checkoutapi.Use(middleware.AuthMiddleware())
	{
		checkoutapi.POST("/checkout", checkout.Checkout)
	}

}
