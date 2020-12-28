package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Response struct {
	Message string `json:"message"`
	Status int32 `json:"status"`
}

type missingParamError struct {
	s string
}

func (e *missingParamError) Error() string{
	return e.s
}

func ErrorResponse(w http.ResponseWriter, err error){
	http.Error(w, err.Error(), http.StatusBadRequest)
	return
}

func ForbiddenResponse(w http.ResponseWriter){
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(&Response{Message: "Unauthorized", Status: http.StatusForbidden}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusForbidden)
	w.Write(buf.Bytes())
}

func ObjectAddedToDatabase (w http.ResponseWriter, m string){
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(&Response{Message: m, Status: http.StatusCreated}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}

func JSONResponse(w http.ResponseWriter, status int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func ParseRequestBody (r *http.Request, body *map[string]string, needed []string) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err = json.Unmarshal(reqBody, body);err != nil {
		return err
	}
	for key, val := range *body {
		if val == ""{
			return &missingParamError{key+" Missing"}
		}
	}
	for i := range needed {
		if (*body)[needed[i]] == ""{
			return &missingParamError{needed[i]+" Missing"}
		}
	}
	return nil
}

func Encrypt(text string) (string, error) {
	plaintext := []byte(text)

	block, err := aes.NewCipher([]byte(os.Getenv("AES_KEY")))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(cryptoText string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher([]byte(os.Getenv("AES_KEY")))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}