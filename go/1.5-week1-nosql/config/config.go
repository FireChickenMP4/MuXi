package config

import "os"

type Config struct {
	MongoURI   string
	MongoDB    string
	PgDSN      string
	ServerPort string
}

func Load() *Config {
	return &Config{
		MongoURI:   getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:    getEnv("MONGO_DB", "week1_nosql"),
		PgDSN:      getEnv("PG_DSN", "host=localhost user=postgres password=postgres dbname=week1_nosql port=5432 sslmode=disable"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
