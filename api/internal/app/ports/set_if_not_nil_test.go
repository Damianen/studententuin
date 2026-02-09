package ports

import (
	"errors"
	"testing"
)

func TestSetIfNotNil_NilValue(t *testing.T) {
	updates := map[string]any{}

	err := SetIfNotNil(updates, "key", (*string)(nil), func(s string) (any, error) {
		return s, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updates) != 0 {
		t.Fatalf("expected empty map, got %v", updates)
	}
}

func TestSetIfNotNil_NonNilValue(t *testing.T) {
	updates := map[string]any{}
	val := "hello"

	err := SetIfNotNil(updates, "key", &val, func(s string) (any, error) {
		return s + "!", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updates["key"] != "hello!" {
		t.Fatalf("expected 'hello!', got %v", updates["key"])
	}
}

func TestSetIfNotNil_TransformError(t *testing.T) {
	updates := map[string]any{}
	val := ""

	err := SetIfNotNil(updates, "key", &val, func(s string) (any, error) {
		return nil, errors.New("validation failed")
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "validation failed" {
		t.Fatalf("expected 'validation failed', got %q", err.Error())
	}
	if len(updates) != 0 {
		t.Fatalf("expected empty map after error, got %v", updates)
	}
}
