package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// In-memory structure
var alerts []Alert

type Alert struct {
	ID          string `json:"alert_id"`
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	Model       string `json:"model"`
	AlertType   string `json:"alert_type"`
	AlertTS     int64  `json:"alert_ts"`
	Severity    string `json:"severity"`
	TeamSlack   string `json:"team_slack"`
}

type validatedQueryParams struct {
	ServiceID string
	StartTime time.Time
	EndTime   time.Time
}

func main() {
	router := startRouter()
	router.Run()
}

func startRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/alerts", alertsRead)
	router.POST("/alerts", alertsCreate)
	return router
}

// Create an alert and append to the in-memory data storage
func alertsCreate(c *gin.Context) {
	var alert Alert
	if err := c.ShouldBindJSON(&alert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"alert_id": alert.ID,
			"error":    err.Error(),
		})
		return
	}

	if alert.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "alert_id is required",
		})
		return
	}

	alerts = append(alerts, alert)

	c.JSON(http.StatusCreated, gin.H{
		"alert_id": alert.ID,
		"error":    "",
	})
}

// Query alert by service_id specifying by start_ts and end_ts
func alertsRead(c *gin.Context) {
	params, err := validateQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var matchedAlerts []Alert
	for _, alert := range alerts {
		if alert.ServiceID == params.ServiceID &&
			alert.AlertTS >= params.StartTime.Unix() &&
			alert.AlertTS <= params.EndTime.Unix() {
			matchedAlerts = append(matchedAlerts, alert)
		}
	}

	var serviceName string
	if len(matchedAlerts) > 0 {
		serviceName = matchedAlerts[0].ServiceName
	}

	c.JSON(http.StatusOK, gin.H{
		"service_id":   params.ServiceID,
		"service_name": serviceName,
		"alerts":       matchedAlerts,
	})
}

// Validate the query parameters - if any missing - return error
func validateQueryParameters(c *gin.Context) (validatedQueryParams, error) {
	var params validatedQueryParams
	var err error

	params.ServiceID = c.Query("service_id")
	startTS := c.Query("start_ts")
	endTS := c.Query("end_ts")

	if params.ServiceID == "" || startTS == "" || endTS == "" {
		return params, fmt.Errorf("missing parameters")
	}

	params.StartTime, err = time.Parse(time.RFC3339, startTS)
	if err != nil {
		return params, err
	}

	params.EndTime, err = time.Parse(time.RFC3339, endTS)
	if err != nil {
		return params, err
	}

	return params, nil
}
