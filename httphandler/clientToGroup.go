package httphandler

import (
	"commander_server/models"

	"github.com/gin-gonic/gin"
)

var ClientToGroupFields = []string{
	"client_id",
	"client_group_id",
}
var ClientToGroupModel = models.GroupToClient{}

func HandleClientToGroupPost(c *gin.Context) {
	HandlePost(c, ClientToGroupModel)
}

func HandleClientToGroupFind(c *gin.Context) {
	HandleFind(c, ClientToGroupModel)
}

func HandleClientToGroupUpdate(c *gin.Context) {
	HandleUpdate(c, ClientToGroupModel)
}

func HandleClientToGroupDelete(c *gin.Context) {
	HandleDelete(c, ClientToGroupModel)
}
func HandleClientToGroupList(c *gin.Context) {
	HandleList(c, ClientToGroupModel, ClientToGroupFields)
}
