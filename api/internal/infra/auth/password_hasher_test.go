package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcryptHasher_Hash(t *testing.T) {
	h := NewBcryptHasher(bcrypt.MinCost)

	hashed, err := h.Hash("mysecret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte("mysecret")); err != nil {
		t.Fatal("hashed value is not a valid bcrypt hash of the input")
	}
}

func TestBcryptHasher_Compare_Match(t *testing.T) {
	h := NewBcryptHasher(bcrypt.MinCost)

	hashed, err := h.Hash("mysecret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	same, err := h.Compare("mysecret", hashed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !same {
		t.Fatal("expected Compare to return true for matching password")
	}
}

func TestBcryptHasher_Compare_Mismatch(t *testing.T) {
	h := NewBcryptHasher(bcrypt.MinCost)

	hashed, err := h.Hash("mysecret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	same, err := h.Compare("wrongpassword", hashed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if same {
		t.Fatal("expected Compare to return false for wrong password")
	}
}

func TestBcryptHasher_InvalidCost_FallsBackToDefault(t *testing.T) {
	h := NewBcryptHasher(-1)
	if h.cost != bcrypt.DefaultCost {
		t.Fatalf("expected cost %d, got %d", bcrypt.DefaultCost, h.cost)
	}
}
