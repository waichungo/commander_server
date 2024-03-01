package httphandler

import (
	"commander_server/models"

	"github.com/gin-gonic/gin"
)

var ClientProfileFields = []string{
	"client_id",
}
var ClientProfileModel = models.ClientProfile{}

func HandleClientProfilePost(c *gin.Context) {
	HandlePost(c, ClientProfileModel)
}

func HandleClientProfileFind(c *gin.Context) {
	HandleFind(c, ClientProfileModel)
}

func HandleClientProfileUpdate(c *gin.Context) {
	HandleUpdate(c, ClientProfileModel)
}

func HandleClientProfileDelete(c *gin.Context) {
	HandleDelete(c, ClientProfileModel)
}
func HandleClientProfileList(c *gin.Context) {
	HandleList(c, ClientProfileModel, ClientProfileFields)
}
