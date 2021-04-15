package validate

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"jirm.cz/gwc-server/internal/config"
)

// Create cookie hmac
func cookieSignature(config config.Configs, email, expires string) string {
	hash := hmac.New(sha256.New, []byte(config.Cookie.Secret))
	hash.Write([]byte(config.Cookie.Domain))
	hash.Write([]byte(email))
	hash.Write([]byte(expires))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

// ValidateCookie verifies that a cookie matches the expected format of:
// Cookie = hash(secret, cookie domain, email, expires)|expires|email
func ValidateCookie(config config.Configs, cookie string) (bool, string) {

	// Check cookie format
	parts := strings.Split(cookie, "|")
	if len(parts) != 3 {
		return false, "Wrong cookie format"
	}

	mac, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return false, "Unable to decode cookie mac"
	}

	expectedSignature := cookieSignature(config, parts[2], parts[1])
	expected, err := base64.URLEncoding.DecodeString(expectedSignature)
	if err != nil {
		return false, "Failed request by " + parts[2] + " - unable to generate mac"
	}

	// Valid token?
	if !hmac.Equal(mac, expected) {
		return false, "Failed request by " + parts[2] + " - invalid cookie mac"
	}

	expires, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return false, "Failed request by " + parts[2] + " - unable to parse cookie expiry"
	}

	// Token expired ?
	if time.Unix(expires, 0).Before(time.Now()) {
		return false, "Failed request by " + parts[2] + " - cookie has expired"
	}

	// Token is valid - return true, email
	return true, parts[2]
}
