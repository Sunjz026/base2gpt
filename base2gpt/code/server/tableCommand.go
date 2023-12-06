package server

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

func ListTables(c *gin.Context) {
	app_token := c.Param("app_token")
	// 创建 Client
	// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
	client := lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(true))

	// 创建请求对象
	req := larkbitable.NewListAppTableReqBuilder().
		AppToken(app_token).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Bitable.AppTable.List(context.Background(), req)
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

func GetHeaders(tableData TableData) []*larkbitable.AppTableCreateHeader {
	headers := make([]*larkbitable.AppTableCreateHeader, 0, len(tableData.Table.Fields))
	// 遍历 Fields 数组
	for _, field := range tableData.Table.Fields {
		header := larkbitable.NewAppTableCreateHeaderBuilder().FieldName(field.FieldName).Type(field.Type).Build()
		// 对每个元素执行构建操作
		headers = append(headers, header)
	}
	return headers
}

func AddTable(c *gin.Context) {
	app_token := c.Param("app_token")
	// var tableData TableData

	// 使用Gin绑定功能解析JSON请求体
	// if err := c.BindJSON(&tableData); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// 处理请求并返回响应
	// 创建 Client
	// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
	client := lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(true))

	// 创建请求对象 动态
	// headers := GetHeaders(tableData)
	// req := larkbitable.NewCreateAppTableReqBuilder().
	// 	AppToken(app_token).
	// 	Body(larkbitable.NewCreateAppTableReqBodyBuilder().
	// 		Table(larkbitable.NewReqTableBuilder().
	// 			Name(tableData.Table.Name).
	// 			DefaultViewName(tableData.Table.DefaultViewName).
	// 			Fields(headers).
	// 			Build()).
	// 		Build()).
	// 	Build()

	// 创建请求对象 写死
	req := larkbitable.NewCreateAppTableReqBuilder().
		AppToken(app_token).
		Body(larkbitable.NewCreateAppTableReqBodyBuilder().
			Table(larkbitable.NewReqTableBuilder().
				Name(`图片表`).
				DefaultViewName(`默认的表格视图`).
				Fields([]*larkbitable.AppTableCreateHeader{
					larkbitable.NewAppTableCreateHeaderBuilder().
						FieldName(`Prompt`).
						Type(1).
						Build(),
					larkbitable.NewAppTableCreateHeaderBuilder().
						FieldName(`图片`).
						Type(17).
						Build(),
					larkbitable.NewAppTableCreateHeaderBuilder().
						FieldName(`标签`).
						Type(4).
						Build(),
					larkbitable.NewAppTableCreateHeaderBuilder().
						FieldName(`创建时间`).
						Type(1001).
						Build(),
				}).
				Build()).
			Build()).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Bitable.AppTable.Create(context.Background(), req)

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
