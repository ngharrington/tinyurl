package main

import (
	"context"
	"flag"
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

func runServer(ctx context.Context, cfg *tinyurl.ServerConfig) {
	store, _ := store.NewSqliteUrlStore("./db.sqlite")
	app := tinyurl.NewApp("v0", store)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
	})
	addr := cfg.Host + ":" + cfg.Port
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tiny/{code}", tinyurl.MakeHandleFunction(app, tinyurl.HandleRedirect)).Methods("GET")
	router.HandleFunc("/tiny", tinyurl.MakeHandleFunction(app, tinyurl.HandleTinyfication)).Methods("POST")
	router.HandleFunc("/", tinyurl.MakeHandleFunction(app, tinyurl.HandleRoot)).Methods("GET")
	srv := http.Server{Addr: addr, Handler: c.Handler(router)}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe Error: %v", err)
		}
	}()
	<-ctx.Done()

	srv.Shutdown(ctx)

}

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
			Name:        "server",
			Description: "runs the web service",
			ExecFunc: func(ctx context.Context, args []string) error {
				cfg := tinyurl.NewConfig()
				fs := flag.NewFlagSet("tinyurl", flag.PanicOnError)
				fs.StringVar(&cfg.Port, "port", "5000", "")
				fs.StringVar(&cfg.Host, "host", "localhost", "")
				if err := fs.Parse(args); err != nil {
					fmt.Println(err)
					return err
				}
				fmt.Println(cfg)
				runServer(ctx, cfg)
				return nil
			},
		},
	}

	// all the acmd.Config fields are optional
	r := acmd.RunnerOf(cmds, acmd.Config{
		AppName:        "tinyurl",
		AppDescription: "A tinyurl web service",
		Version:        "0.0.1",
	})

	if err := r.Run(); err != nil {
		fmt.Println(err)
		r.Exit(err)
	}
}
