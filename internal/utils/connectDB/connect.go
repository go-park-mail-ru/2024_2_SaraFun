package connectDB

import (
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/config"
)

func GetConnectURL(env config.EnvConfig) (string, error) {
	connStr := "host=" + env.DbHost + " port=" + env.DbPort + " user=" + env.DbUser + " password=" + env.DbPassword +
		" dbname=" + env.DbName + " sslmode=" + env.DbSSLMode
	//db, err := sql.Open("postgres", connStr)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if err = db.Ping(); err != nil {
	//	return nil, err
	//}
	//fmt.Println("Successfully connected to PostgreSQL!")
	//return db, nil
	return connStr, nil
}
