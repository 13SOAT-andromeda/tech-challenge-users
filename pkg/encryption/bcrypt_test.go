package encryption

import "testing"

func TestHash_GeneratesNonEmptyHash(t *testing.T) {
	hash, err := Hash("mypassword")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == "" {
		t.Error("expected non-empty hash")
	}
}

func TestCompare_CorrectPassword(t *testing.T) {
	hash, _ := Hash("mypassword")
	if !Compare(hash, "mypassword") {
		t.Error("expected Compare to return true for correct password")
	}
}

func TestCompare_WrongPassword(t *testing.T) {
	hash, _ := Hash("mypassword")
	if Compare(hash, "wrongpassword") {
		t.Error("expected Compare to return false for wrong password")
	}
}

func TestHash_DifferentHashesSameInput(t *testing.T) {
	h1, _ := Hash("samepassword")
	h2, _ := Hash("samepassword")
	if h1 == h2 {
		t.Error("expected different hashes for same input due to bcrypt salting")
	}
}
