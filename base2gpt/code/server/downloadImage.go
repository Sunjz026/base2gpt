package server

import (
	"context"
	//"fmt"
	"log"

	"github.com/gin-gonic/gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
)

func DownloadImage(c *gin.Context) {
	// 创建 Client
	// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
	client := lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(true))

	// 创建请求对象
	req := larkdrive.NewDownloadMediaReqBuilder().
		FileToken(`JYpobG3uro63hjxA62bcFN07nQd`).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Drive.Media.Download(context.Background(), req)

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

	err = resp.WriteFile("output.jpg")
	if err != nil {
		panic(err)
	}

	log.Println(larkcore.Prettify(resp))
	c.JSON(200, gin.H{
		"messege": larkcore.Prettify(resp),
	})
	resp.WriteFile("/Users/bytedance/Downloads")
}
