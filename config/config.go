package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     int
	SSLMode  string
}

func GetConfig() *Config {
	log.Printf("DBUSER: %s, DBPASSWORD: %s", os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"))
	return &Config{
		User:     os.Getenv("DBUSER"),
		Password: os.Getenv("DBPASSWORD"),
		DBName:   os.Getenv("DBNAME"),
		Host:     os.Getenv("DBHOST"),
		// Port: func() int {
		// 	p, err := strconv.Atoi(os.Getenv("PORT"))
		// 	if err != nil {
		// 		log.Panicf("port %+v", err)
		// 	}

		// 	return p
		// }(),
		Port: func() int {
			portStr := os.Getenv("PORT")
			if portStr == "" {
				portStr = "5432" // Default value
			}

			p, err := strconv.Atoi(portStr)
			if err != nil {
				log.Panicf("Invalid port value: %+v", err)
			}
			return p
		}(),
		SSLMode: os.Getenv("SSL"),
	}
}

func (c *Config) DBDataSource() string {
	log.Printf("postgres://host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
	log.Printf("postgres://%s:%s@%s:%d/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
