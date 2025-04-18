package cache_test

import (
	"context"
	"testing"

	"github.com/bilte-co/bilte/internal/cache"
)

func TestInMemoryCache(t *testing.T) {
	ctx := context.Background()
	cache, err := cache.NewCache(ctx)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Test Set and Get
	ok := cache.Set("foo", "bar")
	if !ok {
		t.Fatalf("Failed to set value: %v", err)
	}

	value, ok := cache.Get("foo")
	if !ok {
		t.Fatalf("Expected key 'foo' to exist")
	}
	if value != "bar" {
		t.Fatalf("Expected 'bar', got %v", value)
	}

	// Test Delete
	cache.Delete("foo")

	_, ok = cache.Get("foo")
	if ok {
		t.Fatalf("Expected key 'foo' to be deleted")
	}
}
