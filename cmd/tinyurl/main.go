package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cristalhq/acmd"
	"github.com/gorilla/mux"
	"github.com/ngharrington/tinyurl"
	"github.com/ngharrington/tinyurl/store"
	"github.com/rs/cors"
)

func runServer(ctx context.Context) {
	store, _ := store.NewSqliteUrlStore("./db.sqlite")
	app := tinyurl.NewApp("v0", store)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
	})
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tiny/{code}", tinyurl.MakeHandleFunction(app, tinyurl.HandleRedirect)).Methods("GET")
	router.HandleFunc("/tiny", tinyurl.MakeHandleFunction(app, tinyurl.HandleTinyfication)).Methods("POST")
	router.HandleFunc("/", tinyurl.MakeHandleFunction(app, tinyurl.HandleRoot)).Methods("GET")
	srv := http.Server{Addr: ":5000", Handler: c.Handler(router)}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe Error: %v", err)
		}
	}()
	<-ctx.Done()

	srv.Shutdown(ctx)

}

// func main() {
// store, _ := store.NewSqliteUrlStore("./db.sqlite")
// app := tinyurl.NewApp("v0", store)
// c := cors.New(cors.Options{
// 	AllowedOrigins: []string{"*"}, // All origins
// })
// router := mux.NewRouter().StrictSlash(true)
// router.HandleFunc("/tiny/{code}", tinyurl.MakeHandleFunction(app, tinyurl.HandleRedirect)).Methods("GET")
// router.HandleFunc("/tiny", tinyurl.MakeHandleFunction(app, tinyurl.HandleTinyfication)).Methods("POST")
// router.HandleFunc("/", tinyurl.MakeHandleFunction(app, tinyurl.HandleRoot)).Methods("GET")
// runServer(c.Handler(router))
// }

func main() {
	cmds := []acmd.Command{
		{
			Name:        "now",
			Description: "prints current time",
			ExecFunc: func(ctx context.Context, args []string) error {
				fmt.Printf("now: %s\n", time.Now().Format("15:04:05"))
				return nil
			},
		},
		{
			Name:        "status",
			Description: "prints status of the system",
			ExecFunc: func(ctx context.Context, args []string) error {
				// do something with ctx :)
				return nil
			},
		},
		{
			Name:        "sleep",
			Description: "prints status of the system",
			ExecFunc: func(ctx context.Context, args []string) error {
				log.Println("starting server")
				runServer(ctx)
				return nil
			},
		},
	}

	// all the acmd.Config fields are optional
	r := acmd.RunnerOf(cmds, acmd.Config{
		AppName:        "acmd-example",
		AppDescription: "Example of acmd package",
		Version:        "the best v0.x.y",
		// Context - if nil `signal.Notify` will be used
		// Args - if nil `os.Args[1:]` will be used
		// Usage - if nil default print will be used
	})

	if err := r.Run(); err != nil {
		r.Exit(err)
	}
}
