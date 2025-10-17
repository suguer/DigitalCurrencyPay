package admin

import (
	"DigitalCurrency/internal/model/mdb"
	"DigitalCurrency/internal/service/configuration"
	"DigitalCurrency/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// ConfigurationIndex godoc
// @Summary 配置列表
// @Description 获取配置列表
// @Tags 后台配置管理
// @Accept json
// @Produce json
// @Param data body object true "分页参数"
// @Success 200 {object} util.Response "成功"
// @Failure 400 {object} util.Response "请求错误"
// @Router /admin/configuration/index [post]
func ConfigurationIndex(c *gin.Context) {
	var request struct {
		mdb.Pagination
	}
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	configurationList, pagination, err := configuration.Index(request.Current, request.PageSize, nil)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.SuccessData(c, configurationList, pagination)
}
