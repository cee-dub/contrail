package contrail

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func ExampleFormat_Category() {
	fmt.Println(format("module", ""))
	// Output: ctx="module"
}

func ExampleFormat_CategoryTrace() {
	fmt.Println(format("module", "trace-id"))
	// Output: ctx="module" trace="trace-id"
}

func TestNewWriter(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewWriter("module", b)
	l.Info("test")
	if !strings.Contains(b.String(), `ctx="module"]`) {
		t.Errorf("expected ctx in log: %s", b.String())
	}
}

func TestNewTrace(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewWriter("module", b).NewTrace("trace-id")
	l.Info("test")
	if !strings.Contains(b.String(), `ctx="module" trace="trace-id"]`) {
		t.Errorf("expected ctx and trace in log: %s", b.String())
	}
}

func TestCallerLocation(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewWriter("module", b)
	_, path, line, _ := runtime.Caller(0)
	l.Info("test")
	callsite := fmt.Sprintf("%s:%d", filepath.Base(path), line+1) // Caller one line above Info
	if !strings.Contains(b.String(), callsite) {
		t.Errorf("expected callsite %s in log: %s", callsite, b.String())
	}
}
