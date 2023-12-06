package server

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

func ListDashBoards(c *gin.Context) {
	app_token := c.Param("app_token")
	// 创建 Client
	// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
	client := lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(true))

	req := larkbitable.NewListAppDashboardReqBuilder().
		AppToken(app_token).
		Build()
	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Bitable.AppDashboard.List(context.Background(), req)

	// 处理错误
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		log.Println(resp.Code, resp.Msg, resp.RequestId())
		c.JSON(404, gin.H{
			"error": resp.Msg,
		})
		return
	}

	// 业务处理
	log.Println(larkcore.Prettify(resp))
	c.JSON(200, gin.H{
		"messege": larkcore.Prettify(resp),
	})
}
