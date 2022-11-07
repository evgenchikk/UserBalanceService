package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/evgenchikk/avito-tech-internship_backend_2022/api"
	"github.com/evgenchikk/avito-tech-internship_backend_2022/app/models"
	"github.com/evgenchikk/avito-tech-internship_backend_2022/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

// @Summary      Approve money from user's reserved balance (debit money to the company's revenue)
// @Description  Responds with the "approve" request body if OK.
// @Tags         user balance
// @Accept       json
// @Param        ApproveSchemaJSON body api.ApproveSchemaJSON true "Request body"
// @Produce      json
// @Success      201 {object} api.ApproveSchemaJSON
// @Failure	     400,500 {object} api.ErrorResponseJSON
// @Router       /approve [post]
func (h handler) approve(c *gin.Context) {
	var requestJSON api.ApproveSchemaJSON
	var foundUser models.Users
	var foundOrder models.Orders
	var alreadyPaidAmount float64
	var query string

	if err := utils.CheckBinding(c, &requestJSON); err != nil {
		return
	}

	tx, err := h.DB.Begin(context.Background())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
	}
	defer tx.Rollback(context.Background())

	// get the requsted user
	query = "select user_id, real_balance, reserved_balance from userbalancedb.users where user_id=$1"
	if err := tx.QueryRow(context.Background(), query, requestJSON.UserID).Scan(&foundUser.UserID, &foundUser.RealBalance, &foundUser.ReservedBalance); err != nil {
		// check if there is an error: no rows or smth else
		if err != pgx.ErrNoRows {
			log.Printf("Select failed: %s", err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"no content": "user not found"})
		}
		return
	}

	// get the order with the requsted params: user_id, service_id, order_id
	query = "select order_id, user_id, service_id, activated_at, price from userbalancedb.orders where user_id=$1 and service_id=$2 and order_id=$3"
	if err := tx.QueryRow(context.Background(), query, requestJSON.UserID, requestJSON.ServiceID, requestJSON.OrderID).Scan(&foundOrder.OrderID, &foundOrder.UserID, &foundOrder.ServiceID, &foundOrder.ActivatedAt, &foundOrder.Price); err != nil {
		// check if there is an error: no rows or smth else
		if err != pgx.ErrNoRows {
			log.Printf("Select failed: %s", err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"no content": "order not found"})
		}
		return
	}

	// get the already paid amount for found order (to avoid overpayment)
	query = "select sum(amount) from userbalancedb.payment_history where order_id = $1 group by order_id"
	if err := tx.QueryRow(context.Background(), query, foundOrder.OrderID).Scan(&alreadyPaidAmount); err != nil {
		// check if there is an error: no rows or smth else
		if err != pgx.ErrNoRows {
			log.Printf("Select failed: %s", err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
			return
		}
		// in another case - the case when there are now rows in payment_history for requested order. alreadyPaidAmount = 0 (default int value)
	}

	// check if the amount is too large
	if alreadyPaidAmount+requestJSON.Amount > foundOrder.Price {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"bad request":                 "amount is too large for the requested order",
			"remaining amount to be paid": foundOrder.Price - alreadyPaidAmount,
		})
		return
	}

	// update found user's balance (debit money from its reserved balance)
	query = "update userbalancedb.users set reserved_balance = $1 where user_id = $2"
	if _, err := tx.Exec(context.Background(), query, foundUser.ReservedBalance-requestJSON.Amount, foundUser.UserID); err != nil {
		log.Printf("Update failed: %s", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		return
	}

	// register payment transaction
	query = "insert into userbalancedb.payment_history (order_id, payment_date, amount) values ($1, $2, $3)"
	if _, err := tx.Exec(context.Background(), query, requestJSON.OrderID, time.Now(), requestJSON.Amount); err != nil {
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
