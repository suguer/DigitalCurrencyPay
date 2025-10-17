package router

import (
	"DigitalCurrency/internal/service/user"
	"DigitalCurrency/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserInstance(c *gin.Context) {
	userI, err := user.Instance(c.GetInt("user_id"))
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	util.SuccessResponse(c, gin.H{
		"Id":        userI.ID,
		"Username":  userI.Username,
		"Secret":    userI.Secret,
		"CreatedAt": userI.CreatedAt,
	})
}
