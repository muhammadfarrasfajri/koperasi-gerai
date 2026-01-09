package models

type BaseUser struct {
	ID        int    `json:"id"`
	GoogleUID string `json:"google_uid"`
	IDMember  string `json:"id_member"`
	// PERHATIKAN: Ada tambahan tag form:"..." di sebelah kanan
	Name         string `json:"name" form:"name"`
	Email        string `json:"email"` // Email biarkan kosong form-nya, karena dari Google
	NIK          string `json:"nik" form:"nik"`
	NPWP         string `json:"npwp" form:"npwp"`
	Gender       string `json:"gender" form:"gender"`
	Religion     string `json:"religion" form:"religion"`
	PlaceOfBirth string `json:"placeofbirth" form:"placeofbirth"`

	// Tanggal Lahir (String dulu biar aman)
	Birth string `json:"birth" form:"birth"`

	Address          string `json:"address" form:"address"`
	RegisterLocation string `json:"register_location" form:"register_location"`
	RegisterIP       string `json:"register_ip" form:"register_ip"`
	Job              string `json:"job" form:"job"`
	MaritalStatus    string `json:"marital_status" form:"marital_status"`
	Citizenship      string `json:"citizenship" form:"citizenship"`
	PhoneNumber      string `json:"phone_number" form:"phone_number"`

	GooglePicture  string `json:"google_picture"`
	ProfilePicture string `json:"profile_picture"` // Tidak perlu form tag, diisi backend
	KtpPicture     string `json:"ktp_image_path"`  // Tidak perlu form tag, diisi backend
}