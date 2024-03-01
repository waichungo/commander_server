package httphandler

import (
	"commander_server/models"

	"github.com/gin-gonic/gin"
)

var DownloadFields = []string{
	"client_id",
	"progress",
	"name",
	"status",
}
var DownloadModel = models.Download{}

func HandleDownloadPost(c *gin.Context) {
	HandlePost(c, DownloadModel)
}

func HandleDownloadFind(c *gin.Context) {
	HandleFind(c, DownloadModel)
}

func HandleDownloadUpdate(c *gin.Context) {
	HandleUpdate(c, DownloadModel)
}

func HandleDownloadDelete(c *gin.Context) {
	HandleDelete(c, DownloadModel)
}
func HandleDownloadList(c *gin.Context) {
	HandleList(c, DownloadModel, DownloadFields)
}
