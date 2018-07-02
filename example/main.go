package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/web-services/config"
	"github.com/linkernetworks/web-services/web"
)

func main() {

	web, err := web.New(&config.Config{
		Logger: logger.LoggerConfig{
			Level: "debug",
		},
	})
	if err != nil {
		logger.Fatalln("create web failed. err: [%v]", err)
	}

	http.HandleFunc("/signin", web.SignIn)
	http.HandleFunc("/signup", web.SignUp)

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			logger.Fatalln("Start HTTP server failed. err: [%v]", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
