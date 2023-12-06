package main

import (
	"demo/server"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 初始化表格
	r.GET("/init/table", server.InitializeTable)

	// 创建多维表格
	r.POST("/create/bitable", server.CreateBitable)
	// 为多维表格设置权限
	r.PATCH("/get/permission/:app_token", server.GetPermission)

	//  操作数据表
	r.GET("/list/tables/:app_token", server.ListTables)
	r.POST("/add/table/:app_token", server.AddTable)

	// 操作仪表盘
	r.GET("/list/dashboards/:app_token", server.ListDashBoards)

	// 操作视图
	r.GET("/list/views/:app_token/:table_id", server.ListViews)
	r.GET("/find/view/:app_token/:table_id/:view_id", server.FindView)
	r.POST("/add/view/:app_token/:table_id", server.AddView)

	// 操作字段
	r.GET("/list/fields/:table_id", server.ListFields)
	r.POST("/add/field/:table_id", server.AddField)

	// 操作记录
	r.POST("/add/record/:app_token/:table_id", server.AddRecord)
	r.GET("/list/records/:app_token/:table_id", server.ListRecords)
	r.GET("/find/record/:app_token/:table_id/:record_id", server.FindRecord)

	// 前置操作
	r.POST("/get/file_token/:app_token", server.GetFileToken)
	r.POST("/generate/image", server.GenerateImage)
	r.GET("/download/image", server.DownloadImage)

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run(":9000")
}
