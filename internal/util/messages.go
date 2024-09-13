package util

var MESSAGES = map[string]string{
    "GET_SUCCESS":         	"GET data berhasil",
    "REGISTER_SUCCESS":		"Register data berhasil",
    "REGISTER_FAILED":		"Register data gagal",
    "METHOD_NOT_ALLOWED":  	"Method tidak diizinkan",
    "HOST_NOT_ALLOWED":  	"Host tidak diizinkan",
    "SAVE_SUCCESS":        	"Input data berhasil",
    "UPDATE_SUCCESS":      	"Ubah data berhasil",
    "DELETE_SUCCESS":      	"Hapus data berhasil",
    "LOGIN_FAILED":        	"Email atau password salah",
    "LOGIN_SUCCESS":       	"Login berhasil",
    "LOGOUT_SUCCESS":      	"Logout berhasil",
    "MIN_8":               	"Kolom ini harus minimal 8 karakter",
    "ERROR":               	"Kesalahan Sistem",
    "ERROR_REQUIRED":      	"Kolom ini harus diisi",
    "DATA_NOT_FOUND":      	"Data tidak ditemukan",
    "ERROR_VALIDATION":    	"Validasi gagal",
    "EMAIL_USED":          	"Email sudah digunakan",
    "WRONG_EMAIL_FORMAT":	"Format email salah !",
    "UNIQUE":           	"Data sudah ada di database",
    "PASSWORD_MISMATCH":   	"Password dan konfirmasi password tidak cocok",
    "PASSWORD_MIN_8":      	"Password minimal 8 karakter",
    "INVALID_DATE":        	"Format tanggal salah. gunakan format YYY-mm-dd",
    "TOKEN_NOTFOUND":       "Token tidak ditemukan",
    "INVALID_TOKEN":       	"Token tidak valid",
    "PROCESS_SUCCESS":     	"Proses Berhasil",
    "APIKEY_REQUIRED":     	"API Key harus diisi",
    "APIKEY_INVALID":      	"API Key tidak valid",
    "INVALID_PAYLOAD":      "Invalid request payload",
    "PASSWORD_CRITERIA":   	"Password harus terdiri dari minimal 1 huruf besar, 1 huruf kecil, 1 angka, dan 1 spesial karakter",
    "TITLE_CRITERIA":   	"Title salah. gunakan code: Mr, Mrs, or Miss",
    "GENDER_CRITERIA":   	"Gender salah. gunakan code: L for Man or P for Woman",
}
