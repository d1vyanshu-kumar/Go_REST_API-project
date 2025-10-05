package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// and the http_server part is nested so we are going to create a another struct for that
type HTTPServer struct { // private struct
	Addr string    `yaml:"address" env-required:"true"` // yaml tag is used to map the yaml key to struct field
}
// env-default:"production
// struct tax
type Config struct { // public struct
	Env        string      `yaml:"env" env:"ENV" env-required:"true"` // yaml tag is used to map the yaml key to struct field
	StoragePath string      `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

func MustLoad() *Config { // as the name sugests we will need to must load this func if there is any error will be here then we dont have start the application
	// lets get the config path

	var configPath string

	configPath = os.Getenv("CONFIG_PATH") // we are getting the config path from the environment variable

	if configPath == "" { // if the config path is not set then we will use the default path
		// and if the config path is empty then we will checking the flag if there is any config path provided in the command line argument
		// what is flag and argument when we try to run through terminal we can provide some argument to the command like go run main.go --config=config.yaml here --config is the flag and config.yaml is the argument
		flags := flag.String("config", "", "path to config file") // we are using flag package to get the config path from the command line argument
		flag.Parse()

		configPath = *flags // dereferencing the pointer to get the value

		// and after all of that if the config path is still empty then here we go:

		if configPath == "" {
			log.Fatal("config path not provided") // if the config path is not provided then we will log the error and exit the application
		}
	}

	// after that if condition we will now checking the config path which we will set is there any file exists in that path or not

	if _, err := os.Stat(configPath); os.IsNotExist(err) { // os.Stat is used to check if the file exists or not

		log.Fatalf("config file does not exist: %s", configPath) // if the file does not exist then we will log the error and exit the application
	}

	var cfg Config // creating a variable of type Config struct

	err := cleanenv.ReadConfig(configPath, &cfg) // reading the config file and unmarshalling it into the config struct

	if err != nil {
		log.Fatalf("failed to read config file: %s, error: %v", configPath, err.Error()) // if there is any error then we will log the error and exit the application
	}

	return &cfg // returning the config struct

}