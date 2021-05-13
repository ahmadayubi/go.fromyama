package utils

import "testing"

func TestEncrypt(t *testing.T) {
	got, _ := AESEncrypt("TEST STRING")
	want := ""

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestDecrypt(t *testing.T) {
	got, _ := AESDecrypt("TEST STRING")
	want := ""

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
