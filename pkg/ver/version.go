// Package ver provides version information for the application.
//
// The Version variable can be overridden at build time using ldflags:
// go build -ldflags "-X github.com/gantries/knife/pkg/ver.Version=0.1.0"
package ver

// Version : can be changed with `go build -ldflags "-X github.com/gantries/knife/pkg/ver.Version=0.1.0"`.
var Version string = "0.0.0"
