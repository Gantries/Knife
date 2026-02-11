package synch

import (
	"context"
	"testing"
	"time"
)

type mockLock struct {
	locked   bool
	unlocked bool
}

func (m *mockLock) Lock(ctx context.Context, source, owner string, timeout time.Duration) (bool, error) {
	m.locked = true
	return true, nil
}

func (m *mockLock) Unlock(ctx context.Context, source, owner string) (bool, error) {
	m.unlocked = true
	return true, nil
}

func TestLock(t *testing.T) {
	ctx := context.Background()
	m := &mockLock{}

	locked, err := m.Lock(ctx, "test", "owner", time.Second)
	if err != nil {
		t.Errorf("Lock() unexpected error: %v", err)
	}
	if !locked {
		t.Error("Lock() returned false, want true")
	}
	if !m.locked {
		t.Error("Lock() was not called")
	}
}

func TestUnlock(t *testing.T) {
	ctx := context.Background()
	m := &mockLock{}

	unlocked, err := m.Unlock(ctx, "test", "owner")
	if err != nil {
		t.Errorf("Unlock() unexpected error: %v", err)
	}
	if !unlocked {
		t.Error("Unlock() returned false, want true")
	}
	if !m.unlocked {
		t.Error("Unlock() was not called")
	}
}
