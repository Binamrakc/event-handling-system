package main

import (
	middleware "crud/Middleware"
	"crud/controller"
	"crud/intializer"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.Use(middleware.ErrorHandler())
	intializer.Loadenv()
	intializer.ConnectDB()
	intializer.DBmigrate()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))
	r.Static("/uploads", "./uploads")

	protected := r.Group("/")
	protected.Use(middleware.Protected())
	r.Use(middleware.Ratelimit())
	//cors-done
	//http timeout-done
	//controller,server,repositry archeitures
	//pagination,
	// caching(redis)
	// ratelimit(middleware)-done
	//dependency injection

	s := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "userlogin",
		})
	})

	r.POST("/submit", controller.Submit)
	// r.GET("/view", controller.View)

	// r.GET("/update/:id", controller.Take)
	// r.POST("/update/:id", controller.Update)

	r.DELETE("/users/:id", controller.Delete)

	// r.GET("/register", controller.Showregister)
	// r.POST("/register", controller.Register)

	// r.GET("/auth", controller.Showauth)
	// r.POST("/auth", controller.Auth)
	// r.GET("/search", controller.Search)
	r.GET("/users", controller.Users)
	protected.POST("/events", controller.CreateEvent)
	r.GET("/events", controller.FetchEvent)
	protected.PUT("/events/:id/approve", controller.Approve)
	protected.DELETE("admin/events/:id", controller.Reject)
	r.GET("/dashboard/events", controller.Dashboard)
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	protected.GET("/profile/:id", controller.Profile)
	r.GET("/check-email/:email", controller.Verifyemail)
	protected.PUT("/change-password", controller.Changepass)
	s.ListenAndServe()
}

//localhost:8080/signup
