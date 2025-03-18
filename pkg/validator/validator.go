package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	// Регулярные выражения для валидаций
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)
	passwordRegex = regexp.MustCompile(`^.{8,}$`) // Минимум 8 символов
)

func init() {
	validate = validator.New()

	// Регистрация кастомных валидаторов
	_ = validate.RegisterValidation("username", validateUsername)
	_ = validate.RegisterValidation("password", validatePassword)

	// Регистрация функции для получения имени поля из json-тега
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateStruct валидирует структуру с помощью тегов
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		// Преобразуем ошибки валидатора в более читаемый формат
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			messages := make([]string, 0, len(validationErrors))
			for _, e := range validationErrors {
				messages = append(messages, formatError(e))
			}
			return fmt.Errorf("validation error: %s", strings.Join(messages, "; "))
		}
		return err
	}
	return nil
}

// formatError форматирует ошибку валидации в удобочитаемое сообщение
func formatError(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()
	param := e.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, param)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "username":
		return fmt.Sprintf("%s must be between 3-20 characters and contain only letters, numbers, underscore or hyphen", field)
	case "password":
		return fmt.Sprintf("%s must be at least 8 characters long", field)
	default:
		return fmt.Sprintf("%s failed on %s validation", field, tag)
	}
}

// validateUsername проверяет, что имя пользователя соответствует требованиям
func validateUsername(fl validator.FieldLevel) bool {
	return usernameRegex.MatchString(fl.Field().String())
}

// validatePassword проверяет, что пароль соответствует требованиям
func validatePassword(fl validator.FieldLevel) bool {
	return passwordRegex.MatchString(fl.Field().String())
}
