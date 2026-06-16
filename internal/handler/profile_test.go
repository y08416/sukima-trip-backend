package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupProfileRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := &ProfileHandler{}
	r.PUT("/profile", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		h.UpdateProfile(c)
	})
	return r
}

func TestUpdateProfile_Validation(t *testing.T) {
	r := setupProfileRouter()

	cases := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
		wantError  string
	}{
		{
			name:       "名前51文字はエラー",
			body:       map[string]interface{}{"name": strings.Repeat("あ", 51), "gender": "male"},
			wantStatus: http.StatusBadRequest,
			wantError:  "名前は50文字以内で入力してください",
		},
		{
			name:       "名前50文字はOK（repoがnilでも通過確認）",
			body:       map[string]interface{}{"name": strings.Repeat("あ", 50), "gender": "male"},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "不正なgender値はエラー",
			body:       map[string]interface{}{"name": "田中", "gender": "unknown"},
			wantStatus: http.StatusBadRequest,
			wantError:  "性別の値が不正です",
		},
		{
			name:       "genderが空はOK",
			body:       map[string]interface{}{"name": "田中", "gender": ""},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "genderがmaleはOK",
			body:       map[string]interface{}{"name": "田中", "gender": "male"},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, _ := json.Marshal(tc.body)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/profile", bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d (body: %s)", w.Code, tc.wantStatus, w.Body.String())
			}
			if tc.wantError != "" {
				var res map[string]string
				json.Unmarshal(w.Body.Bytes(), &res)
				if res["error"] != tc.wantError {
					t.Errorf("error: got %q, want %q", res["error"], tc.wantError)
				}
			}
		})
	}
}
