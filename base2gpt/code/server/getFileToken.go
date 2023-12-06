package server

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
)

func GetFileToken(c *gin.Context) {
	app_token := c.Param("app_token")
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

	// 创建一个临时文件

	// 设置随机数种子，通常应该只设置一次
	rand.Seed(time.Now().UnixNano())

	// 生成一个随机整数
	// randomNumber := rand.Intn(100)                                                 // 生成一个0到99之间的随机整数
	tempFile, err := os.CreateTemp("", "image") // 生成一个随机文件名
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tempFile.Close()

	// 将图片内容写入临时文件
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 将临时文件转换为os.File
	imageFile, err := os.Open(tempFile.Name())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer imageFile.Close()

	// 获取文件信息
	fileInfo, err := imageFile.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	// 创建 Client
	// 如需SDK自动管理租户Token的获取与刷新，可调用lark.WithEnableTokenCache(true)进行设置
	client := lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(true))

	// 创建请求对象
	UploadImageRequest := larkdrive.NewUploadAllMediaReqBuilder().
		Body(larkdrive.NewUploadAllMediaReqBodyBuilder().
			FileName(tempFile.Name()).
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

	// 业务处理

	c.JSON(200, gin.H{
		"messege": larkcore.Prettify(UploadImageResponse),
	})

}
