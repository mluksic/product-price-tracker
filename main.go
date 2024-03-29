package main

import (
	"fmt"
	"github.com/mluksic/product-price-tracker/cmd"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("unable to get working dir %s", err))
	}
	viper.SetConfigFile(dir + "/.env")
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// initialize commands
	cmd.Execute()
}
