// Copyright 2013 Cameron Walters. All Rights Reserved.

// Package contrail implements thin wrappers around glog to add useful
// properties to logged messages.
package contrail

import (
	"fmt"
	"io"

	"github.com/cee-dub/glog"
)

// Logger is the set of logging methods supported by all contrail logger types.
// See glog's documentation for more detail.
type Logger interface {
	// HeaderTag returns the logging tag in use by the Logger. Useful to integrate
	// with other logging systems.
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
