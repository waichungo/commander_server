package httphandler

import (
	"commander_server/models"

	"github.com/gin-gonic/gin"
)

var UserFields = []string{
	"name",
	"email",
}
var UserModel = models.User{}

func HandleUserPost(c *gin.Context) {
	HandlePost(c, UserModel)
}

func HandleUserFind(c *gin.Context) {
	HandleFind(c, UserModel)
}

func HandleUserUpdate(c *gin.Context) {
	HandleUpdate(c, UserModel)
}

func HandleUserDelete(c *gin.Context) {
	HandleDelete(c, UserModel)
}
func HandleUserList(c *gin.Context) {
	HandleList(c, UserModel, UserFields)
}
