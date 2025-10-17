package router

import (
	"DigitalCurrency/internal/service/deposit"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DepositIndex(c *gin.Context) {
	userID, _ := c.Get("user_id")
	depositList, err := deposit.DepositList(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, depositList)
}
