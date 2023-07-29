package main

import (
	user "note_app_server/routers"
	db "note_app_server/services"

	"github.com/gin-gonic/gin"
)

func main() {
	db.CheckConnectionToDB()
	router := gin.Default()
	router.GET("/users", user.GetUsers)
	router.GET("/user/:username", user.GetUserByUsername)
	router.POST("/registration/", user.AddUser)
	router.DELETE("/user/", user.DeleteUser)
	router.POST("/login/", user.Login)
	router.GET("/notes/", user.GetNotes)
	router.POST("/note/", user.AddNote)
	router.DELETE("/note/", user.DeleteNote)
	router.Run("localhost:8000")
}
