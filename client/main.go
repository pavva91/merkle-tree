/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/pavva91/merkle-tree/client/cmd"
	"github.com/spf13/viper"
)

func main() {

	viper.AutomaticEnv()
	viper.SetDefault("DEFAULT_STORAGE_FOLDER", "./storage")
	viper.SetDefault("DEFAULT_UPLOAD_FOLDER", "./testfiles")
	viper.SetDefault("DEFAULT_DOWNLOAD_FOLDER", "./downloads")
	viper.SetDefault("SERVER_URL", "http://localhost:8080")

	cmd.Execute()
}
