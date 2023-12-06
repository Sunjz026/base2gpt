package server

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义处理POST请求的路由
func GenerateImage(c *gin.Context) {

	var imageDescription ImageDescription

	// 使用Gin绑定功能解析JSON请求体
	if err := c.BindJSON(&imageDescription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置API端点URL和请求体
	apiUrl := "https://api.openai.com/v1/images/generations"
	requestBody := []byte(`{
            "model": "dall-e-3",
            "prompt": "` + imageDescription.Description + `",
            "n": 1,
            "size": "1024x1024"
        }`)

	// 创建HTTP请求客户端
	client := &http.Client{}

	// 创建POST请求
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置请求头，包括Authorization头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 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")

	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 将响应体作为JSON返回
	c.Data(http.StatusOK, "application/json", responseBody)
}
