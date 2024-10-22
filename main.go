package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	running      = flag.Bool("running", true, "If the application should be running. If set to false, the application will exit.")
	runningAfter = flag.Duration("running-after", 2*time.Second, "The duration after which the application will serve content.")
	healthy      = flag.Bool("healthy", true, "If the application should be healthy.")
	healthyAfter = flag.Duration("healthy-after", 10*time.Second, "The duration after which the application will serve 200 to the /health endpoint.")
	exitCode     = flag.Int("exit-code", 0, "The exit code of the application.")

	port = flag.String("port", "80", "Server listening port")

	startingTime = time.Now()
)

func init() {
	flag.Parse()
}

func main() {
	os.Exit(run())
}

func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if !*running {
		log.Printf("Application exiting because -running=false. Exit code is %d", *exitCode)
		return *exitCode
	}

	log.Printf("Application is starting... Should start in %.0f seconds.", runningAfter.Seconds())
	time.AfterFunc(*runningAfter, server)

	// Listen for the interrupt signal.
	<-ctx.Done()
	return *exitCode
}

func server() {
	mux := http.NewServeMux()
	mux.Handle("/health", withCLF(health))
	mux.Handle("/", withCLF(hello))

	log.Printf("Starting up on port %s (started in %.0f seconds)", *port, time.Since(startingTime).Seconds())

	log.Fatal(http.ListenAndServe(":"+*port, mux))
}

func withCLF(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next(w, r)

		// <remote_IP_address> - [<timestamp>] "<request_method> <request_path> <request_protocol>" -
		log.Printf("%s - - [%s] \"%s %s %s\" - -", r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.URL.Path, r.Proto)
	})
}

func hello(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := rw.Write([]byte("Mimic says hello!"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func health(rw http.ResponseWriter, _ *http.Request) {
	// Starting
	if *healthy && time.Since(startingTime) < *healthyAfter {
		_, err := rw.Write([]byte("starting"))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		return
	}

	// healthy
	if *healthy && time.Since(startingTime) > *healthyAfter {
		_, err := rw.Write([]byte("healthy"))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		return
	}

	// Unhealthy
	if !*healthy {
		_, err := rw.Write([]byte("unhealthy"))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusServiceUnavailable)
		return
	}
}
