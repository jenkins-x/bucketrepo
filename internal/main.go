package main

func main() {
	InitLogger()
	config := NewConfig()

	storage := NewStorage(config.Storage)
	repository := NewRepository(config.Repository)

	controller := NewFileController(storage, repository)

	InitHttp(config.HTTP, controller)
}
