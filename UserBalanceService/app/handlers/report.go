package handlers

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/evgenchikk/avito-tech-internship_backend_2022/api"
	"github.com/evgenchikk/avito-tech-internship_backend_2022/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type reportRow struct {
	service_id int
	amount     float64
}

// write slice data to csv file
// returns file name and error
func (h handler) writeReport(reportData []reportRow) (string, error) {
	filename := uuid.New().String() + ".csv"

	file, err := os.Create(fmt.Sprintf("%s/%s", h.ReportDirPath, filename))
	if err != nil {
		return "", err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	w.Comma = ';'
	defer w.Flush()

	for i, row := range reportData {
		if i == 0 {
			w.Write([]string{"Service", "Total revenue"})
		}

		if err := w.Write([]string{fmt.Sprintf("%v", row.service_id), fmt.Sprintf("%.3f", row.amount)}); err != nil {
			return filename, err
		}
	}

	return filename, nil
}

// @Summary      Create report according to the selected period
// @Description  Responds with the link to report download page if OK.
// @Tags         report
// @Accept       json
// @Param        ReportSchemaJSON body api.ReportSchemaJSON true "Request body"
// @Produce      json
// @Success      201 {object} api.ReportResponseJSON
// @Success      204 {object} api.ErrorResponseJSON
// @Failure	     400,500 {object} api.ErrorResponseJSON
// @Router       /report [post]
func (h handler) report(c *gin.Context) {
	var requestJSON api.ReportSchemaJSON
	var query string
	var reportData []reportRow

	if err := utils.CheckBinding(c, &requestJSON); err != nil {
		return
	}

	period := strings.Split(requestJSON.Period, "-")
	var year, month int
	var err error

	// checks
	if year, err = strconv.Atoi(period[0]); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Bad request": "wrong year format"})
		return
	}
	if year < 0 || year > 9999 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Bad request": "wrong year format"})
		return
	}
	if month, err = strconv.Atoi(period[1]); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Bad request": "wrong month format"})
		return
	}
	if month < 1 || month > 12 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Bad request": "wrong month format"})
		return
	}
	firstMonthTimestamp := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())
	lastMonthTimestamp := firstMonthTimestamp.AddDate(0, 1, 0).Add(-1 * time.Second)

	// get revenue info about from payment_history table (stores all transactions)
	query = "select o.service_id, sum(amount) from userbalancedb.payment_history p join userbalancedb.orders o on p.order_id = o.order_id where payment_date between $1 and $2 group by service_id order by service_id"
	var row reportRow
	_, err = h.DB.QueryFunc(context.Background(),
		query,
		[]interface{}{firstMonthTimestamp, lastMonthTimestamp},
		[]interface{}{&row.service_id, &row.amount},
		func(pgx.QueryFuncRow) error {
			// add each row data to the reportData slice
			reportData = append(reportData, row)
			return nil
		})
	// check if there is an error: no rows or smth else
	if err != nil {
		if err != pgx.ErrNoRows {
			log.Printf("Select failed: %s", err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
		} else {
			c.IndentedJSON(http.StatusNoContent, gin.H{"no content": "no data with requested period"})
		}
		return
	}

	var reportFileName string
	if reportFileName, err = h.writeReport(reportData); err != nil {
		log.Printf("Writing report failed: %s", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"internal server error": err.Error()})
	}

	// respone with the report file name and the file download link
	c.IndentedJSON(http.StatusCreated, api.ReportResponseJSON{
		ReportFilename: reportFileName,
		ReportFileLink: h.Socket + "/report/" + reportFileName,
	})
}

// @Summary      Download report
// @Description  download report file if OK.
// @Tags         report
// @Param        filename path string true "download file"
// @Produce      json
// @Success      200
// @Success      404 {object} api.ErrorResponseJSON
// @Router       /report/{filename} [get]
func (h handler) downloadReport(c *gin.Context) {
	filename := c.Param("filename")
	filepath := h.ReportDirPath + "/" + filename

	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.FileAttachment(filepath, "report.csv")
}
