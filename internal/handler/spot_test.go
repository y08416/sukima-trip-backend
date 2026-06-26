package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupArriveRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := &SpotHandler{}
	r.POST("/spots/:id/arrive", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		h.Arrive(c)
	})
	return r
}

func TestArrive_MissingPlaceName(t *testing.T) {
	r := setupArriveRouter()

	b, _ := json.Marshal(map[string]interface{}{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/spots/dummy/arrive", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
	var res map[string]string
	json.Unmarshal(w.Body.Bytes(), &res)
	if res["error"] != "リクエストが不正です" {
		t.Errorf("error: got %q, want %q", res["error"], "リクエストが不正です")
	}
}
