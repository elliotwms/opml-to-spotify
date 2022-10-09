package pkce

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

func TestNewVerifier(t *testing.T) {
	for _, i := range []int{
		LenMin, LenMax,
	} {
		t.Run(fmt.Sprintf("Len%d", i), func(t *testing.T) {
			v := NewVerifier(i)

			if len(v) != i {
				t.Fatalf("unexpected verifier length: %d", len(v))
			}

			m, err := regexp.Match(`[a-zA-Z0-9\-._~]{`+strconv.Itoa(i)+`}`, v)
			if err != nil {
				t.Fatal(err)
			}

			if !m {
				t.Fatal("did not match")
			}
		})
	}
}

func TestNewVerifierPanicsWithInvalidLength(t *testing.T) {
	for _, i := range []int{
		0,
		LenMin - 1,
		LenMax + 1,
		512,
	} {
		t.Run(fmt.Sprintf("InvalidLen%d", i), func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("did not panic")
				}
			}()

			NewVerifier(i)
		})
	}
}

func TestVerifier_Challenge(t *testing.T) {
	// values from https://www.oauth.com/playground/authorization-code-with-pkce.html
	v := Verifier("6mk7z_YTkAjaxUbx06ro9rEx30EeKZtMm3s1u-EZL_KULH4Y")

	expected := "6PUMCbtSovfEx9UxUIyi3qpbNImUYZ4cqIiqh0ggbqI"
	actual := v.Challenge()

	if string(actual) != expected {
		t.Fatalf("invalid challenge. Expected %s, got %s", expected, actual)
	}
}
