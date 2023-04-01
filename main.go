package main

import (
	"github.com/diemus/azure-openai-proxy/pkg/azure"
	"github.com/diemus/azure-openai-proxy/pkg/openai"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

var (
	Address   = "0.0.0.0:8080"
	ProxyMode = "azure"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	if v := os.Getenv("AZURE_OPENAI_PROXY_ADDRESS"); v != "" {
		Address = v
	}
	if v := os.Getenv("AZURE_OPENAI_PROXY_MODE"); v != "" {
		ProxyMode = v
	}
	log.Printf("loading azure openai proxy address: %s", Address)
	log.Printf("loading azure openai proxy mode: %s", ProxyMode)
}

func main() {
	r := gin.Default()
	r.Any("*path", func(c *gin.Context) {
		if ProxyMode == "azure" {
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
		} else {
			server := openai.NewOpenAIReverseProxy()
			server.ServeHTTP(c.Writer, c.Request)
		}
	})

	r.Run(Address)

}
