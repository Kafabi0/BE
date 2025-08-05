package main

import (
	"github.com/gin-gonic/gin"
	"rumahsakit/db"
	"rumahsakit/routes"
)

func main() {
	r := gin.Default()
	r.Use(corsMiddleware())

	db.ConnectDB()
	routes.PasienRoutes(r)
	routes.AuthRoutes(r) // <- hanya register & login

	r.Run(":8080")
}

func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")

        if origin == "http://localhost:5173" {
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
        }

        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}


