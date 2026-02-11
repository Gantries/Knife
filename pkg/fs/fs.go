// Package fs contains file iostream related utilities
package fs

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path"
)

// With open file named as name with permission perm and flags flag, then consume the file in fn, finally close this file
func With(name string, flag int, perm os.FileMode, fn func(*os.File) error) error {
	// #nosec G304 - file path comes from function parameter
	c, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return err
	}
	defer func() { _ = (*c).Close() }()
	return fn(c)
}

// WithReader reads file with name and consume it in fn and close it finally.
func WithReader(name string, fn func(*os.File) error) error {
	return With(name, os.O_RDONLY, 0, fn)
}

// WithScanner reads file into line and invoke fn one by one
func WithScanner[T any](name string, maximum int, stoppable func(line string) (bool, T)) (found bool, value T, err error) {
	var r io.Reader
	// #nosec G304 - file path comes from function parameter
	if r, err = os.Open(name); err != nil {
		return
	}
	scanner := bufio.NewScanner(r)
	defer func() { err = scanner.Err() }()
	scanner.Buffer(make([]byte, maximum), maximum)
	for scanner.Scan() {
		if found, value = stoppable(scanner.Text()); found {
			return found, value, nil
		}
	}
	return
}

// WithWriter open file with name for writing and consume the opened file in fn
func WithWriter(name string, fn func(*os.File) error) error {
	return With(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm, fn)
}

// WithCloser consume resource r in action and close it finally.
func WithCloser[T any](r T, close func(T) error, action func(T) error) error {
	defer func() { _ = close(r) }()
	return action(r)
}

// WithReadCloser read content in resource rc in action and close it finally.
func WithReadCloser(rc io.ReadCloser, action func([]byte) error) (err error) {
	defer func() { err = errors.Join(rc.Close(), err) }()
	b, err := io.ReadAll(rc)
	if err != nil {
		return err
	}
	return action(b)
}

// WithWriteCloser write anything in action to rc and close rc finally.
func WithWriteCloser(wc io.WriteCloser, action func(wc io.WriteCloser) error) (err error) {
	defer func() { err = errors.Join(wc.Close(), err) }()
	return action(wc)
}

// WithOpen open resource with open function, consume it in action and close it finally.
func WithOpen[T any](open func() (T, error), close func(T) error, action func(T) error) error {
	inst, err := open()
	if err != nil {
		return err
	}
	defer func() { _ = close(inst) }()
	return action(inst)
}

// Save read content from r and write it into file located in folder.
func Save[T io.Reader](r T, folder, file string) error {
	return With(path.Join(folder, file), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm, func(c *os.File) error {
		_, err := io.Copy(c, r)
		return err
	})
}
