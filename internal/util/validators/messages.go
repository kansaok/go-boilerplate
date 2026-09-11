package validators

import "github.com/kansaok/go-boilerplate/internal/util"

// ValidationMessages menyediakan pesan kustom untuk setiap kesalahan validasi
var ValidationMessages = map[string]string{
	"required":        	util.MESSAGES["ERROR_REQUIRED"],
	"email":           	util.MESSAGES["WRONG_EMAIL_FORMAT"],
	"unique":          	util.MESSAGES["UNIQUE"],
	"min8":             util.MESSAGES["MIN_8"],
	"eqfield": 			util.MESSAGES["PASSWORD_MISMATCH"],
	"password": 		util.MESSAGES["PASSWORD_CRITERIA"],
	"title": 			util.MESSAGES["TITLE_CRITERIA"],
	"gender": 			util.MESSAGES["GENDER_CRITERIA"],
	"date_format": 		util.MESSAGES["INVALID_DATE"],
}
