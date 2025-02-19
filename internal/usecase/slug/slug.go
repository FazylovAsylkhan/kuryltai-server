package slug

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/gosimple/slug"
	"github.com/jaswdr/faker"
)

func randomString(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)[:n]
}

func GenerateRandomSlug() string {
	f := faker.New()
	randomFirstName := f.Person().FirstName()
	randomPart := randomString(4)

	return slug.Make(randomFirstName + randomPart)
}

func GetUsernameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func GenerateSlug() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}