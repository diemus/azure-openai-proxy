package main

import (
	"github.com/diemus/azure-openai-proxy/pkg/azure"
	"github.com/diemus/azure-openai-proxy/pkg/openai"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var (
	Address   = "0.0.0.0:8080"
	ProxyMode = "azure"
)

func init() {
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
	gin.SetMode(gin.ReleaseMode)
	r.Any("*path", func(c *gin.Context) {
		if ProxyMode == "azure" {
			server := azure.NewOpenAIReverseProxy()
			server.ServeHTTP(c.Writer, c.Request)
		} else {
			server := openai.NewOpenAIReverseProxy()
			server.ServeHTTP(c.Writer, c.Request)
		}
	})

	r.Run(Address)

}
