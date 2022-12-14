package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var cmdConf string

func init() {
	if flag.Lookup("foo") == nil {
		flag.StringVar(&cmdConf, "config_path", ".env", "Set config path eg: .env")
	}
}

func initFlags() {
	cmdConf = flag.Lookup("config_path").Value.(flag.Getter).Get().(string)
}

func Config(key string) string {
	flag.Parse()
	initFlags()

	// Load .env file
	err := godotenv.Load(cmdConf)

	// Check if failed to load file
	if err != nil {
		fmt.Print("Error loading .env file")
	}

	// Take env file parameters
	return os.Getenv(key)
}
