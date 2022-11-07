package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckBinding(c *gin.Context, requestJSON interface{}) error {
	if err := c.BindJSON(requestJSON); err != nil {
		log.Println(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"comment": "binding request json failed",
		})
		return err
	}
	return nil
}
