package utils

import "testing"

func TestEncrypt(t *testing.T) {
	got, _ := Encrypt("TEST STRING")
	want := ""

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestDecrypt(t *testing.T) {
	got, _ := Decrypt("TEST STRING")
	want := ""

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
