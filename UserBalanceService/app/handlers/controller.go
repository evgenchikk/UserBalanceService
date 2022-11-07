package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type handler struct {
	Socket        string
	DB            *pgx.Conn
	ReportDirPath string
}

func RegisterRoutes(r *gin.Engine, db *pgx.Conn, reportDirPath string, socket string) {
	h := &handler{
		Socket:        socket,
		DB:            db,
		ReportDirPath: reportDirPath,
	}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/add", h.add)
	r.POST("/transfer", h.transfer)
	r.POST("/reserve", h.reserve)
	r.POST("/dereserve", h.dereserve)
	r.POST("/approve", h.approve)
	r.POST("/balance", h.balance)
	r.POST("/report", h.report)
	r.GET("/report/:filename", h.downloadReport)
}
