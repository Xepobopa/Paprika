package app

import (
	"Paprika/process"
	"Paprika/publisher"
	"Paprika/router"
	"Paprika/utils"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

// customWriter - writer for logs
func Run(config *utils.Config, customWriter io.Writer) error {
	server := gin.Default()
	server.Use(gin.LoggerWithWriter(customWriter))
	server.Use(gin.Recovery())

	ctx := context.Background()

	pub, err := publisher.New("5011")
	if err != nil {
		return fmt.Errorf("failed to start grpc server: %v", err)
	}
	go func() {
		if err := pub.Listen(); err != nil {
			log.Printf("Publisher error: %v", err)
			pub.Stop()
			return
		}
	}()

	pm := process.NewManager(ctx)
	pm.Spawn("analyzer", process.SpawnAnalyzerProcess(pub))

	router.ApplyExchangeRouter(server, pm)

	return server.Run(":8080")
}
