package handler

import "errors"

func validateNameAndGender(name, gender string) error {
	if len([]rune(name)) > 50 {
		return errors.New("名前は50文字以内で入力してください")
	}
	if gender != "" && gender != "male" && gender != "female" && gender != "other" {
		return errors.New("性別の値が不正です")
	}
	return nil
}
