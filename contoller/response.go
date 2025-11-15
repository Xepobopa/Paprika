package controller

import (
	"Paprika/process"

	"github.com/gin-gonic/gin"
)

// general
func responseWithError(err error) gin.H {
	return gin.H{
		"status": "error",
		"reason": err.Error(),
	}
}

// Exchange (CEX)
func responseGetStatsByName(stat *process.ProcessInformation) gin.H {
	return gin.H{
		"status": "success",
		"result": map[string]interface{}{
			"start":       stat.Start,
			"end":         stat.End,
			"status":      stat.Status.String(),
			"processName": stat.ProcessName,
			"error":       stat.Err,
		},
	}
}
func responseStopByName() gin.H {
	return gin.H{
		"status": "success",
	}
}
