package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("pengguna tidak ditemukan")
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrUserInactive       = errors.New("akun Anda telah dinonaktifkan. Hubungi administrator.")
	ErrUserLocked         = errors.New("akun terkunci karena terlalu banyak percobaan. Coba lagi dalam 30 menit.")
	ErrEmailRequired      = errors.New("email wajib diisi")
	ErrInvalidEmail       = errors.New("format email tidak valid")
	ErrPasswordRequired   = errors.New("password wajib diisi")
)
