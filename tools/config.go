package tools

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tanya-lyubimaya/WB_L0/internal/domain"
	"os"
)

var instance *domain.Config

func ConfigInstance() *domain.Config {
	if instance == nil {
		instance = new(domain.Config)
		err := godotenv.Load()
		if err != nil {
			panic("no .env file provided!")
		}
		instance.Host = os.Getenv("PGHOST")
		instance.Port = os.Getenv("PORT")
		instance.DBName = os.Getenv("DBNAME")
		instance.Username = os.Getenv("PGUSER")
		instance.Password = os.Getenv("PGPASSWORD")
		instance.NATSUrl = os.Getenv("NATS")
		fmt.Printf("%v\n", *instance)
	}
	return instance
}
