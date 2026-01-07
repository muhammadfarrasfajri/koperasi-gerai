package bootstrap

import (
	"github.com/muhammadfarrasfajri/login-google/database"
)

func InitDatabase() {
	database.ConnectMySQL()
}
