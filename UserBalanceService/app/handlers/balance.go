package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/evgenchikk/avito-tech-internship_backend_2022/api"
	_ "github.com/evgenchikk/avito-tech-internship_backend_2022/docs"
	"github.com/evgenchikk/avito-tech-internship_backend_2022/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

// @Summary      Get user balance
// @Description  Responds with the user balance as JSON.
// @Tags         user balance
// @Accept       json
// @Param        BalanceSchemaJSON body api.BalanceSchemaJSON true "Request body"
// @Produce      json
// @Success      200 {object} api.BalanceResponseJSON
// @Failure	     400,500 {object} api.ErrorResponseJSON
// @Router       /balance [post]
func (h handler) balance(c *gin.Context) {
	var requestJSON api.BalanceSchemaJSON
	var balance float64
	var query string

	if err := utils.CheckBinding(c, &requestJSON); err != nil {
		return
	}

	// get requsted user's balance
	query = "select real_balance from userbalancedb.users where user_id=$1"
	if err := h.DB.QueryRow(context.Background(), query, requestJSON.UserID).Scan(&balance); err != nil {
		// check if there is an error: no rows or smth else
		if err != pgx.ErrNoRows {
			log.Printf("Select failed: %s", err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"bad request": "user not found"})
		}
		return
	}

	// c.IndentedJSON(http.StatusOK, gin.H{"user balance": balance})
	c.IndentedJSON(http.StatusOK, api.BalanceResponseJSON{Balance: balance})
}
