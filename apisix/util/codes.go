package util

import "google.golang.org/grpc/codes"

const (
	// Success status
	Success codes.Code = 200
	//SuccessCreated status
	SuccessCreated codes.Code = 201
	// SuccessNoContent status
	SuccessNoContent codes.Code = 204
	// InvalidArgument status
	InvalidArgument codes.Code = 400
	// Unauthorized status
	Unauthorized codes.Code = 401
	// Forbidden status
	Forbidden codes.Code = 403
	// NotFound status
	NotFound codes.Code = 404
	// Cancelled status
	Cancelled codes.Code = 405
	// RequestTimeout status
	RequestTimeout codes.Code = 408
	// InactiveAccount status
	InactiveAccount codes.Code = 410
	// InvalidToken status
	InvalidToken codes.Code = 411
	// InvalidAPIKey status
	InvalidAPIKey codes.Code = 412
	// InvalidSession status
	InvalidSession codes.Code = 413
	// ResourceExhausted status
	ResourceExhausted codes.Code = 414
	// InvalidSubdomain status
	InvalidSubdomain codes.Code = 420
	// InactiveSubdomain status
	InactiveSubdomain codes.Code = 421
	// SuspendedSubdomain status
	SuspendedSubdomain codes.Code = 422
	// InvalidTransaction status
	InvalidTransaction codes.Code = 430
	// DuplicateTransaction status
	DuplicateTransaction codes.Code = 431
	// InternalError status
	InternalError codes.Code = 500
	// ProcessingError status
	ProcessingError codes.Code = 502
)

// StatusMessage represent string message for code
var StatusMessage = map[codes.Code]string{
	Success:              "Berhasil",
	SuccessCreated:       "Berhasil, Data Tersimpan",
	SuccessNoContent:     "Berhasil, Data tidak ditemukan",
	InvalidArgument:      "Parameter tidak valid",
	Unauthorized:         "Username atau password kamu salah",
	Forbidden:            "Akses tidak dibolehkan atau kamu tidak memiliki akses",
	NotFound:             "Data tidak ditemukan",
	Cancelled:            "Permintaan dibatalkan",
	RequestTimeout:       "Permintaan melebihi batas waktu",
	InactiveAccount:      "Akun tidak aktif",
	InvalidToken:         "Akses tidak valid karena token tidak cocok atau login sudah kadaluarsa",
	InvalidAPIKey:        "API key tidak valid",
	InvalidSession:       "Sesi tidak valid atau sudah berakhir",
	ResourceExhausted:    "Sudah mencapai batas limit",
	InvalidSubdomain:     "Nama Toko kamu salah",
	InactiveSubdomain:    "Toko belum diaktifkan",
	SuspendedSubdomain:   "Toko kamu diblokir, silahkan hubungi Qasir melalui email hello@qasir.id",
	InvalidTransaction:   "Data transaksi tidak valid atau tidak sesuai",
	DuplicateTransaction: "Data transaksi duplikat",
	InternalError:        "Error dari server",
}
