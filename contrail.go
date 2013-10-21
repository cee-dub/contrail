// Copyright 2013 Cameron Walters. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package contrail implements thin wrappers around glog to add useful
// properties to logged messages.
//
// Example (assuming the date is 2006-01-02 15:04:05):
//	// Log messages with a context tag of "trackme"
//	log := contrail.New("trackme")
//	log.Warningln("I can log stuff, neat")
//	// Output: W0102 15:04:05.678901 threadid contrail.go:10 ctx="trackme"] I can log stuff, neat
//
// Log output is buffered and written periodically using Flush. Programs
// should call Flush before exiting to guarantee all log output is written.
//
// If you need to log to an io.Writer w instead of auto-rotated files, use
// contrail.NewWriter(w).
//
// Example of tracing:
//	// Follow a specific identifier as it is handled in some part of your code.
//	// It maintains the context of the parent.
//	tracer := log.NewTrace('some-identifier')
//	tracer.Errorln("unexpected error")
//	// Output: E0102 15:04:05.678901 threadid contrail.go:17 ctx="trackme" trace="some-identifier"] unexpected error
//
// The above examples also apply to any Logger returned by NewWriter.
package contrail

import (
	"bytes"
	"fmt"
	"io"

	"github.com/cee-dub/glog"
)

// Logger is the set of logging methods supported by all contrail logger types.
// See glog's documentation for more detail.
type Logger interface {
	// NewTrace returns a logger that includes trace="id" in logged messages.
	NewTrace(id string) *log

	// HeaderTag returns the logging tag in use by the Logger. Useful to integrate
	// with other logging systems so tagged output can be matched.
	HeaderTag() string

	// V is used to protect verbose logs from printing unless enabled.
	V(level glog.Level) glog.TagVerbose

	InfoLogger

	// These function signatures Copyright 2013 Google Inc. All Rights Reserved.
	Warning(args ...interface{})
	Warningln(args ...interface{})
	Warningf(format string, args ...interface{})
	Error(args ...interface{})
	Errorln(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalln(args ...interface{})
	Fatalf(format string, args ...interface{})
}

// InfoLogger is the subset of logging methods pertaining to more verbose logging.
type InfoLogger interface {
	// These function signatures Copyright 2013 Google Inc. All Rights Reserved.

	Info(args ...interface{})
	Infoln(args ...interface{})
	Infof(format string, args ...interface{})
}

// implementation of Logger interface, via embedded *glog.Tag
type log struct {
	name string
	*glog.Tag
}

// New returns a logger that includes ctx="name" in logged messages.
func New(name string) *log {
	return &log{
		name: name,
		Tag:  glog.NewTag(format(name, "")),
	}
}

// NewWriter returns a logger writing all messages to w that includes ctx="name"
// in logged messages.
func NewWriter(name string, w io.Writer) *log {
	return &log{
		name: name,
		Tag:  glog.NewTagWriter(format(name, ""), w),
	}
}

// NewTrace returns a logger based on l and its ctx name, that also includes
// trace="id" in logged messages.
func (l *log) NewTrace(id string) *log {
	return &log{
		name: l.name,
		Tag:  l.Tag.New(format(l.name, id)),
	}
}

// HeaderTag is part of the Logger interface.
func (l *log) HeaderTag() string {
	return l.Tag.String()
}

// always include ctx="ctx" and add trace="trace" if necessary.
func format(ctx, trace string) string {
	if trace == "" {
		return fmt.Sprintf("ctx=%q", ctx)
	}
	return fmt.Sprintf("ctx=%q trace=%q", ctx, trace)
}

// LogWriter adapts a log-level function (Infoln, Warninln, etc.) to io.Writer.
func LogWriter(logFunc func(args ...interface{})) *logWriter {
	return &logWriter{f: logFunc}
}

type logWriter struct {
	f func(args ...interface{})
}

func (l *logWriter) Write(p []byte) (int, error) {
	for _, line := range bytes.Split(p, []byte{'\n'}) {
		l.f(string(line))
	}
	return len(p), nil
}
