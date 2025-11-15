package controller

import (
	"Paprika/process"
	"Paprika/service"

	"github.com/gin-gonic/gin"
)

// ExchangeGETList returns all avaliable exchanges
func ExchangeGETList(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "success",
		"list":   service.GetAllAvailableExchanges(),
	})
}

func ExchangeSpawnByName(c *gin.Context, pm *process.Manager) {
	name := c.Param("name")

	if err := service.SpawnExchangeProcess(name, pm); err != nil {
		c.JSON(400, responseWithError(err))
	} else {
		c.JSON(200, gin.H{
			"status": "success",
		})
	}
}

func ExchangeGetStatsByName(c *gin.Context, pm *process.Manager) {
	name := c.Param("name")

	stats, err := service.GetExchangeStats(name, pm)
	if err != nil {
		c.JSON(400, responseWithError(err))
	} else {
		c.JSON(200, responseGetStatsByName(stats))
	}
}

func ExchangeGETProcessList(c *gin.Context, pm *process.Manager) {
	c.JSON(200, gin.H{
		"status": "success",
		"list":   service.GetAllAvailableExchanges(),
	})
}

func ExchangeStopByName(c *gin.Context, pm *process.Manager) {
	name := c.Param("name")

	if err := service.StopExchangeByName(name, pm); err != nil {
		c.JSON(400, responseWithError(err))
	} else {
		c.JSON(200, responseStopByName())
	}
}