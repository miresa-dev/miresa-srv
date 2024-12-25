package verifier

import (
	"crypto/rand"
	"math/big"
	"maps"

	"github.com/Kaamkiya/nanoid-go"
)

// TODO: hash captcha responses so bots can't use this map.
// OR: get the servers to make their own captchas. More i18n friendly.
var captchas = map[string]string{
	"Is ice hot or cold?": "cold",
	"If today is Sunday, what day is tomorrow?": "monday",
	"How many months in a year?": "twelve",
}

// pairs stores SIDs and CAPTCHAs as sid: captcha.
var pairs = map[string]string{}

// GenCaptchaSIDPair is called when a GET request is made to /init. It returns
// a random SID and CAPTCHA.
func GenCaptchaSIDPair() (string, string, error) {
	sid := nanoid.Nanoid(64, nanoid.DefaultCharset)

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(captchas))))
	if err != nil {
		return "", "", err
	}

	keys := make([]string, len(captchas))
	i := 0
	for key := range maps.Keys(captchas) {
		keys[i] = key
		i++
	}

	captchaQuestion := keys[n.Int64()]

	pairs[sid] = captchaQuestion

	return sid, captchaQuestion, nil
}

// This is called when a POST request is made to /u and when a user tries to
// log in. It prevents writing scripts that allow users to infinitely make bot
// accounts.
func VerifyPair(sid, captcha string) bool {
	for question, answer := range captchas {
		if answer == captcha && pairs[sid] == question {
			return true
		}
	}
	return false
}
