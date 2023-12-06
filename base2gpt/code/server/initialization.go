package server

import (
	"context"

	"github.com/gin-gonic/gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
)

func InitializeTable(c *gin.Context) {
	client := lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(true))

	// 创建请求对象
	CreateBitableRequest := larkbitable.NewCreateAppReqBuilder().
		ReqApp(larkbitable.NewReqAppBuilder().
			Name(`一篇新的多维表格`).
			Build()).
		Build()

	// 发起请求
	CreateBitableResponse, err := client.Bitable.App.Create(context.Background(), CreateBitableRequest)
	// 处理错误
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
	}

	// 服务端错误处理
	if !CreateBitableResponse.Success() {
		c.JSON(404, gin.H{
			"error": CreateBitableResponse.Msg,
		})
	}
	app_token := CreateBitableResponse.Data.App.AppToken

	GetPermissionRequest := larkdrive.NewPatchPermissionPublicReqBuilder().
		Token(*app_token).
		Type(`bitable`).
		PermissionPublicRequest(larkdrive.NewPermissionPublicRequestBuilder().
			ExternalAccess(true).
			SecurityEntity(`anyone_can_view`).
			CommentEntity(`anyone_can_view`).
			ShareEntity(`anyone`).
			LinkShareEntity(`anyone_editable`).
			InviteExternal(true).
			Build()).
		Build()

	GetPermissionResponse, err := client.Drive.PermissionPublic.Patch(context.Background(), GetPermissionRequest)

	// 处理错误
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
	}

	// 服务端错误处理
	if !GetPermissionResponse.Success() {
		c.JSON(404, gin.H{
			"error": GetPermissionResponse.Msg,
		})
	}

	AddTableRequest := larkbitable.NewCreateAppTableReqBuilder().
		AppToken(*app_token).
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

	AddTableResponse, err := client.Bitable.AppTable.Create(context.Background(), AddTableRequest)

	// 处理错误
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
	}

	// 服务端错误处理
	if !AddTableResponse.Success() {
		c.JSON(404, gin.H{
			"error": AddTableResponse.Msg,
		})
	}
	table_id := AddTableResponse.Data.TableId

	AddViewRequest := larkbitable.NewCreateAppTableViewReqBuilder().
		AppToken(*app_token).
		TableId(*table_id).
		ReqView(larkbitable.NewReqViewBuilder().
			ViewName("画册视图").
			ViewType("gallery").
			Build()).
		Build()

	AddViewResponse, err := client.Bitable.AppTableView.Create(context.Background(), AddViewRequest)

	// 处理错误
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
	}

	// 服务端错误处理
	if !AddTableResponse.Success() {
		c.JSON(404, gin.H{
			"error": AddViewResponse.Msg,
		})
	}
	url := "https://bytedance.larkoffice.com/base/" + *app_token + "?table=" + *table_id
	Response := make(map[string]string)
	Response["app_token"] = *app_token
	Response["table_id"] = *table_id
	Response["url"] = url
	// 业务处理
	c.JSON(200, gin.H{
		"messege": larkcore.Prettify(Response),
	})
}
