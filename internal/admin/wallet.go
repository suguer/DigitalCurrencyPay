package admin

import (
	"DigitalCurrency/internal/model/mdb"
	"DigitalCurrency/internal/service/wallet"
	"DigitalCurrency/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// WalletIndex godoc
// @Summary 钱包列表
// @Description 获取钱包列表
// @Tags 后台钱包管理
// @Accept json
// @Produce json
// @Param data body object true "分页参数"
// @Success 200 {object} util.Response "成功"
// @Failure 400 {object} util.Response "请求错误"
// @Router /admin/wallet/index [post]
func WalletIndex(c *gin.Context) {
	var request struct {
		mdb.Pagination
	}
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	walletList, pagination, err := wallet.Index(request.Current, request.PageSize, nil)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.SuccessData(c, walletList, pagination)
}
