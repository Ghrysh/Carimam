package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func proxyTo(targetURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(targetURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Target URL tidak valid"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)

		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = c.Request.URL.Path
			req.URL.RawQuery = c.Request.URL.RawQuery
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	r := gin.Default()

	identityService := "http://localhost:8080"
	productService := "http://localhost:8081"
	orderService := "http://localhost:8082"

	r.POST("/register", proxyTo(identityService))
	r.POST("/login", proxyTo(identityService))
	r.GET("/api/profile", proxyTo(identityService))

	r.Any("/api/products/*filepath", proxyTo(productService))
	r.Any("/api/products", proxyTo(productService))
	r.Static("/uploads", "../product-service/uploads")

	r.Any("/api/orders", proxyTo(orderService))
	r.Any("/api/orders/*filepath", proxyTo(orderService))

	r.Any("/api/cart", proxyTo(orderService))
	r.Any("/api/cart/*filepath", proxyTo(orderService))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong! CariMam API Gateway is Running 🌐🚀"})
	})

	log.Println("API Gateway berjalan di Port 8000...")
	r.Run(":8000")
}