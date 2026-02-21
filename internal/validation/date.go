package validation

import (
	"reflect"
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// Формат даты подписки: MM-YYYY (месяц 01–12, год — 4 цифры).
var monthYearRegex = regexp.MustCompile(`^(0[1-9]|1[0-2])-\d{4}$`)

// MonthYear проверяет, что строка в формате MM-YYYY и месяц корректен (01–12).
// Поддерживает string и *string (для опциональных полей).
func MonthYear(fl validator.FieldLevel) bool {
	var s string
	switch fl.Field().Kind() {
	case reflect.String:
		s = fl.Field().String()
	case reflect.Ptr:
		if fl.Field().IsNil() {
			return true
		}
		s = fl.Field().Elem().String()
	default:
		return false
	}
	if s == "" {
		return true
	}
	if !monthYearRegex.MatchString(s) {
		return false
	}
	month, _ := strconv.Atoi(s[:2])
	return month >= 1 && month <= 12
}

// RegisterMonthYear регистрирует кастомный тег "month_year" в валидаторе.
func RegisterMonthYear(v *validator.Validate) error {
	return v.RegisterValidation("month_year", MonthYear)
}

// IsValidMonthYear возвращает true, если s в формате MM-YYYY (01–12 и 4 цифры года).
func IsValidMonthYear(s string) bool {
	if s == "" {
		return false
	}
	return monthYearRegex.MatchString(s) && func() bool {
		month, _ := strconv.Atoi(s[:2])
		return month >= 1 && month <= 12
	}()
}
