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

	if err := db.AddUserDB(newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, newUser)
}

func DeleteUser(c *gin.Context) {
	var deleteUser db.User

	if err := c.BindJSON(&deleteUser); err != nil {
		return
	}

	if err := db.DeleteUserDB(deleteUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"Deleted user": deleteUser.Username})
}

func Login(c *gin.Context) {
	var loginUser db.User

	if err := c.BindJSON(&loginUser); err != nil {
		return
	}

	if err := db.LoginDB(loginUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"Logged in as": loginUser.Username})
}
