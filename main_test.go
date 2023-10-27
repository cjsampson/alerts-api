package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlertsCreate(t *testing.T) {
	router := startRouter()
	w := httptest.NewRecorder()

	alert := Alert{
		ID:          "testing-id-123",
		ServiceID:   "testing-service-id",
		ServiceName: "testing-service-name",
		Model:       "testing-model",
		AlertType:   "info",
		AlertTS:     1695734400,
		Severity:    "low",
		TeamSlack:   "testing-slack-team",
	}

	alertJSON, err := json.Marshal(alert)
	if err != nil {
		t.Fail()
	}

	req, _ := http.NewRequest("POST", "/alerts", bytes.NewBuffer(alertJSON))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "testing-id-123", response["alert_id"])
	assert.Empty(t, response["error"])
}

func TestAlertsRead(t *testing.T) {
	router := startRouter()
	w := httptest.NewRecorder()

	alerts = []Alert{
		{
			ID:          "alert-id-1",
			ServiceID:   "testing-service",
			ServiceName: "testing-service-name-a",
			Model:       "testing-model-a",
			AlertType:   "Info",
			AlertTS:     1695734400,
			Severity:    "low",
			TeamSlack:   "testing-team-a",
		},
		{
			ID:          "alert-id-2",
			ServiceID:   "testing-service",
			ServiceName: "testing-service-b",
			Model:       "testing-model-b",
			AlertType:   "warn",
			AlertTS:     1695734400,
			Severity:    "medium",
			TeamSlack:   "testing-team-b",
		},
		{
			ID:          "alert-id-3",
			ServiceID:   "testing-service",
			ServiceName: "testing-service-c",
			Model:       "testing-model-c",
			AlertType:   "warn",
			AlertTS:     1695734400,
			Severity:    "high",
			TeamSlack:   "testing-team-c",
		},
	}

	req, err := http.NewRequest("GET", "/alerts?service_id=testing-service&start_ts=2022-10-20T00:00:00Z&end_ts=2023-10-27T00:00:00Z", nil)
	if err != nil {
		t.Fail()
	}
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		assert.Fail(t, "failed to retrieve an appropriate response")
	}
	assert.Nil(t, err)

	alertsList, ok := response["alerts"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, alertsList, 3)
}
