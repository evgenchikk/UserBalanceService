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

// @Summary      De-reserve money from user's reserved balance
// @Description  Responds with the "dereserve" request body if OK.
// @Tags         user balance
// @Accept       json
// @Param        ReserveSchemaJSON body api.DereserveSchemaJSON true "Request body"
// @Produce      json
// @Success      201 {object} api.DereserveSchemaJSON
// @Failure	     400,500 {object} api.ErrorResponseJSON
// @Router       /dereserve [post]
func (h handler) dereserve(c *gin.Context) {
	var requestJSON api.DereserveSchemaJSON
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

	// get the requseted order entity
	query = "select order_id, user_id, price from userbalancedb.orders where order_id=$1"
	if err := tx.QueryRow(context.Background(), query, requestJSON.OrderID).Scan(&foundOrder.OrderID, &foundOrder.UserID, &foundOrder.Price); err != nil {
		if err != pgx.ErrNoRows {
			log.Printf("Select failed: %s", err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"no content": "order not found"})
		}
		return
	}

	// get user who owns the requsted order
	query = "select user_id, real_balance, reserved_balance from userbalancedb.users where user_id = $1"
	if err := tx.QueryRow(context.Background(), query, foundOrder.UserID).Scan(&foundUser.UserID, &foundUser.RealBalance, &foundUser.ReservedBalance); err != nil {
		if err != pgx.ErrNoRows {
			log.Printf("Select failed: %s", err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"no content": "user not found"})
		}
		return
	}

	// check if the order is already in progress
	query = "select sum(amount) from userbalancedb.payment_history where order_id = $1 group by order_id"
	if err := tx.QueryRow(context.Background(), query, foundOrder.OrderID).Scan(&alreadyPaidAmount); err != nil {
		if err != pgx.ErrNoRows {
			log.Printf("Select failed: %s", err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
			return
		}
	}

	if alreadyPaidAmount > 0 {
		log.Println("Can't de-reserve money for requsted order. The order is being executed.")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"bad requset": "Can't de-reserve money for requsted order. The order is being executed."})
		return
	}

	// update user balance
	query = "update userbalancedb.users set real_balance = $1, reserved_balance = $2 where user_id = $3"
	if _, err := tx.Exec(context.Background(), query, foundUser.RealBalance+foundOrder.Price, foundUser.ReservedBalance-foundOrder.Price, foundUser.UserID); err != nil {
		log.Printf("Update failed: %s", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		return
	}

	// delete order record
	query = "delete from userbalancedb.orders where order_id=$1"
	if _, err := tx.Exec(context.Background(), query, requestJSON.OrderID); err != nil {
		log.Printf("Delete failed: %s", err.Error())
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
