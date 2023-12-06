package server

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
)

func CreateBitable(c *gin.Context) {

	client := lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(true))

	// 创建请求对象
	req := larkbitable.NewCreateAppReqBuilder().
		ReqApp(larkbitable.NewReqAppBuilder().
			Name(`一篇新的多维表格`).
			Build()).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Bitable.App.Create(context.Background(), req)

	// 处理错误
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
	}

	// 服务端错误处理
	if !resp.Success() {
		c.JSON(404, gin.H{
			"error": resp.Msg,
		})
	}

	// 业务处理
	log.Println(larkcore.Prettify(resp))
	c.JSON(200, gin.H{
		"messege": larkcore.Prettify(resp),
	})
}

func GetPermission(c *gin.Context) {

	app_token := c.Param("app_token")

	client := lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(true))

	// 创建请求对象
	req := larkdrive.NewPatchPermissionPublicReqBuilder().
		Token(app_token).
		Type(`bitable`).
		PermissionPublicRequest(larkdrive.NewPermissionPublicRequestBuilder().
			ExternalAccess(true).
			SecurityEntity(`anyone_can_view`).
			CommentEntity(`anyone_can_view`).
			ShareEntity(`anyone`).
			LinkShareEntity(`tenant_readable`).
			InviteExternal(true).
			Build()).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Drive.PermissionPublic.Patch(context.Background(), req)

	// 处理错误
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
	}

	// 服务端错误处理
	if !resp.Success() {
		c.JSON(404, gin.H{
			"error": resp.Msg,
		})
	}

	// 业务处理
	log.Println(larkcore.Prettify(resp))
	c.JSON(200, gin.H{
		"messege": larkcore.Prettify(resp),
	})
}
