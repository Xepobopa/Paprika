package app

import (
	"Paprika/process"
	"Paprika/router"
	"Paprika/utils"
	"context"
	"io"

	"github.com/gin-gonic/gin"
)

// customWriter - writer for logs
func Run(config *utils.Config, customWriter io.Writer) error {
	server := gin.Default()
	server.Use(gin.LoggerWithWriter(customWriter))
	server.Use(gin.Recovery())

	// pb, err := publisher.NewPublisher(config.Nats.URL)
	// if err != nil {
	// 	return fmt.Errorf("failed to connect to the nats server: %v", err)
	// }

	ctx := context.Background()
	pm := process.NewManager(ctx)
	pm.Spawn("anallyzer", process.SpawnAnalyzerProcess())

	router.ApplyExchangeRouter(server, pm)

	return server.Run(":8080")
}
