package main

func main() {
	InitLogger()
	config := NewConfig()

	storage := NewStorage(config.Storage)
	controller := NewFileController(storage)

	InitHttp(config.HTTP, controller)
}
