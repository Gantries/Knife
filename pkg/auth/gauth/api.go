package gauth

type Whitelist interface {
	In(id string) bool
}
