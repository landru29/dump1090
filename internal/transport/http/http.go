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
	mutex        sync.Mutex
	formaters    map[string]serialize.Serializer
}

func New(ctx context.Context, addr string, apiPath string, formaters []serialize.Serializer) (*Transporter, error) {
	subFS, _ := fs.Sub(staticFiles, "public")

	output := Transporter{
		aircraftPool: make(map[uint32]*dump.Aircraft),
		formaters:    map[string]serialize.Serializer{},
	}

	for _, elt := range formaters {
		output.formaters[elt.MimeType()] = elt
	}

	router := mux.NewRouter()
	router.HandleFunc(apiPath, output.serveData)
	router.HandleFunc("/config.js", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("window.apiPath='%s'", apiPath)))
	})
	if apiPath != "/" {
		router.PathPrefix("/").Handler(http.FileServer(http.FS(subFS)))
	}

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
	requestedMimeType := req.Header.Get("Accept")

	formater, ok := t.formaters[requestedMimeType]
	if !ok {
		formater = t.formaters["application/json"]
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	dataArray := []*dump.Aircraft{}
	for _, elt := range t.aircraftPool {
		dataArray = append(dataArray, elt)
	}

	output, err := formater.Serialize(dataArray)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}

	writer.Header().Set("content-type", requestedMimeType)
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
