package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/evgenchikk/avito-tech-internship_backend_2022/api"
	"github.com/evgenchikk/avito-tech-internship_backend_2022/app/models"
	"github.com/evgenchikk/avito-tech-internship_backend_2022/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

// @Summary      Reserve money from user's real balance (deposit money to user's reserved balance)
// @Description  Responds with the "reserve" request body if OK.
// @Tags         user balance
// @Accept       json
// @Param        ReserveSchemaJSON body api.ReserveSchemaJSON true "Request body"
// @Produce      json
// @Success      201 {object} api.ReserveSchemaJSON
// @Failure	     400,500 {object} api.ErrorResponseJSON
// @Router       /reserve [post]
func (h handler) reserve(c *gin.Context) {
	var requestJSON api.ReserveSchemaJSON
	var foundUser models.Users
	var query string

	if err := utils.CheckBinding(c, &requestJSON); err != nil {
		return
	}

	tx, err := h.DB.Begin(context.Background())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
	}
	defer tx.Rollback(context.Background())

	query = "select user_id, real_balance, reserved_balance from userbalancedb.users where user_id = $1"
	if err := tx.QueryRow(context.Background(), query, requestJSON.UserID).Scan(&foundUser.UserID, &foundUser.RealBalance, &foundUser.ReservedBalance); err != nil {
		if err != pgx.ErrNoRows {
			log.Printf("Select failed: %s", err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"no content": "user not found"})
		}
		return
	}

	query = "update userbalancedb.users set real_balance = $1, reserved_balance = $2 where user_id = $3"
	if _, err := tx.Exec(context.Background(), query, foundUser.RealBalance-requestJSON.Price, foundUser.ReservedBalance+requestJSON.Price, foundUser.UserID); err != nil {
		log.Printf("Update failed: %s", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		return
	}

	query = "insert into userbalancedb.orders (order_id, user_id, service_id, price) values ($1, $2, $3, $4)"
	if _, err := tx.Exec(context.Background(), query, requestJSON.OrderID, requestJSON.UserID, requestJSON.ServiceID, requestJSON.Price); err != nil {
		log.Printf("Insert failed: %s", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		log.Printf("Transaction commit failed: %s", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, requestJSON)
}
