package router

import (
	controller "Paprika/contoller"
	"Paprika/process"

	"github.com/gin-gonic/gin"
)

//TODO: analyzer

// TODO: add smth like "remove history by name(id)" endpoint, that will delete completed or failed process from the process manager
func ApplyExchangeRouter(server *gin.Engine, pm *process.Manager) {
	router := server.Group("/cex")

	router.GET("/list", controller.ExchangeGETList)
	router.GET("/processList", controller.ExchangeGETList)
	router.GET("/start/:name", func(ctx *gin.Context) {
		controller.ExchangeSpawnByName(ctx, pm)
	})
	router.GET("/stop/:name", func(ctx *gin.Context) {
		controller.ExchangeStopByName(ctx, pm)
	})
	router.GET("/stats/:name", func(ctx *gin.Context) {
		controller.ExchangeGetStatsByName(ctx, pm)
	})
}
