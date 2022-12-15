package main

import (
	API "api"
	"env"
)


// @title           IP Address API
// @version         1.0

// @contact.name   API Support
// @contact.url    http://t.me/onorridg

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	env.InitEnvVars()
	API.IPInfo()
}
