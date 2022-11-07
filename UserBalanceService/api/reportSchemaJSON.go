package api

type ReportSchemaJSON struct {
	Period string `json:"period" binding:"required" example:"2022-10"`
}
