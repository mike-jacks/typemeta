package typemeta

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func resetRegistry() {
	mu.Lock()
	defer mu.Unlock()
	registry = make(map[reflect.Type]map[string]string)
}

type User struct {
	Name string
	Age  int
}

var userType = reflect.TypeOf((*User)(nil)).Elem()

func Test_typemeta(t *testing.T) {
	t.Run("Register Success", func(t *testing.T) {
		resetRegistry()

		Register[User]("table", "user")
		Register[User]("plural", "users")

		mu.RLock()
		defer mu.RUnlock()

		if got := registry[userType]["table"]; got != "user" {
			t.Errorf("expected '%s' for key '%s', got '%s'", "user", "table", got)
		}

		if got := registry[userType]["plural"]; got != "users" {
			t.Errorf("expected '%s' for key '%s', got '%s'", "users", "plural", got)
		}

	})
	t.Run("Meta Success", func(t *testing.T) {
		resetRegistry()
		Register[User]("table", "user")
		value, ok := Meta[User]("table")
		if !ok {
			t.Errorf("expected ok to be: '%t'", true)
		}
		if value != "user" {
			t.Errorf("expected value to be: `%s`", "user")
		}
	})
}

func ExampleRegister() {
	resetRegistry()
	Register[User]("table", "user")
	fmt.Println(registry[userType]["table"])

	// Output: user
}

func ExampleMeta_success() {
	resetRegistry()
	Register[User]("table", "user")
	value, ok := Meta[User]("table")
	fmt.Println(value, ok)

	// Output: user true
}

func ExampleMeta_fail() {
	resetRegistry()
	Register[User]("table", "user")
	value, ok := Meta[User]("version")
	fmt.Println(value, ok)
	// Output: false
}

func ExampleMust_success() {
	resetRegistry()
	Register[User]("table", "user")
	val := Must[User]("table")
	fmt.Println(val)
	// Output: user

}
func ExampleMust_fail() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic")
		}
	}()
	resetRegistry()
	Register[User]("table", "user")
	Must[User]("missing")
	// Output: panic
}

func ExampleMustWithLog_success() {
	resetRegistry()
	Register[User]("table", "user")
	val := Must[User]("table")
	fmt.Println(val)
	// Output: user

}
func ExampleMustWithLog_fail() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic")
		}
	}()
	resetRegistry()
	Register[User]("table", "user")
	MustWithLog[User]("missing")
	// Output:
	// typemeta: missing metadata key "missing" for type typemeta.User
	// panic
}

func ExampleList() {
	resetRegistry()
	Register[User]("table", "user")
	Register[User]("plural", "users")
	results := List()
	slices.SortFunc(results, func(a, b Entry) int {
		return cmp.Compare(a.Key, b.Key)
	})
	fmt.Println(results)
	// Output: [{typemeta.User plural users} {typemeta.User table user}]
}

func Benchmark_Register(b *testing.B) {
	resetRegistry()
	for b.Loop() {
		Register[User]("table", "user")
	}
}

func Benchmark_Meta(b *testing.B) {
	resetRegistry()
	Register[User]("table", "user")
	for b.Loop() {
		Meta[User]("table")
	}
}

func Benchmark_Must(b *testing.B) {
	resetRegistry()
	Register[User]("table", "user")
	for b.Loop() {
		Must[User]("table")
	}
}

func Benchmark_MustWithLog(b *testing.B) {
	resetRegistry()
	Register[User]("table", "user")
	for b.Loop() {
		MustWithLog[User]("table")
	}
}

func Benchmark_List(b *testing.B) {
	resetRegistry()
	Register[User]("table", "user")
	Register[User]("plural", "users")
	for b.Loop() {
		List()
	}
}

func TestMust_panicsOnMissingKey(t *testing.T) {
	resetRegistry()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Must did not panic on missing key")
		}
	}()
	Must[User]("missing")
}

func TestMustWithLog_panicsOnMissingKeyAndLogs(t *testing.T) {
	resetRegistry()
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = old
	}()
	defer func() {
		w.Close()
	}()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustWithLog did not panic on missing key")
		}
		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		if !strings.Contains(buf.String(), "typemeta: missing metadata key") {
			t.Errorf("Expected log output for missing key, got: %s", buf.String())
		}
	}()
	MustWithLog[User]("missing")
}
