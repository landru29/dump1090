// Package http is the http display.
package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/landru29/dump1090/internal/dump"
	"github.com/landru29/dump1090/internal/serialize"
)

const (
	cleanDelay time.Duration = time.Second * 10

	outOfDateAC time.Duration = time.Minute
)

// Transporter is the http transporter.
type Transporter struct {
	aircraftPool map[uint32]*dump.Aircraft
	serializer   serialize.Serializer
	mutex        sync.Mutex
}

func New(ctx context.Context, serializer serialize.Serializer, port int) (*Transporter, error) {
	output := Transporter{
		aircraftPool: make(map[uint32]*dump.Aircraft),
		serializer:   serializer,
	}

	router := mux.NewRouter()
	router.HandleFunc("/", output.serveData)

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
	}

	go func() {
		srv.ListenAndServe()
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
