package rtl28xxx

import (
	"context"
	"math/rand"
	"sync"
	"time"
	"unsafe"

	"github.com/landru29/dump1090/internal/source"
)

/*
  #include "context.h"
  #include <malloc.h>

*/
import "C"

var (
	contextMapper map[string]any
	random        *rand.Rand
	mutex         sync.Mutex
)

func init() {
	mutex.Lock()
	defer mutex.Unlock()

	contextMapper = map[string]any{}
	source := rand.NewSource(time.Now().UnixNano())
	random = rand.New(source)
}

type Context struct {
	Context  context.Context
	Key      string
	Ccontext unsafe.Pointer
}

// New creates a C Context with a processor
func New(ctx context.Context, processor source.Processer) *Context {
	if ctx == nil {
		ctx = context.Background()
	}

	output := &Context{
		Context: context.WithValue(ctx, deviceInContext{}, processor),
	}

	output.Key = saveContext(output)

	cKey := unsafe.Pointer(C.CString(string(output.Key)))

	output.Ccontext = unsafe.Pointer(C.newContext(cKey))

	return output
}

type deviceInContext struct{}

func saveContext(ctx context.Context) string {
	mutex.Lock()
	defer mutex.Unlock()

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[random.Intn(len(letterRunes))]
	}

	key := string(b)

	contextMapper[key] = ctx

	return key
}

// ContextFromKey gets a global context.
func ContextFromKey(key string) context.Context {
	return contextMapper[key].(context.Context)
}

// ContextFromPtr gets a context from a key pointer.
func ContextFromPtr(ptr unsafe.Pointer) context.Context {
	cContext := (*C.context)(ptr)
	contextPointer := (*C.char)(cContext.goContext)

	return ContextFromKey(C.GoString(contextPointer))
}

// Processor gets the processor from the context.
func Processor(ctx context.Context) source.Processer {
	return ctx.Value(deviceInContext{}).(source.Processer)
}

// DisposeContext Garbage collect the context.
func DisposeContext(key string) {
	delete(contextMapper, key)
}

// Deadline implements the context.Context interface.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
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
