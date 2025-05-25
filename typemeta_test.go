package typemeta

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_typemeta(t *testing.T) {
	type User struct {
		name string
		age  int
	}
	t.Run("Register Success", func(t *testing.T) {
		// Clear the registry before testing
		mu.Lock()
		registry = make(map[reflect.Type]map[string]string)
		mu.Unlock()

		Register[User]("table", "user")
		Register[User]("plural", "users")

		userType := reflect.TypeOf((*User)(nil)).Elem()

		mu.RLock()
		defer mu.RUnlock()

		if got := registry[userType]["table"]; got != "user" {
			t.Errorf("expected '%s' for key '%s', got '%s'", "user", "table", got)
		}

		if got := registry[userType]["plural"]; got != "users" {
			t.Errorf("expected '%s' for key '%s', got '%s'", "users", "plural", got)
		}

	})
}

func ExampleRegister() {
	type User struct {
		name string
		age  int
	}

	mu.Lock()
	registry = make(map[reflect.Type]map[string]string)
	mu.Unlock()
	Register[User]("table", "user")
	userType := reflect.TypeOf((*User)(nil)).Elem()
	fmt.Println(registry[userType]["table"])

	// Output: user
}

func Benchmark_Register(b *testing.B) {
	type User struct {
		name string
		age  int
	}
	mu.Lock()
	registry = make(map[reflect.Type]map[string]string)
	mu.Unlock()
	for b.Loop() {
		Register[User]("table", "user")
	}
}
