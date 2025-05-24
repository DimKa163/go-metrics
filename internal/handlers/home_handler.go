package handlers

import (
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomeHandler(repository persistence.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		metrics := repository.GetAll()
		c.Writer.Header().Set("Content-Type", "text/html")
		c.HTML(http.StatusOK, "home.tmpl", gin.H{
			"metrics": metrics,
		})
	}
}
