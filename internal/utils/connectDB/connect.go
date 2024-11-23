package connectDB

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/config"
)

func ConnectDB(env config.EnvConfig) (*sql.DB, error) {
	//connStr := "host=sparkit-postgres port=5432 user=reufee password=sparkit dbname=sparkitDB sslmode=disable"
	//connStr := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname+ " sslmode=" + sslmode
	connStr := "host=" + env.DbHost + " port=" + env.DbPort + " user=" + env.DbUser + " password=" + env.DbPassword +
		" dbname=" + env.DbName + " sslmode=" + env.DbSSLMode
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to PostgreSQL!")
	return db, nil
}
