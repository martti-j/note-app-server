package main

import (
	//user "note_app_server/routers"
	db "note_app_server/services"

	"github.com/gin-gonic/gin"
)

func main() {
	db.CheckConnectionToDB()
	router := gin.Default()
	router.GET("/users", db.GetUsernames)
	router.Run("localhost:8000")
}
