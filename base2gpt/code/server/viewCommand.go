package server

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

func ListViews(c *gin.Context) {
	app_token := c.Param("app_token")
	table_id := c.Param("table_id")

	// 创建 Client
	// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
	client := lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(true))

	// 创建请求对象
	req := larkbitable.NewListAppTableViewReqBuilder().
		AppToken(app_token).
		TableId(table_id).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Bitable.AppTableView.List(context.Background(), req)

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

func FindView(c *gin.Context) {
	app_token := c.Param("app_token")
	table_id := c.Param("table_id")
	view_id := c.Param("view_id")

	// 创建 Client
	// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
	client := lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(true))

	// 创建请求对象
	req := larkbitable.NewGetAppTableViewReqBuilder().
		AppToken(app_token).
		TableId(table_id).
		ViewId(view_id).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Bitable.AppTableView.Get(context.Background(), req)

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

func AddView(c *gin.Context) {

	app_token := c.Param("app_token")
	table_id := c.Param("table_id")

	// var viewData ViewData

	// 使用Gin绑定功能解析JSON请求体
	// if err := c.BindJSON(&viewData); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// 创建 Client
	// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
	client := lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(true))

	// 创建请求对象 传参
	// req := larkbitable.NewCreateAppTableViewReqBuilder().
	// 	AppToken(app_token).
	// 	TableId(table_id).
	// 	ReqView(larkbitable.NewReqViewBuilder().
	// 		ViewName(viewData.Name).
	// 		ViewType(viewData.Type).
	// 		Build()).
	// 	Build()

	// 创建请求对象 写死
	req := larkbitable.NewCreateAppTableViewReqBuilder().
		AppToken(app_token).
		TableId(table_id).
		ReqView(larkbitable.NewReqViewBuilder().
			ViewName("画册视图").
			ViewType("gallery").
			Build()).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Bitable.AppTableView.Create(context.Background(), req)

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
