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

func TestArrive_CoordValidation(t *testing.T) {
	r := setupArriveRouter()

	cases := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
		wantError  string
	}{
		{
			name:       "lat > 90",
			body:       map[string]interface{}{"place_name": "test", "lat": 999, "lng": 135.0},
			wantStatus: http.StatusBadRequest,
			wantError:  "座標が不正です",
		},
		{
			name:       "lat < -90",
			body:       map[string]interface{}{"place_name": "test", "lat": -91, "lng": 135.0},
			wantStatus: http.StatusBadRequest,
			wantError:  "座標が不正です",
		},
		{
			name:       "lng > 180",
			body:       map[string]interface{}{"place_name": "test", "lat": 35.0, "lng": 181},
			wantStatus: http.StatusBadRequest,
			wantError:  "座標が不正です",
		},
		{
			name:       "lng < -180",
			body:       map[string]interface{}{"place_name": "test", "lat": 35.0, "lng": -181},
			wantStatus: http.StatusBadRequest,
			wantError:  "座標が不正です",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, _ := json.Marshal(tc.body)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/spots/dummy/arrive", bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", w.Code, tc.wantStatus)
			}
			var res map[string]string
			json.Unmarshal(w.Body.Bytes(), &res)
			if res["error"] != tc.wantError {
				t.Errorf("error: got %q, want %q", res["error"], tc.wantError)
			}
		})
	}
}
