package validate

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/sirupsen/logrus"

	"jirm.cz/gwc-server/internal/config"
)

// Create cookie hmac
func cookieSignature(domain, email, expires string, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(domain))
	hash.Write([]byte(email))
	hash.Write([]byte(expires))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

// ValidateCookie
func ValidateCookie(log *logrus.Logger, config config.Configs, hash string, expires string, userMail string) bool {

	//mac, err := base64.URLEncoding.DecodeString("IMU7VWeJ3-KGd7h5yBbkBkPYrtf0qEG9UaAkrHL2TA4=")
	mac, err := base64.URLEncoding.DecodeString(hash)
	if err != nil {
		log.Error("Unable to decode cookie mac")
		return false
	}

	expectedSignature := cookieSignature(config.Cookie.Domain, userMail, expires, config.Cookie.Secret)

	//expectedSignature := cookieSignature("safetica.cloud", "jiri.matejicek@safetica.com", "1587430262")
	expected, err := base64.URLEncoding.DecodeString(expectedSignature)
	if err != nil {
		log.Error("Unable to generate mac")
		return false
	}

	// Valid token?
	if !hmac.Equal(mac, expected) {
		log.Error("Invalid cookie mac")
		return false
	} else {
		log.Info("Cookie is valid!")
		return true
	}

}
