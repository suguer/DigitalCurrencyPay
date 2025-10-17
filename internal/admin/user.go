package admin

import (
	"DigitalCurrency/internal/model/mdb"
	"DigitalCurrency/internal/service/user"
	"DigitalCurrency/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func UserIndex(c *gin.Context) {
	var request struct {
		mdb.Pagination
	}
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userList, pagination, err := user.Index(request.Current, request.PageSize, nil)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.SuccessData(c, userList, pagination)
}

// Login godoc
// @Summary 后台创建用户
// @Description 后台创建用户
// @Tags 后台用户管理
// @Accept json
// @Produce json
// @Param data body object true "登录信息"
// @Success 200 {object} util.Response "成功"
// @Failure 400 {object} util.Response "请求错误"
// @Router /admin/user/create [post]
func UserCreate(c *gin.Context) {
	c.JSON(200, gin.H{"message": "admin user create"})
	var request struct {
		mdb.User
	}
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := user.Create(&request.User); err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.SuccessResponse(c, request.User)
}
