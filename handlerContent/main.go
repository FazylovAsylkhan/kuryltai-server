package handlerContent

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func UploadFile(userId uuid.UUID, base64Str string) (link string, err error) {
	contentType, ext, err := getMimeType(base64Str)
	filename := fmt.Sprintf("%s.%s", getRandomFileName(), ext)
	path := fmt.Sprintf("/assets/%s", contentType)
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("uploadFile os.MkdirAll error: %s", err)
	}
	fileOnServer, err := os.Create(fmt.Sprintf(".%s/%s", path, filename))
	if err != nil {
		return "", fmt.Errorf("uploadFile os.Create error: %s", err)
	}
	defer fileOnServer.Close()

	parts := strings.Split(base64Str, ",")
	decodedData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("base64.DecodeString error: %s", err)
	}

	_, err = fileOnServer.Write(decodedData)
	if err != nil {
		return "", fmt.Errorf("uploadFile os.Create error: %s", err)
	}
	
	domenString := os.Getenv("DOMEN")
	portString := os.Getenv("PORT")
	return domenString + ":" + portString + path + "/" + filename, nil
}

func getMimeType(base64 string) (contentType string, ext string, err error) {
	if !strings.HasPrefix(base64, "data:") {
		return "", "", errors.New("invalid base64 data string")
	}
	parts := strings.SplitN(base64, ";base64,", 2)
	if len(parts) < 2 {
		return "", "", errors.New("invalid base64 string format")
	}
	mimeType := strings.TrimPrefix(parts[0], "data:")
	arrMimeType := strings.Split(mimeType, "/")

	return arrMimeType[0], arrMimeType[1], nil
}

func getRandomFileName() (string) {
	t := time.Now().Format("2006-01-02")
	randUuid := uuid.Must(uuid.NewUUID()).String()
	return strings.ToLower(t + "_" + randUuid)
}