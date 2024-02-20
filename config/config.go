package config

import "fmt"

const (
	DBHost     = "localhost"
	DBPort     = 5432
	DBUser     = "user"
	DBPassword = "user"
	DBName     = "postgres"
)

// var ConnectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 	DBHost, DBPort, DBUser, DBPassword, DBName)

func GetConnectionString() (string, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName)
	return connectionString, nil
}
