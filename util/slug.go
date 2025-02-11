package util

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

func GenerateSlugByName(name string) string {
	slugName := strings.ReplaceAll(slug.Make(name), "-", "")

	return slugName
}

func GetUsernameFrom(email string) string {
	arrEmail := strings.Split(email, "@")
	return arrEmail[0]
}