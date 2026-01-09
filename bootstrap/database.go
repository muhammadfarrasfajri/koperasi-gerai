package bootstrap

import "github.com/muhammadfarrasfajri/koperasi-gerai/database"

func InitDatabase() {
	database.ConnectMySQL()
}
