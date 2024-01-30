package config

import (
	//"fmt"
	"os"

	// "github.com/joho/godotenv"
	// "github.com/spf13/cast"
)
 
type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
}

func Load() Config {

	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println("error while loading godotenv!", err)
	// }
	
	cfg := Config{}

	// cfg.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	// cfg.PostgresPort = cast.ToString(getOrReturnDefault("POSTGRES_PORT", "5432"))
	// cfg.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "the_khoji"))
	// cfg.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRESS_PASSWORD", "546944"))
	// cfg.PostgresDB = cast.ToString(getOrReturnDefault("POSTGRES_DB", "mentor_store"))

	cfg.PostgresHost = "localhost"
	cfg.PostgresPort = "5432"
	cfg.PostgresUser = "the_khoji"
	cfg.PostgresPassword = "546944"
	cfg.PostgresDB = "mentor_store"

	return cfg
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}
