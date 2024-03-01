package httphandler

import (
	"commander_server/models"

	"github.com/gin-gonic/gin"
)

var InfoFields = []string{
	"machine_id",
}
var infoModel = models.MachineInfo{}

func HandleInfoPost(c *gin.Context) {
	HandlePost(c, infoModel)
}

func HandleInfoFind(c *gin.Context) {
	HandleFind(c, infoModel)
}

func HandleInfoUpdate(c *gin.Context) {
	HandleUpdate(c, infoModel)
}

func HandleInfoDelete(c *gin.Context) {
	HandleDelete(c, infoModel)
}
func HandleInfoList(c *gin.Context) {

	HandleList(c, infoModel, InfoFields)
}
