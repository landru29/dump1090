// Package http is the http display.
package http

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/landru29/dump1090/internal/dump"
	"github.com/landru29/dump1090/internal/serialize"
)

//go:embed public/*
var staticFiles embed.FS

const (
	cleanDelay time.Duration = time.Second * 10

	outOfDateAC time.Duration = time.Minute
)

// Transporter is the http transporter.
type Transporter struct {
	aircraftPool map[uint32]*dump.Aircraft
	serializer   serialize.Serializer
	mutex        sync.Mutex

	staticHandler http.Handler
}

func New(ctx context.Context, serializer serialize.Serializer, addr string) (*Transporter, error) {
	subFS, _ := fs.Sub(staticFiles, "public")

	output := Transporter{
		aircraftPool: make(map[uint32]*dump.Aircraft),
		serializer:   serializer,

		staticHandler: http.FileServer(http.FS(subFS)),
	}

	router := mux.NewRouter()
	router.HandleFunc("/", output.serveData)

	srv := &http.Server{
		Handler: router,
		Addr:    addr,
	}

	go func() {
		fmt.Printf("Serving on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("ERR: %s", err)
		}
	}()

	go func() {
		<-ctx.Done()
		srv.Shutdown(ctx)
	}()

	go func(app *Transporter) {
		app.acCleaner(ctx)
	}(&output)

	return &output, nil
}

// Transport implements the transport.Transporter interface.
func (t *Transporter) Transport(ac *dump.Aircraft) error {
	t.mutex.Lock()
	t.aircraftPool[ac.Addr] = ac
	t.mutex.Unlock()

	return nil
}

func (t *Transporter) serveData(writer http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Accept") != t.serializer.MimeType() {
		t.staticHandler.ServeHTTP(writer, req)
		return
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	output, err := t.serializer.Serialize(t.aircraftPool)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}

	writer.Header().Set("content-type", t.serializer.MimeType())
	writer.Write([]byte(output))
}

func (t *Transporter) acCleaner(ctx context.Context) {
	ticker := time.NewTicker(cleanDelay)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.mutex.Lock()
			for idx, ac := range t.aircraftPool {
				if time.Since(ac.Seen) > outOfDateAC {
					delete(t.aircraftPool, idx)
				}
			}
			t.mutex.Unlock()
		}
	}
}
