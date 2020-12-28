package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Response struct {
	Message string `json:"message"`
	Status int32 `json:"status"`
}

type missingParamError struct {
	s string
}

type UrlResponse struct {
	Url string `json:"url"`
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

func GenerateOAuthHeader (method, path, consumerKey, consumerSecret, token, tokenSecret string, extra map[string]string) (string, string) {
	vals := url.Values{}
	vals.Add("oauth_nonce", uuid.New().String())
	vals.Add("oauth_consumer_key", consumerKey)
	vals.Add("oauth_signature_method", "HMAC-SHA1")
	vals.Add("oauth_timestamp", strconv.Itoa(int(time.Now().Unix())))
	vals.Add("oauth_version", "1.0")
	vals.Add("oauth_token", token)

	for key, val := range extra{
		vals.Add(key, val)
	}

	if len(strings.Split(path, "?")) > 1 {
		params := strings.Split(strings.Split(path, "?")[1], "&")
		for i := range params {
			pair := strings.Split(params[i], "=")
			key, _ := url.QueryUnescape(pair[0])
			val, _ := url.QueryUnescape(pair[1])
			vals.Add(key, val)
		}
	}

	parameterString := strings.Replace(vals.Encode(), "+", "%20", -1)

	signatureBase := strings.ToUpper(method) + "&" + url.QueryEscape(strings.Split(path, "?")[0]) + "&" + url.QueryEscape(parameterString)
	signingKey := url.QueryEscape(consumerSecret)+"&"+url.QueryEscape(tokenSecret)
	signature := calculateSignature(signatureBase, signingKey)

	return "OAuth oauth_consumer_key=\"" + url.QueryEscape(vals.Get("oauth_consumer_key")) + "\", oauth_nonce=\"" + url.QueryEscape(vals.Get("oauth_nonce")) +
		"\", oauth_signature=\"" + url.QueryEscape(signature) + "\", oauth_signature_method=\"" + url.QueryEscape(vals.Get("oauth_signature_method")) +
		"\", oauth_timestamp=\"" + url.QueryEscape(vals.Get("oauth_timestamp")) + "\", oauth_token=\"" + url.QueryEscape(vals.Get("oauth_token")) +
		"\", oauth_version=\"" + url.QueryEscape(vals.Get("oauth_version")) + "\"", vals.Get("oauth_nonce")
}

func GenerateRequestOAuthHeader (method, path, callback, consumerKey, consumerSecret string) (string, string) {
	vals := url.Values{}
	vals.Add("oauth_nonce", uuid.New().String())
	vals.Add("oauth_consumer_key", consumerKey)
	vals.Add("oauth_signature_method", "HMAC-SHA1")
	vals.Add("oauth_timestamp", strconv.Itoa(int(time.Now().Unix())))
	vals.Add("oauth_version", "1.0")
	vals.Add("oauth_callback",callback)


	params := strings.Split(strings.Split(path, "?")[1], "&")
	for i := range params {
		pair := strings.Split(params[i],"=")
		key, _ := url.QueryUnescape(pair[0])
		val, _ := url.QueryUnescape(pair[1])
		vals.Add(key, val)
	}

	parameterString := strings.Replace(vals.Encode(), "+", "%20", -1)

	signatureBase := strings.ToUpper(method) + "&" + url.QueryEscape(strings.Split(path, "?")[0]) + "&" + url.QueryEscape(parameterString)
	signingKey := url.QueryEscape(consumerSecret)+"&"
	signature := calculateSignature(signatureBase, signingKey)

	return "OAuth oauth_callback=\"" + url.QueryEscape(vals.Get("oauth_callback")) +
		"\", oauth_consumer_key=\"" + url.QueryEscape(vals.Get("oauth_consumer_key")) +
		"\", oauth_nonce=\"" + url.QueryEscape(vals.Get("oauth_nonce")) +
		"\", oauth_signature=\"" + url.QueryEscape(signature) +
		"\", oauth_signature_method=\"" + url.QueryEscape(vals.Get("oauth_signature_method")) +
		"\", oauth_timestamp=\"" + url.QueryEscape(vals.Get("oauth_timestamp")) +
		"\", oauth_version=\"" + url.QueryEscape(vals.Get("oauth_version")) + "\"", vals.Get("oauth_nonce")
}

func calculateSignature(base, key string) string {
	hash := hmac.New(sha1.New, []byte(key))
	hash.Write([]byte(base))
	signature := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}