package httphandler

import (
	"commander_server/models"

	"github.com/gin-gonic/gin"
)

var taskFields = []string{
	"client_id",
	"client_group_id",
	"type",
}
var taskModel = models.Task{}

func HandleTaskPost(c *gin.Context) {
	HandlePost(c, taskModel)
}

func HandleTaskFind(c *gin.Context) {
	HandleFind(c, taskModel)
}

func HandleTaskUpdate(c *gin.Context) {
	HandleUpdate(c, taskModel)
}

func HandleTaskDelete(c *gin.Context) {
	HandleDelete(c, taskModel)
}
func HandleTaskList(c *gin.Context) {
	HandleList(c, taskModel, taskFields)
}
