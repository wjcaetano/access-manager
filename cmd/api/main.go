package main

import "access-manager/cmd/api/modules"

func main() {
	app := modules.NewApp()
	app.Run()
}
