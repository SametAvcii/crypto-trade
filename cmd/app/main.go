package main

// @title xxx API
// @version 1.0
// @description xxx API Documentation

// @host localhost:8000
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	/*ex, err := os.Executable()
	if err != nil {
		log.Fatalf("cannot get executable path: %v", err)
	}
	rootPath := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(ex)))) // /app/cmd/app -> /app
	if err := os.Chdir(rootPath); err != nil {
		log.Fatalf("failed to set working directory: %v", err)
	}*/
	StartApp()
}
