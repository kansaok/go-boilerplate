package validators

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/kansaok/go-boilerplate/internal/db"
	"gorm.io/gorm"
)

var Validate *validator.Validate

// UniqueValidation adalah validasi custom untuk mengecek field unik di database
func UniqueValidation(fl validator.FieldLevel) bool {
    params := strings.Split(fl.Param(), "=")
    if len(params) != 2 {
        return false
    }

    table := params[0]
    field := params[1]
    value := fl.Field().String()

	dbConn, ok := db.GetDB().(*gorm.DB)
    if !ok {
        return false
    }

    var count int64
    result := dbConn.Table(table).Where(fmt.Sprintf("%s = ?", field), value).Count(&count)
    if result.Error != nil {
        return false
    }

    return count == 0
}

func Min8Validation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return len(value) >= 8
}

func GetJSONFieldName(s interface{}, fieldName string) string {
	r := reflect.TypeOf(s)
	for i := 0; i < r.NumField(); i++ {
		field := r.Field(i)
		if field.Name == fieldName {
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" || jsonTag == "-" {
				return field.Name
			}
			return strings.Split(jsonTag, ",")[0]
		}
	}
	return fieldName
}

func ValidateAndMapErrors(v *validator.Validate, req interface{}) map[string]string {
	errors := make(map[string]string)

	// Validasi struct
	if err := v.Struct(req); err != nil {
		for _, validationErr := range err.(validator.ValidationErrors) {
			// Ambil nama field berdasarkan tag JSON
			jsonFieldName := GetJSONFieldName(req, validationErr.Field())
			tag := validationErr.Tag()
			// Isi pesan error
			errors[jsonFieldName] = ValidationMessages[tag]
		}
	}
	return errors
}


func ValidatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()

    if len(password) < 8 {
        return false
    }
    if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
        return false
    }
    if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
        return false
    }
    if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
        return false
    }
    if matched, _ := regexp.MatchString(`[!@#\$%\^&\*]`, password); !matched {
        return false
    }
    return true
}

// Validasi email
func ValidateEmail(email string) error {
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !re.MatchString(email) {
        return errors.New("format email tidak valid")
    }
    return nil
}

func ValidateEnum(enumValues []string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		for _, v := range enumValues {
			if value == v {
				return true
			}
		}
		return false
	}
}

func ValidateDateFormat(fl validator.FieldLevel) bool {
	// Regular expression for YYYY-dd-mm format
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return re.MatchString(fl.Field().String())
}

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("unique", UniqueValidation)
	v.RegisterValidation("min8", Min8Validation)
	v.RegisterValidation("password", ValidatePassword)
	v.RegisterValidation("title", ValidateEnum([]string{"Mr", "Mrs", "Miss"}))
	v.RegisterValidation("gender", ValidateEnum([]string{"L", "P"}))
	v.RegisterValidation("date_format", ValidateDateFormat)
}

func LoadValidatorConfig() *validator.Validate {
    Validate = validator.New()
    RegisterCustomValidators(Validate)
    return Validate
}
