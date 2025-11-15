package app

import (
	"Paprika/process"
	"Paprika/publisher"
	"Paprika/router"
	"Paprika/utils"
	"context"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

// customWriter - writer for logs
func Run(config *utils.Config, customWriter io.Writer) error {
	server := gin.Default()
	server.Use(gin.LoggerWithWriter(customWriter))
	server.Use(gin.Recovery())

	pb, err := publisher.NewPublisher(config.Nats.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to the nats server: %v", err)
	}
	
	ctx := context.Background()
	pm := process.NewManager(ctx, pb)

	router.ApplyExchangeRouter(server, pm)

	return server.Run(":8080")
}
