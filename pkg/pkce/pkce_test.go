package pkce

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestNewVerifier(t *testing.T) {
	for _, i := range []int{
		LenMin, LenMax,
	} {
		t.Run(fmt.Sprintf("Len%d", i), func(t *testing.T) {
			v := NewVerifier(i)

			require.Len(t, v, i)
			require.Regexp(t, `[a-zA-Z0-9\-._~]{`+strconv.Itoa(i)+`}`, string(v))
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
			require.Panics(t, func() {
				NewVerifier(i)
			})
		})
	}
}

func TestVerifier_Challenge(t *testing.T) {
	// values from https://www.oauth.com/playground/authorization-code-with-pkce.html
	v := Verifier("6mk7z_YTkAjaxUbx06ro9rEx30EeKZtMm3s1u-EZL_KULH4Y")

	require.Equal(t, "6PUMCbtSovfEx9UxUIyi3qpbNImUYZ4cqIiqh0ggbqI", v.Challenge())
}
