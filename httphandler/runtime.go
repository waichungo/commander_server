package httphandler

import (
	"commander_server/models"

	"github.com/gin-gonic/gin"
)

var runtimeFields = []string{
	"group_id",
	"name",
	"machine_id",
}
var runtimeModel = models.Runtime{}

func HandleRuntimePost(c *gin.Context) {
	HandlePost(c, runtimeModel)
}

func HandleRuntimeFind(c *gin.Context) {
	HandleFind(c, runtimeModel)
}

func HandleRuntimeUpdate(c *gin.Context) {
	HandleUpdate(c, runtimeModel)
}

func HandleRuntimeDelete(c *gin.Context) {
	HandleDelete(c, runtimeModel)
}
func HandleRuntimeList(c *gin.Context) {

	HandleList(c, runtimeModel, runtimeFields)
}
