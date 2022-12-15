package main

import (
	API "api"
	"env"
)


func main() {
	env.InitEnvVars()
	API.IPInfo()
}
