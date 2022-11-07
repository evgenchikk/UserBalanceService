package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/evgenchikk/avito-tech-internship_backend_2022/api"
	"github.com/evgenchikk/avito-tech-internship_backend_2022/app/models"
	"github.com/evgenchikk/avito-tech-internship_backend_2022/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

// @Summary      Transfer money from one user to another (creates user if not exists)
// @Description  Responds with the "add" request body if OK.
// @Tags         user balance
// @Accept       json
// @Param        AddSchemaJSON body api.TransferSchemaJSON true "Request body"
// @Produce      json
// @Success      201 {object} api.AddSchemaJSON
// @Failure	     400,500 {object} api.ErrorResponseJSON
// @Router       /transfer [post]
func (h handler) transfer(c *gin.Context) {
	var requestJSON api.TransferSchemaJSON
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

	// check user exisiting
	query = "select user_id, real_balance from userbalancedb.users where user_id = $1"
	if err := tx.QueryRow(context.Background(), query, requestJSON.FromUserID).Scan(&foundUser.UserID, &foundUser.RealBalance); err != nil {
		if err != pgx.ErrNoRows {
			log.Printf("Select failed: %s", err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"no content": "user not found"})
		}
		return
	}

	// try to debit money from user's balance
	query = "update userbalancedb.users set real_balance = $1 where user_id = $2"
	if _, err := tx.Exec(context.Background(), query, foundUser.RealBalance-requestJSON.Amount, foundUser.UserID); err != nil {
		log.Printf("Update failed: %s", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		return
	}

	// prepare the new request body to /add
	jsonBody := []byte(fmt.Sprintf("{\"user_id\":%v, \"amount\":%v}", requestJSON.ToUserID, requestJSON.Amount))
	bodyReader := bytes.NewReader(jsonBody)
	c.Request.Body = io.NopCloser(bodyReader)

	// set the current transaction to executable context
	c.Set("tx", tx)

	// call add method to add money
	h.add(c)
}
