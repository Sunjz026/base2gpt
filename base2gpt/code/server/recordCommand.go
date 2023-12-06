package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
)

func ListRecords(c *gin.Context) {
	table_id := c.Param("table_id")
	// 创建 Client
	// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
	client := lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(true))

	// 创建请求对象
	req := larkbitable.NewListAppTableRecordReqBuilder().
		AppToken(appToken).
		TableId(table_id).
		PageSize(20).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Bitable.AppTableRecord.List(context.Background(), req)

	for _, record := range resp.Data.Items {
		if _, ok := record.Fields["Attachment"]; ok {
			// 如果 Fields 包含键 "Attachment"，则调用函数 Y
			// 创建请求对象
			mediaToken, ok := record.Fields["Attachment"].(string)
			if ok {
				req2 := larkdrive.NewDownloadMediaReqBuilder().FileToken(mediaToken).Build()
				resp2, err := client.Drive.Media.Download(context.Background(), req2)
				// 处理错误
				if err != nil {
					fmt.Println(err)
					return
				}

				// 服务端错误处理
				if !resp.Success() {
					fmt.Println(resp2.Code, resp2.Msg, resp2.RequestId())
					return
				}

				record.Fields["Attachment"] = resp2.File
			}
		} else {
			fmt.Println("Value is not a string")
		}

	}

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

func FindRecord(c *gin.Context) {
	app_token := c.Param("app_token")
	table_id := c.Param("table_id")
	record_id := c.Param("record_id")
	// 创建 Client
	// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
	client := lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(true))

	// 创建请求对象
	req := larkbitable.NewGetAppTableRecordReqBuilder().
		AppToken(app_token).
		TableId(table_id).
		RecordId(record_id).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Bitable.AppTableRecord.Get(context.Background(), req)

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

func AddRecord(c *gin.Context) {
	app_token := c.Param("app_token")
	table_id := c.Param("table_id")

	var mediaData MediaData

	// 使用Gin绑定功能解析JSON请求体
	if err := c.BindJSON(&mediaData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取URL
	imageURL := mediaData.FileUrl

	// 使用HTTP客户端下载图片
	resp, err := http.Get(mediaData.FileUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to download image from %s", imageURL)})
		return
	}

	// Create a file

	// 设置随机数种子，通常应该只设置一次
	rand.Seed(time.Now().UnixNano())

	// 生成一个随机整数
	randomNumber := rand.Intn(100000) // 生成一个0到100000之间的随机整数
	file, err := os.Create("image_" + strconv.Itoa(randomNumber) + ".jpg")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create file")
		return
	}
	defer file.Close()

	// Write the response body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to save image to file")
		return
	}

	imageFile, err := os.Open(file.Name())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer imageFile.Close()

	fileInfo, err := imageFile.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	client := lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(true))

	// 创建请求对象
	UploadImageRequest := larkdrive.NewUploadAllMediaReqBuilder().
		Body(larkdrive.NewUploadAllMediaReqBodyBuilder().
			FileName(fileInfo.Name()).
			ParentType(`bitable_image`).
			ParentNode(app_token).
			Size(int(fileInfo.Size())).
			File(imageFile).
			Build()).
		Build()
	UploadImageResponse, err := client.Drive.Media.UploadAll(context.Background(), UploadImageRequest)

	// 处理错误
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
	}

	// 服务端错误处理
	if !UploadImageResponse.Success() {
		c.JSON(404, gin.H{
			"error": UploadImageResponse.Msg,
		})
	}
	var arr []map[string]string
	arr = append(arr, map[string]string{
		"file_token": *UploadImageResponse.Data.FileToken,
	})

	tag_arr := mediaData.Tags

	// 创建请求对象
	AddRecordRequest := larkbitable.NewCreateAppTableRecordReqBuilder().
		AppToken(app_token).
		TableId(table_id).
		AppTableRecord(larkbitable.NewAppTableRecordBuilder().
			Fields(map[string]interface{}{
				`Prompt`: mediaData.Prompt,
				`图片`:     arr,
				`标签`:     tag_arr,
			}).
			Build()).
		Build()

	AddRecordResponse, err := client.Bitable.AppTableRecord.Create(context.Background(), AddRecordRequest)

	// 处理错误
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
	}

	// 服务端错误处理
	if !AddRecordResponse.Success() {
		c.JSON(404, gin.H{
			"error": AddRecordResponse.Msg,
		})
	}

	// 业务处理
	log.Println(larkcore.Prettify(AddRecordResponse))
	c.JSON(200, gin.H{
		"messege": larkcore.Prettify(AddRecordResponse.Data),
	})

}
