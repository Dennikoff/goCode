package model

import validation "github.com/go-ozzo/ozzo-validation"

func passwordRequiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond == true {
			return validation.Validate(value, validation.Required)
		}
		return nil
	}
}
