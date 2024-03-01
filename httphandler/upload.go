package httphandler

import (
	"commander_server/models"

	"github.com/gin-gonic/gin"
)

var uploadFields = []string{
	"username",
	"machine",
	"machine_id",
}
var uploadModel = models.Client{}

func HandleUploadPost(c *gin.Context) {
	HandlePost(c, uploadModel)
}

func HandleUploadFind(c *gin.Context) {
	HandleFind(c, uploadModel)
}

func HandleUploadUpdate(c *gin.Context) {
	HandleUpdate(c, uploadModel)
}

func HandleUploadDelete(c *gin.Context) {
	HandleDelete(c, uploadModel)
}
func HandleUploadList(c *gin.Context) {
	HandleList(c, uploadModel, uploadFields)
}
