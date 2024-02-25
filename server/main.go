package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/pavva91/merkle-tree/server/config"
	"github.com/pavva91/merkle-tree/server/internal/middleware"
	"github.com/pavva91/merkle-tree/server/internal/router"

	// "github.com/swaggo/http-swagger" // http-swagger middleware
	_ "github.com/pavva91/merkle-tree/server/docs" // docs is generated by Swag CLI, you have to import it.
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//	@title			File Server with Merkle Tree
//	@version		1.0
//	@description	HTTP File Server with Merkle Tree to ensure that the downloaded data is not tampered

// @host	localhost:8080
func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	isDebug := false
	if len(os.Args) == 2 {
		debugArg := os.Args[1]
		if debugArg == "d" || debugArg == "debug" {
			os.Setenv("SERVER_ENVIRONMENT", "dev")
			isDebug = true
		}
	}
	log.Printf("debug mode: %t", isDebug)

	router.NewRouter()

	router.Router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	useEnvVar := os.Getenv("USE_ENVVAR")
	log.Printf("Using envvar value: %s\n", useEnvVar)
	log.Println("must be USE_ENVVAR=\"true\" to run with environment variable, otherwise will use a config file by default")

	if useEnvVar == "true" {
		config.Values.Server.Host = os.Getenv("SERVER_HOST")
		config.Values.Server.Port = os.Getenv("SERVER_PORT")
		bulkUploadMaxSize, err := strconv.Atoi(os.Getenv("MAX_BULK_UPLOAD_SIZE"))
		if err != nil {
			log.Panicf("insert int value for max size of bulk upload")
		}
		config.Values.Server.MaxBulkUploadSize = bulkUploadMaxSize
		fileUploadMaxSize, err := strconv.Atoi(os.Getenv("MAX_UPLOAD_FILE_SIZE"))
		if err != nil {
			log.Panicf("insert int value for max size of file upload")
		}
		config.Values.Server.MaxUploadFileSize = fileUploadMaxSize
		config.Values.Server.UploadFolder = os.Getenv("UPLOAD_FOLDER")
	} else {
		env := os.Getenv("SERVER_ENVIRONMENT")

		log.Printf("Running Environment: %s", env)

		switch env {
		case "dev":
			setConfig("./config/dev-config.yml")
			// setConfig("/home/bob/work/task/config/dev-config.yml")
		case "stage":
			log.Panicf("Incorrect Dev Environment: %s\nInterrupt execution", env)
		case "prod":
			log.Panicf("Incorrect Dev Environment: %s\nInterrupt execution", env)
		default:
			log.Panicf("Incorrect Dev Environment: %s\nRun with environment variable: SERVER_ENVIRONMENT=\"dev\" go run main.go\nInterrupt execution", env)
		}
	}

	// run the server
	fmt.Printf("Server is running on port %s\n", config.Values.Server.Port)
	// addr := fmt.Sprint("127.0.0.1:" + config.ServerConfigValues.Server.Port)
	// addr := fmt.Sprint("0.0.0.0:" + config.ServerConfigValues.Server.Port)
	addr := fmt.Sprint(config.Values.Server.Host + ":" + config.Values.Server.Port)

	srv := &http.Server{
		// Addr: "0.0.0.0:8080",
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		// Handler:      router.Router, // Pass our instance of gorilla/mux in.
		Handler: middleware.Limit(router.Router), // Pass instance of gorilla/mux with http reqeusts limiter
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Panicln("error shutting down gracefully, panic")
		return
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}

func setConfig(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config.Values)
	if err != nil {
		log.Fatal(err)
	}
}
