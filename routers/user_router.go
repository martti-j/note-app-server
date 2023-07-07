package routers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUsernames(c *gin.Context) {
	filter := bson.D{{}}
	fmt.Println(filter)
	c.IndentedJSON(http.StatusOK, filter)
}
