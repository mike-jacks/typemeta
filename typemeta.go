package typemeta

import (
	"fmt"
	"reflect"
	"sync"
)

type Entry struct {
	TypeName string
	Key      string
	Value    string
}

var (
	mu       sync.RWMutex
	registry = make(map[reflect.Type]map[string]string)
)

func Register[T any](key, value string) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	mu.Lock()
	defer mu.Unlock()
	if registry[t] == nil {
		registry[t] = make(map[string]string)
	}
	registry[t][key] = value
}

func Meta[T any](key string) (string, bool) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	mu.RLock()
	defer mu.RUnlock()
	vals, ok := registry[t]
	if !ok {
		return "", false
	}
	val, ok := vals[key]
	return val, ok
}

func Must[T any](key string) string {
	if val, ok := Meta[T](key); ok {
		return val
	}
	panic(fmt.Sprintf("typemeta: key %q not found for type %T", key, *new(T)))
}

func MustWithLog[T any](key string) string {
	if val, ok := Meta[T](key); ok {
		return val
	}
	fmt.Printf("typemeta: missing metadata key %q for type %T\n", key, *new(T))
	panic("exiting due to missing typemeta")
}

func List() []Entry {
	mu.RLock()
	defer mu.RUnlock()
	var result []Entry
	for typ, pairs := range registry {
		for k, v := range pairs {
			result = append(result, Entry{
				TypeName: typ.String(),
				Key:      k,
				Value:    v,
			})
		}
	}
	return result
}

func ResetRegistry() {
	mu.Lock()
	defer mu.Unlock()
	registry = make(map[reflect.Type]map[string]string)
}
