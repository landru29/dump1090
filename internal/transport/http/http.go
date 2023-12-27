// Package http is the http display.
package http

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/landru29/dump1090/internal/database"
	"github.com/landru29/dump1090/internal/model"
	"github.com/landru29/dump1090/internal/serialize"
)

//go:embed public/*
var staticFiles embed.FS

const (
	readHeaderTimeout time.Duration = time.Second * 10
)

// Transporter is the http transporter.
type Transporter struct {
	aircraftDB *database.ElementStorage[model.ICAOAddr, model.Aircraft]
	formaters  map[string]serialize.Serializer
}

// New creates an http transporter.
func New(
	ctx context.Context,
	addr string,
	apiPath string,
	aircraftDB *database.ElementStorage[model.ICAOAddr, model.Aircraft],
	formaters []serialize.Serializer,
) (*Transporter, error) {
	subFS, _ := fs.Sub(staticFiles, "public")

	output := Transporter{
		aircraftDB: aircraftDB,
		formaters:  map[string]serialize.Serializer{},
	}

	for _, elt := range formaters {
		output.formaters[elt.MimeType()] = elt
	}

	router := mux.NewRouter()

	router.HandleFunc(apiPath, output.serveData)

	router.HandleFunc("/config.js", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf("window.apiPath='%s'", apiPath)))
	})

	if apiPath != "/" {
		router.PathPrefix("/").Handler(http.FileServer(http.FS(subFS)))
	}

	srv := &http.Server{
		Handler:           router,
		Addr:              addr,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		fmt.Printf("Serving on %s\n", addr) //nolint: forbidigo

		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("ERR: %s", err) //nolint: forbidigo
		}
	}()

	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(ctx)
	}()

	return &output, nil
}

// Transport implements the transport.Transporter interface.
func (t *Transporter) Transport(_ *model.Aircraft) error {
	return nil
}

func (t *Transporter) serveData(writer http.ResponseWriter, req *http.Request) {
	requestedMimeType := req.Header.Get("Accept")

	formater, ok := t.formaters[requestedMimeType]
	if !ok {
		formater = t.formaters["application/json"]
	}

	dataArray := []*model.Aircraft{}
	for _, addr := range t.aircraftDB.Keys() {
		dataArray = append(dataArray, t.aircraftDB.Element(addr))
	}

	output, err := formater.Serialize(dataArray)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}

	writer.Header().Set("content-type", requestedMimeType)
	_, _ = writer.Write(output)
}

// String implements the transport.Transporter interface.
func (t *Transporter) String() string {
	return "http"
}
