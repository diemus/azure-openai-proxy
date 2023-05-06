package main

import (
	"fmt"
	"github.com/diemus/azure-openai-proxy/pkg/azure"
	"github.com/diemus/azure-openai-proxy/pkg/openai"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	Address   = "0.0.0.0:8080"
	ProxyMode = "azure"
	LogPath   = "log/proxy.log"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	if v := os.Getenv("AZURE_OPENAI_PROXY_ADDRESS"); v != "" {
		Address = v
	}
	if v := os.Getenv("AZURE_OPENAI_PROXY_MODE"); v != "" {
		ProxyMode = v
	}
	if v := os.Getenv("AZURE_OPENAI_PROXY_LOG_PATH"); v != "" {
		LogPath = v
	}
	logFile, err := os.OpenFile(LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	multi := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multi)
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
	log.Printf("loading azure openai proxy address: %s", Address)
	log.Printf("loading azure openai proxy mode: %s", ProxyMode)
}

func main() {
	router := gin.Default()
	if ProxyMode == "azure" {
		router.GET("/v1/models", handleGetModels)
		router.OPTIONS("/v1/*path", handleOptions)

		router.POST("/v1/chat/completions", handleAzureProxy)
		router.POST("/v1/completions", handleAzureProxy)
		router.POST("/v1/embeddings", handleAzureProxy)
	} else {
		router.Any("*path", handleOpenAIProxy)
	}

	_ = router.Run(Address)

}

func handleGetModels(c *gin.Context) {
	// BUGFIX: fix options request, see https://github.com/diemus/azure-openai-proxy/issues/3
	models := []string{"gpt-4", "gpt-4-0314", "gpt-4-32k", "gpt-4-32k-0314", "gpt-3.5-turbo", "gpt-3.5-turbo-0301", "text-davinci-003", "text-embedding-ada-002"}
	result := azure.ListModelResponse{
		Object: "list",
	}
	for _, model := range models {
		result.Data = append(result.Data, azure.Model{
			ID:      model,
			Object:  "model",
			Created: 1677649963,
			OwnedBy: "openai",
			Permission: []azure.ModelPermission{
				{
					ID:                 "",
					Object:             "model",
					Created:            1679602087,
					AllowCreateEngine:  true,
					AllowSampling:      true,
					AllowLogprobs:      true,
					AllowSearchIndices: true,
					AllowView:          true,
					AllowFineTuning:    true,
					Organization:       "*",
					Group:              nil,
					IsBlocking:         false,
				},
			},
			Root:   model,
			Parent: nil,
		})
	}
	c.JSON(200, result)
}

func handleOptions(c *gin.Context) {
	// BUGFIX: fix options request, see https://github.com/diemus/azure-openai-proxy/issues/1
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Status(200)
	return
}

func handleAzureProxy(c *gin.Context) {
	// BUGFIX: fix options request, see https://github.com/diemus/azure-openai-proxy/issues/1
	if c.Request.Method == http.MethodOptions {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Status(200)
		return
	}

	server := azure.NewOpenAIReverseProxy()
	server.ServeHTTP(c.Writer, c.Request)
	//BUGFIX: try to fix the difference between azure and openai
	//Azure's response is missing a \n at the end of the stream
	//see https://github.com/Chanzhaoyu/chatgpt-web/issues/831
	if c.Writer.Header().Get("Content-Type") == "text/event-stream" {
		if _, err := c.Writer.Write([]byte("\n")); err != nil {
			log.Printf("rewrite azure response error: %v", err)
		}
	}
}

func handleOpenAIProxy(c *gin.Context) {
	server := openai.NewOpenAIReverseProxy()
	server.ServeHTTP(c.Writer, c.Request)
}
