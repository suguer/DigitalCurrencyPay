package router

import (
	"DigitalCurrency/internal/constant"
	"DigitalCurrency/internal/model/mdb"
	"DigitalCurrency/internal/service/transaction"
	"DigitalCurrency/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func TransactionCreate(c *gin.Context) {
	instance := mdb.Transaction{}
	if err := c.ShouldBindBodyWith(&instance, binding.JSON); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	instance.UserId = userID.(uint)
	response, err := transaction.Create(&instance)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if response.Chain == constant.ChainTron || response.Chain == constant.ChainTronShasta {
		toAddr, _ := util.HexString2Address(response.ToAddress)
		response.ToAddress = toAddr
		contractAddr, _ := util.HexString2Address(response.ContractAddress)
		response.ContractAddress = contractAddr
	}

	util.SuccessResponse(c, gin.H{
		"id":               response.ID,
		"merchant_id":      response.UserId,
		"out_trade_no":     response.OutTradeNo,
		"address":          response.ToAddress,
		"amount":           response.Amount,
		"chain":            response.Chain,
		"contract_address": response.ContractAddress,
		"status":           constant.TransactionStatusMap(response.Status),
	})
}

func TransactionInstance(c *gin.Context) {
	outTradeNo := c.Param("out_trade_no")
	userID, _ := c.Get("user_id")
	tx, err := transaction.InstanceByOutTradeNo(outTradeNo, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tx.Chain == "tron" {
		tronlinkAddress, _ := util.HexString2Address(tx.ToAddress)
		tx.ToAddress = tronlinkAddress
		tronlinkAddress, _ = util.HexString2Address(tx.ContractAddress)
		tx.ContractAddress = tronlinkAddress
	}
	c.JSON(http.StatusOK, tx)

}

func TransactionRepair(c *gin.Context) {
	userID, _ := c.Get("user_id")
	outTradeNo := c.Param("out_trade_no")
	err := transaction.TransactionRepair(outTradeNo, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})

}
