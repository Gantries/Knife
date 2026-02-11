package cache

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gantries/knife/pkg/lists"
	"github.com/stretchr/testify/assert"
)

var properties = Properties{
	Type:       Redis,
	Addresses:  *lists.Of[string]("127.0.0.1:6379"),
	Database:   0,
	Credential: Credential{},
	Pool:       Pool{},
}

func TestNew(t *testing.T) {
	ctxt := context.Background()
	c, m := New(ctxt, &properties)

	// Skip Redis operations if cache is not properly initialized
	if c == nil || !m.Fine() || c.Ping(ctxt) != nil {
		t.Skip("Redis not available, skipping cache operations test")
		return
	}

	assert.Nil(t, c.Push(ctxt, "knife:key", "value"))
	v, err := c.Pop(ctxt, "knife:key")
	assert.Nil(t, err)
	assert.True(t, v == "value")
}

func TestCache_Lock_Repeat(t *testing.T) {
	ctx := context.Background()
	c, err := NewRedis(ctx, &properties)
	if err != nil {
		t.Skip("Redis not available")
		return
	}
	flag1, _ := c.Lock(ctx, "keyOne", "name1", time.Second*30)
	if flag1 {
		println("get lock success")
	} else {
		println("get lock fail")
	}
	flag2, _ := c.Unlock(ctx, "keyOne", "name1")
	if flag2 {
		println("unlock success")
	} else {
		println("unlock fail")
	}
}

func TestCache_Lock(t *testing.T) {
	ctx := context.Background()
	c, err := NewRedis(ctx, &properties)
	if err != nil {
		t.Skip("Redis not available")
		return
	}
	wg := &sync.WaitGroup{}
	count := 10
	wg.Add(1)
	go func() {
		defer wg.Done()
		flag, _ := c.Lock(ctx, "keyOne", "name1", time.Second*3)
		if flag {
			println("1--get lock")
			defer func() { _, _ = c.Unlock(ctx, "keyOne", "name1") }()
			println("1--run time start")
			time.Sleep(time.Second) // do something
			println("1--run time end")
			count *= 2
		} else {
			println("1--not get lock")
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		flag, _ := c.Lock(ctx, "keyTwo", "name1", time.Second*3)
		if flag {
			println("2--get lock")
			defer func() { _, _ = c.Unlock(ctx, "keyTwo", "name1") }()
			println("2--run time start")
			time.Sleep(time.Second) // do something
			println("2--run time end")
			count *= 2
		} else {
			println("2--not get lock")
		}
	}()
	wg.Wait()
}

func TestRedisCache_Hash(t *testing.T) {
	ctxt := context.Background()
	c, m := New(ctxt, &properties)

	if c == nil || !m.Fine() || c.Ping(ctxt) != nil {
		t.Skip("Redis not available")
		return
	}

	key := "knife:hash"
	_, e := c.Del(ctxt, key)
	assert.Nil(t, e)
	n, e := c.HSet(ctxt, key, "k1", "v1", "k2", "v2")
	assert.Nil(t, e)
	assert.True(t, n == 2)
	h, e := c.HGetAll(ctxt, key)
	assert.Nil(t, e)
	assert.True(t, h["k1"] == "v1" && h["k2"] == "v2")
	v, e := c.HGet(ctxt, key, "k1")
	assert.Nil(t, e)
	assert.True(t, v == "v1")
	r, e := c.HDel(ctxt, key, "k1")
	assert.Nil(t, e)
	assert.True(t, r == 1)
}

func TestRedisCache_List(t *testing.T) {
	ctxt := context.Background()
	c, m := New(ctxt, &properties)

	if c == nil || !m.Fine() || c.Ping(ctxt) != nil {
		t.Skip("Redis not available")
		return
	}

	key := "knife:list"
	n, e := c.RPush(ctxt, key, "v1", "v2", "v3")
	assert.Nil(t, e)
	assert.True(t, n == 3)
	v, e := c.LPop(ctxt, key)
	assert.Nil(t, e)
	assert.True(t, v == "v1")
}

func TestRedisCache_String(t *testing.T) {
	ctxt := context.Background()
	c, m := New(ctxt, &properties)

	if c == nil || !m.Fine() || c.Ping(ctxt) != nil {
		t.Skip("Redis not available")
		return
	}

	o, e := c.Exists(ctxt, "k1")
	assert.Nil(t, e)
	assert.True(t, o == 0)
	_, e = c.Set(ctxt, "k1", "v1", 10)
	assert.Nil(t, e)
	v, e := c.Get(ctxt, "k1")
	assert.Nil(t, e)
	assert.True(t, v == "v1")
}
