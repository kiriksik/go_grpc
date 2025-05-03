package main

import (
	_ "github.com/lib/pq"
)

func main() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatalf("failed loading enviroment: %s", err)
	// }

	// apiCfg := config.InitializeApiConfig()

	// serveMux := handler.InitializeMux(apiCfg)
	// serverPort := os.Getenv("SERVER_PORT")
	// httpServer := http.Server{
	// 	Addr:    serverPort,
	// 	Handler: serveMux,
	// }
	// log.Println("server started on port", serverPort)
	// err = httpServer.ListenAndServe()
	// if err != nil {
	// 	log.Fatalf("Server failed: %s", err)
	// }

}
