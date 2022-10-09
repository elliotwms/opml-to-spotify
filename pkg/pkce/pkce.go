// Package pkce generates PKCE verifiers and challenges according to https://www.rfc-editor.org/rfc/rfc7636
package pkce

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/oauth2"
)

var encoding = base64.RawURLEncoding

type Verifier []byte

const LenMax = 128
const LenMin = 43

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~"

func NewVerifier(l int) Verifier {
	if l < LenMin || l > LenMax {
		panic(fmt.Sprintf("invalid verifier length: %d", l))
	}

	rand.Seed(time.Now().UnixNano())

	b := make([]byte, l)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return b
}

// Params is a convenience method which generates AuthCodeOptions compatible with the oauth2 package
func (v Verifier) Params() []oauth2.AuthCodeOption {
	return []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_verifier", string(v)),
	}
}

func (v Verifier) Challenge() Challenge {
	n := sha256.New()
	n.Write(v)

	return Challenge(encoding.EncodeToString(n.Sum(nil)))
}

type Challenge string

const challengeMethodS256 = "S256"

// Params is a convenience method which generates AuthCodeOptions compatible with the oauth2 package
func (c Challenge) Params() []oauth2.AuthCodeOption {
	return []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_challenge_method", challengeMethodS256),
		oauth2.SetAuthURLParam("code_challenge", string(c)),
	}
}
