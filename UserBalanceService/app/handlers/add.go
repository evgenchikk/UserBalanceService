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

// @Summary      Add money to user's balance (creates user if not exists)
// @Description  Responds with the "add" request body if OK.
// @Tags         user balance
// @Accept       json
// @Param        AddSchemaJSON body api.AddSchemaJSON true "Request body"
// @Produce      json
// @Success      201 {object} api.AddSchemaJSON
// @Failure	     400,500 {object} api.ErrorResponseJSON
// @Router       /add [post]
func (h handler) add(c *gin.Context) {
	var requestJSON api.AddSchemaJSON
	var foundUser models.Users
	var query string
	var err error
	var tx pgx.Tx

	if err = utils.CheckBinding(c, &requestJSON); err != nil {
		return
	}

	// check if there is an existing transaction in executable context
	if _tx, ok := c.Get("tx"); !ok {
		tx, err = h.DB.Begin(context.Background())
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		}
	} else {
		tx = _tx.(pgx.Tx)
	}

	defer tx.Rollback(context.Background())

	// get the requested user
	query = "select user_id, real_balance, reserved_balance from userbalancedb.users where user_id = $1"
	if err := tx.QueryRow(context.Background(), query, requestJSON.UserID).Scan(&foundUser.UserID, &foundUser.RealBalance, &foundUser.ReservedBalance); err != nil {
		// check if there is an error: no rows or smth else
		if err != pgx.ErrNoRows {
			log.Printf("Select failed: %s", err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
			return
		} else {
			// if user doesn't exist, create it
			newUser := models.Users{UserID: requestJSON.UserID, RealBalance: requestJSON.Amount}

			query = "insert into userbalancedb.users (user_id, real_balance) values ($1, $2)"
			if _, err := tx.Exec(context.Background(), query, newUser.UserID, newUser.RealBalance); err != nil {
				log.Printf("Insert failed: %s", err.Error())
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
				return
			}

			if err := tx.Commit(context.Background()); err != nil {
				log.Printf("Transaction commit failed: %s", err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
				return
			}
			c.IndentedJSON(http.StatusCreated, requestJSON)
			return
		}
	}

	// update found user's balance
	query = "update userbalancedb.users set real_balance = $1 where user_id = $2"
	if _, err := h.DB.Exec(context.Background(), query, foundUser.RealBalance+requestJSON.Amount, foundUser.UserID); err != nil {
		log.Printf("Update failed: %s", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		return
	}

	// commit the transaction
	if err := tx.Commit(context.Background()); err != nil {
		log.Printf("Transaction commit failed: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, requestJSON)
}
