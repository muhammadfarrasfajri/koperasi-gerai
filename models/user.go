package models

type BaseUser struct {
	ID        int    `json:"id"`
	GoogleUID string `json:"google_uid"`

	// PERHATIKAN: Ada tambahan tag form:"..." di sebelah kanan
	Name         string `json:"name" form:"name"`
	Email        string `json:"email"` // Email biarkan kosong form-nya, karena dari Google
	Nik          string `json:"nik" form:"nik"`
	Npwp         string `json:"npwp" form:"npwp"`
	JenisKelamin string `json:"jenis_kelamin" form:"jenis_kelamin"`
	Agama        string `json:"agama" form:"agama"`
	TempatLahir  string `json:"tempat_lahir" form:"tempat_lahir"`

	// Tanggal Lahir (String dulu biar aman)
	TanggalLahir string `json:"tanggal_lahir" form:"tanggal_lahir"`

	AlamatDomisili   string `json:"alamat_domisili" form:"alamat_domisili"`
	RegisterLocation string `json:"register_location" form:"register_location"`
	RegisterIP       string `json:"register_ip" form:"register_ip"`
	Pekerjaan        string `json:"pekerjaan" form:"pekerjaan"`
	StatusPerkawinan string `json:"status_perkawinan" form:"status_perkawinan"`
	WargaNegara      string `json:"warga_negara" form:"warga_negara"`
	NoHp             string `json:"no_hp" form:"no_hp"`

	GooglePicture  string `json:"google_picture"`
	ProfilePicture string `json:"profile_picture"` // Tidak perlu form tag, diisi backend
	KtpImagePath   string `json:"ktp_image_path"`  // Tidak perlu form tag, diisi backend

	Role       string `json:"role"`
	IsLoggedIn int    `json:"is_logged_in"`
}