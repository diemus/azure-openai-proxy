package main

import (
	"github.com/diemus/azure-openai-proxy/pkg/azure"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var (
	Address = "0.0.0.0:8080"
)

func init() {
	if v := os.Getenv("AZURE_OPENAI_PROXY_ADDRESS"); v != "" {
		Address = v
	}
	log.Printf("loading azure openai proxy address: %s", Address)
}

func main() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.Any("*path", func(c *gin.Context) {
		server := azure.NewOpenAIReverseProxy()
		server.ServeHTTP(c.Writer, c.Request)
	})

	r.Run(Address)

}
