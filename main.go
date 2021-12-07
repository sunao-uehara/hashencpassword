package main

import (
	"context"
	"sync"
	"time"

	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	cmn "github.com/sunao-uehara/hashencpassword/common"
	r "github.com/sunao-uehara/hashencpassword/router"
)

func main() {
	// utilize multicore CPUs. enable this if go version is under 1.5
	// runtime.GOMAXPROCS(runtime.NumCPU())

	// router
	wg := &sync.WaitGroup{}
	mux := r.NewRouter(wg)

	srv := &http.Server{
		Addr:    ":" + cmn.PORT,
		Handler: mux,
	}
	// start up http server
	go func() {
		log.Println("listen and serve")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	// wait for SIGTERM signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	s := <-sig
	log.Printf("signal %s received\n", s)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("shut down the server gracefully...")
	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Fatalf("could not gracefully shut down server, %s", err.Error())
	}

	log.Println("waiting all goroutines are finished...")
	wg.Wait()
	log.Println("all done, really closing")
}
