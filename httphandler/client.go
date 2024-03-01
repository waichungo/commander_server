package httphandler

import (
	"commander_server/models"

	"github.com/gin-gonic/gin"
)

var ClientFields = []string{
	"username",
	"machine",
	"machine_id",
}
var clientModel = models.Client{}

func HandleClientPost(c *gin.Context) {
	HandlePost(c, clientModel)
}

func HandleClientFind(c *gin.Context) {
	HandleFind(c, clientModel)
}

func HandleClientUpdate(c *gin.Context) {
	HandleUpdate(c, clientModel)
}

func HandleClientDelete(c *gin.Context) {
	HandleDelete(c, clientModel)
}
func HandleClientList(c *gin.Context) {
	HandleList(c, clientModel, ClientFields)
}
