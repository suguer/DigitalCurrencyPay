package router

import (
	"DigitalCurrency/internal/admin"
	"DigitalCurrency/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Register(r *gin.Engine) *gin.Engine {

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// API 路由组
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		adm := api.Group("/admin")
		adm.Use(middleware.AdminMiddleware())
		{
			walletRouter := adm.Group("/wallet")
			{
				walletRouter.POST("", admin.WalletIndex)
			}
			configurationRouter := adm.Group("/configuration")
			{
				configurationRouter.POST("", admin.ConfigurationIndex)
			}
			userRouter := adm.Group("/user")
			{
				userRouter.POST("", admin.UserIndex)
				userRouter.POST("/create", admin.UserCreate)
			}
		}

		// 交易相关接口
		transactions := api.Group("/transactions")
		{
			// TODO: 添加交易相关路由
			transactions.POST("/create", TransactionCreate)
			transactions.POST("/query/:out_trade_no", TransactionInstance)
			transactions.POST("/repair/:out_trade_no", TransactionRepair)
		}

		depositRouter := api.Group("/deposit")
		{
			depositRouter.POST("", DepositIndex)
		}
		userRouter := api.Group("/user")
		{
			userRouter.POST("", UserInstance)
		}

	}

	return r
}
