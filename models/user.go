package models

type BaseUser struct {
	ID        int    `json:"id"`         //auto_increment
	IDMember  string `json:"id_member"`  //generate
	GoogleUID string `json:"google_uid"` //login user

	// PERHATIKAN: Ada tambahan tag form:"..." di sebelah kanan
	Name  string `json:"name" form:"name"` //input User
	Email string `json:"email"`            // Email biarkan kosong form-nya, karena dari Google
	NPWP  string `json:"npwp" form:"npwp"` // input user (optional)
	NIK   string `json:"nik" form:"nik"`   // input user (mandatory)

	PlaceOfBirth string `json:"place_of_birth" form:"place_of_birth"` // input user (mandatory)
	Birth        string `json:"birth" form:"birth"`               // input user (mandatory)
	Gender       string `json:"gender" form:"gender"`             // input user (mandatory)
	Address      string `json:"address" form:"address"`           // input user (mandatory)
	PosCode      string `json:"pos_code" form:"pos_code"`         // input user (mandatory)
	Religion     string `json:"religion" form:"religion"`         // input user (mandatory)

	MaritalStatus    string `json:"marital_status" form:"marital_status"`       // input user (mandatory)
	Job              string `json:"job" form:"job"`                             // input user (mandatory)
	Citizenship      string `json:"citizenship" form:"citizenship"`             // input user (mandatory)
	Blood_type       string `json:"blood_type" form:"blood_type"`               // input user (mandatory)
	PhoneNumber      string `json:"phone_number" form:"phone_number"`           // input user (mandatory)
	Is_verified      int    `json:"is_verified"`                                // verification by admin
	Rejected_reason  string `json:"rejected_reason"`                            // input from admin if rejected
	RegisterLocation string `json:"register_location" form:"register_location"` // input user (access location permission)
	RegisterIP       string `json:"register_ip" form:"register_ip"`             // input from code

	KtpPicture     string `json:"ktp_picture" binding:"-"`              // input from user (mandatory)
	ProfilePicture string `json:"profile_picture" binding:"-"`          // input  from user (mandatory)
	GooglePicture  string `json:"google_picture"`                       // input from google account
	LastEducation  string `json:"last_education" form:"last_education"` // input user (mandatory)
	ActiveAs       string `json:"active_as" form:"active_as"`           // automatic
	Mother_name    string `json:"mother_name" form:"mother_name"`       // input from user (mandatory)
}