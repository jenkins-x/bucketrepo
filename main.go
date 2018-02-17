package main

import (
	"github.com/spf13/viper"
)

func main() {
	InitConfig()
	InitLogger()

	storageType := viper.GetString("storage.type")
	storage := NewStorage(storageType)
	controller := NewFileController(storage)

	InitHttp(controller)
}
