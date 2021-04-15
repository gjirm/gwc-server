package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"time"
)

// Create cookie hmac
func cookieSignature(domain, email, secret, expires string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(domain))
	hash.Write([]byte(email))
	hash.Write([]byte(expires))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

// makeCookie creates an auth cookie
func makeCookie(domain string, email string, secret string, expireTime time.Duration) string {
	expires := time.Now().Local().Add(expireTime)
	mac := cookieSignature(domain, email, secret, fmt.Sprintf("%d", expires.Unix()))
	value := fmt.Sprintf("%s|%d|%s", mac, expires.Unix(), email)
	return value
}

// Generate hashed cookie
// Cookie = hash(secret, cookie domain, email, expires)|expires|email
func main() {
	// Generate HMAC cookie
	domain := flag.String("domain", "example.com", "domain")

	expire := flag.Int("expire", 600, "expiration in seconds")

	email := flag.String("email", "john.doe@example.com", "email")

	secret := flag.String("secret", "P@ssword", "cookie secret for calculating HMAC")

	flag.Parse()

	expireTime := time.Second * time.Duration(*expire)

	cookie := makeCookie(*domain, *email, *secret, expireTime)

	fmt.Println("Cookie:", cookie)
}
