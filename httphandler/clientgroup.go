package httphandler

import (
	"commander_server/models"

	"github.com/gin-gonic/gin"
)

var clientGroupFields = []string{
	"name",
}
var clientGroupModel = models.ClientGroup{}

func HandleClientGroupPost(c *gin.Context) {
	HandlePost(c, clientGroupModel)
}

func HandleClientGroupFind(c *gin.Context) {
	HandleFind(c, clientGroupModel)
}

func HandleClientGroupUpdate(c *gin.Context) {
	HandleUpdate(c, clientGroupModel)
}

func HandleClientGroupDelete(c *gin.Context) {
	HandleDelete(c, clientGroupModel)
}
func HandleClientGroupList(c *gin.Context) {
	HandleList(c, clientGroupModel, clientGroupFields)
}
