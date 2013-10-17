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
	if !strings.Contains(b.String(), `ctx="module"] test`) {
		t.Errorf("expected ctx in log: %s", b.String())
	}
}

func TestVLogging(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewWriter("module", b)
	l.V(0).Info("test") // same as Info without V
	if !strings.Contains(b.String(), `ctx="module"] test`) {
		t.Errorf("expected ctx in log: %s", b.String())
	}
	l.V(1).Info("V") // higher verbosity not set
	if strings.Contains(b.String(), `ctx="module"] V`) {
		t.Errorf("unexpected verbose message in log: %s", b.String())
	}
}

func TestNewTrace(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewWriter("module", b).NewTrace("trace-id")
	l.Info("test")
	if !strings.Contains(b.String(), `ctx="module" trace="trace-id"] test`) {
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
