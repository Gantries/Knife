// Package kube provides utilities for reading Kubernetes service account secrets.
//
// It enables applications running in Kubernetes to access their namespace,
// service account token, and CA certificate from the standard mounted paths.
package kube

import "os"

const (
	FileNamespace = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	FileToken     = "/var/run/secrets/kubernetes.io/serviceaccount/token" // #nosec G101 - standard K8s service account path
	FileCa        = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

func readfile(file string) (body string, err error) {
	if _, err := os.Stat(file); err != nil {
		return "", err
	}
	b, err := os.ReadFile(file) // #nosec G304 - file path is from constants defined in this package
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Namespace() (body string, err error) {
	return readfile(FileNamespace)
}

func Token() (body string, err error) {
	return readfile(FileToken)
}

func Certificate() (body string, err error) {
	return readfile(FileCa)
}
