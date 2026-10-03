package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonitorCreateRequest_Validation(t *testing.T) {
	req := MonitorCreateRequest{
		Name:     stringPtr("Test"),
		Type:     stringPtr("http"),
		Interval: intPtr(60),
	}
	assert.NotNil(t, req.Name)
	assert.NotNil(t, req.Type)
	assert.NotNil(t, req.Interval)
}

func TestMonitorResponse_TagsEmpty(t *testing.T) {
	resp := MonitorResponse{
		ID:   1,
		Name: stringPtr("Test"),
		Tags: []MonitorTagResponse{},
	}
	assert.NotNil(t, resp.Tags)
	assert.Empty(t, resp.Tags)
}

func TestMonitorResponse_TagsWithData(t *testing.T) {
	resp := MonitorResponse{
		ID:   1,
		Name: stringPtr("Test"),
		Tags: []MonitorTagResponse{
			{ID: 1, MonitorID: 1, TagID: 1, Value: stringPtr("prod")},
		},
	}
	assert.Len(t, resp.Tags, 1)
	assert.Equal(t, "prod", *resp.Tags[0].Value)
}

func TestAPIResponse_Success(t *testing.T) {
	resp := APIResponse{
		Success: true,
		Message: "OK",
		Data:    map[string]string{"key": "value"},
	}
	assert.True(t, resp.Success)
	assert.Equal(t, "OK", resp.Message)
}

func TestAPIResponse_Error(t *testing.T) {
	resp := APIResponse{
		Success: false,
		Error:   "Something went wrong",
	}
	assert.False(t, resp.Success)
	assert.Equal(t, "Something went wrong", resp.Error)
}

func TestPaginatedResponse(t *testing.T) {
	resp := PaginatedResponse{
		Success: true,
		Data:    []MonitorResponse{{ID: 1}, {ID: 2}},
		Total:   2,
		Page:    1,
		Limit:   10,
	}
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, int64(2), resp.Total)
}

func TestHeartbeatResponse(t *testing.T) {
	resp := HeartbeatResponse{
		ID:        1,
		MonitorID: 1,
		Status:    0,
		Time:      "2024-01-01 12:00:00",
		Ping:      int64Ptr(50),
	}
	assert.Equal(t, uint(1), resp.ID)
	assert.NotNil(t, resp.Ping)
	assert.Equal(t, int64(50), *resp.Ping)
}

func TestTagResponse(t *testing.T) {
	resp := TagResponse{
		ID:          1,
		Name:        "test-tag",
		Color:       "#ff0000",
		CreatedDate: "2024-01-01 12:00:00",
	}
	assert.Equal(t, "test-tag", resp.Name)
	assert.Equal(t, "#ff0000", resp.Color)
}

func TestMaintenanceResponse(t *testing.T) {
	resp := MaintenanceResponse{
		ID:          1,
		Title:       "Test",
		Description: "Desc",
		Strategy:    "single",
		Active:      true,
	}
	assert.Equal(t, "Test", resp.Title)
	assert.True(t, resp.Active)
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}