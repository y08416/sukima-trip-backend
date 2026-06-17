package handler

import (
	"errors"
	"net/http"
)

// detectImageMime はバイト列からMIMEタイプを判定する。
// WebPはGo標準のDetectContentTypeが未対応のため個別にチェックする。
func detectImageMime(buf []byte) string {
	if len(buf) >= 12 &&
		buf[0] == 'R' && buf[1] == 'I' && buf[2] == 'F' && buf[3] == 'F' &&
		buf[8] == 'W' && buf[9] == 'E' && buf[10] == 'B' && buf[11] == 'P' {
		return "image/webp"
	}
	return http.DetectContentType(buf)
}

func validateNameAndGender(name, gender string) error {
	if len([]rune(name)) > 50 {
		return errors.New("名前は50文字以内で入力してください")
	}
	if gender != "" && gender != "male" && gender != "female" && gender != "other" {
		return errors.New("性別の値が不正です")
	}
	return nil
}
