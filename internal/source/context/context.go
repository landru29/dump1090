// Package context is the C context.
package context

/*
  #include "context.h"
  #include <malloc.h>
*/
import "C"

import (
	"context"
	"math/rand"
	"sync"
	"time"
	"unsafe"

	"github.com/landru29/dump1090/internal/processor"
)

var (
	contextMapper map[string]any //nolint: gochecknoglobals
	random        *rand.Rand     //nolint: gochecknoglobals
	mutex         sync.Mutex     //nolint: gochecknoglobals
)

func init() { //nolint: gochecknoinits
	mutex.Lock()
	defer mutex.Unlock()

	contextMapper = map[string]any{}
	source := rand.NewSource(time.Now().UnixNano())
	random = rand.New(source)
}

// Context is the C context.
type Context struct {
	Context  context.Context //nolint: containedctx
	Key      string
	Ccontext unsafe.Pointer
}

const randomLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// New creates a C Context with a processor
func New(ctx context.Context, processor processor.Processer) *Context { //nolint: contextcheck
	if ctx == nil {
		ctx = context.Background()
	}

	output := &Context{
		Context: context.WithValue(ctx, deviceInContext{}, processor),
	}

	output.Key = saveContext(output)

	cKey := unsafe.Pointer(C.CString(output.Key))

	output.Ccontext = unsafe.Pointer(C.newContext(cKey)) //nolint: nlreturn

	return output
}

type deviceInContext struct{}

func saveContext(ctx context.Context) string {
	mutex.Lock()
	defer mutex.Unlock()

	letterRunes := []rune(randomLetters)
	b := make([]rune, 10) //nolint: gomnd
	for i := range b {
		b[i] = letterRunes[random.Intn(len(letterRunes))]
	}

	key := string(b)

	contextMapper[key] = ctx

	return key
}

// FromKey gets a global context.
func FromKey(key string) context.Context {
	return contextMapper[key].(context.Context) //nolint: forcetypeassert
}

// FromPtr gets a context from a key pointer.
func FromPtr(ptr unsafe.Pointer) context.Context {
	cContext := (*C.context)(ptr)
	contextPointer := (*C.char)(cContext.goContext)

	return FromKey(C.GoString(contextPointer))
} //nolint: ireturn,nolintlint

// Processor gets the processor from the context.
func Processor(ctx context.Context) processor.Processer { //nolint: ireturn
	return ctx.Value(deviceInContext{}).(processor.Processer) //nolint: forcetypeassert
}

// DisposeContext Garbage collect the context.
func DisposeContext(key string) {
	delete(contextMapper, key)
}

// Deadline implements the context.Context interface.
func (c *Context) Deadline() (time.Time, bool) {
	return c.Context.Deadline()
}

// Done implements the context.Context interface.
func (c *Context) Done() <-chan struct{} {
	return c.Context.Done()
}

// Err implements the context.Context interface.
func (c *Context) Err() error {
	return c.Context.Err()
}

// Value implements the context.Context interface.
func (c *Context) Value(key any) any {
	return c.Context.Value(key)
}
