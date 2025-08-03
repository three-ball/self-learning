package main

import "log"

func main() {
	cfg := config{
		addr: ":8082",
	}

	app := &application{
		config: cfg,
	}

	log.Fatal(app.run(app.mount()))
}
