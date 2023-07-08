package routers

import (
	"net/http"

	db "note_app_server/services"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	users, err := db.GetUsersDB()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to get users"})
		return
	}
	c.IndentedJSON(http.StatusOK, users)
}

func GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := db.GetUserByUsernameDB(username)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Username not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}

func AddUser(c *gin.Context) {
	var newUser db.User

	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	username := newUser.Username
	password := newUser.Password

	if err := db.AddUserDB(username, password); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to add user"})
		return
	}
	c.IndentedJSON(http.StatusOK, newUser)
}
