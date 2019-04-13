package sendshit

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/keybase/go-triplesec"
	"github.com/vincent-petithory/dataurl"
)

// Upload structure
type Upload struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

// Response of API Upload
type Response struct {
	ID string `json:"id"`
}

// Generate a random hex string
func generateRandomString(size int) (string, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)

	return hex.EncodeToString(bytes), err
}

// Encrypt the file
func encryptFile(name string, data []byte, key string) (string, error) {
	var encodedStr string
	cipher, err := triplesec.NewCipher([]byte(key), nil, triplesec.LatestVersion)

	if err != nil {
		return encodedStr, err
	}

	dataURL := dataurl.EncodeBytes(data)

	upload := &Upload{
		URL:  dataURL,
		Name: name}

	json, _ := json.Marshal(upload)

	var encrypted []byte
	encrypted, err = cipher.Encrypt(json)

	if err != nil {
		return encodedStr, err
	}

	encodedStr = hex.EncodeToString(encrypted)

	return encodedStr, nil
}

// Uploaded the encrypted data to the API
func uploadFile(encryptedData string) (Response, error) {
	var response Response
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("upload", "encrypted")
	if err != nil {
		return response, err
	}

	fileWriter.Write([]byte(encryptedData))

	contentType := bodyWriter.FormDataContentType()

	bodyWriter.Close()

	resp, err := http.Post("https://api.sendsh.it/upload", contentType, bodyBuf)

	if err != nil {
		return response, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return response, err
	}

	json.Unmarshal(body, &response)

	return response, nil
}
