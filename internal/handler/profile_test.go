package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
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

func setupAvatarRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := &ProfileHandler{}
	r.POST("/profile/avatar", func(c *gin.Context) {
		c.Set("user_id", "test-user")
		h.UploadAvatar(c)
	})
	return r
}

func TestUploadAvatar_MimeValidation(t *testing.T) {
	r := setupAvatarRouter()

	jpegBytes := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10}
	pngBytes := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	textBytes := []byte("this is not an image")

	cases := []struct {
		name       string
		content    []byte
		wantStatus int
		wantError  string
	}{
		{
			name:       "JPEGは通過",
			content:    jpegBytes,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "PNGは通過",
			content:    pngBytes,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "テキストは拒否",
			content:    textBytes,
			wantStatus: http.StatusBadRequest,
			wantError:  "jpg または png 形式の画像のみアップロードできます",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("avatar", "test.jpg")
			part.Write(tc.content)
			writer.Close()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/profile/avatar", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
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
