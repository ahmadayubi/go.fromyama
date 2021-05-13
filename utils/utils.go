package utils

import (
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

type missingParamError struct {
	s string
}

func (e *missingParamError) Error() string{
	return e.s
}

func ParseRequestBody (r *http.Request, body *map[string]string, needed []string) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
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

func AESEncrypt(text string) (string, error) {
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

func AESDecrypt(cryptoText string) (string, error) {
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
	values := url.Values{}
	values.Add("oauth_nonce", uuid.New().String())
	values.Add("oauth_consumer_key", consumerKey)
	values.Add("oauth_signature_method", "HMAC-SHA1")
	values.Add("oauth_timestamp", strconv.Itoa(int(time.Now().Unix())))
	values.Add("oauth_version", "1.0")
	values.Add("oauth_token", token)

	for key, val := range extra{
		values.Add(key, val)
	}

	if len(strings.Split(path, "?")) > 1 {
		params := strings.Split(strings.Split(path, "?")[1], "&")
		for i := range params {
			pair := strings.Split(params[i], "=")
			key, _ := url.QueryUnescape(pair[0])
			val, _ := url.QueryUnescape(pair[1])
			values.Add(key, val)
		}
	}

	parameterString := strings.Replace(values.Encode(), "+", "%20", -1)

	signatureBase := strings.ToUpper(method) + "&" + url.QueryEscape(strings.Split(path, "?")[0]) + "&" + url.QueryEscape(parameterString)
	signingKey := url.QueryEscape(consumerSecret)+"&"+url.QueryEscape(tokenSecret)
	signature := calculateSignature(signatureBase, signingKey)

	return "OAuth oauth_consumer_key=\"" + url.QueryEscape(values.Get("oauth_consumer_key")) +
		"\", oauth_nonce=\"" + url.QueryEscape(values.Get("oauth_nonce")) +
		"\", oauth_signature=\"" + url.QueryEscape(signature) +
		"\", oauth_signature_method=\"" + url.QueryEscape(values.Get("oauth_signature_method")) +
		"\", oauth_timestamp=\"" + url.QueryEscape(values.Get("oauth_timestamp")) +
		"\", oauth_token=\"" + url.QueryEscape(values.Get("oauth_token")) +
		"\", oauth_version=\"" + url.QueryEscape(values.Get("oauth_version")) +
		"\"", values.Get("oauth_nonce")
}

func GenerateRequestOAuthHeader (method, path, callback, consumerKey, consumerSecret string) (string, string) {
	values := url.Values{}
	values.Add("oauth_nonce", uuid.New().String())
	values.Add("oauth_consumer_key", consumerKey)
	values.Add("oauth_signature_method", "HMAC-SHA1")
	values.Add("oauth_timestamp", strconv.Itoa(int(time.Now().Unix())))
	values.Add("oauth_version", "1.0")
	values.Add("oauth_callback",callback)


	params := strings.Split(strings.Split(path, "?")[1], "&")
	for i := range params {
		pair := strings.Split(params[i],"=")
		key, _ := url.QueryUnescape(pair[0])
		val, _ := url.QueryUnescape(pair[1])
		values.Add(key, val)
	}

	parameterString := strings.Replace(values.Encode(), "+", "%20", -1)

	signatureBase := strings.ToUpper(method) + "&" + url.QueryEscape(strings.Split(path, "?")[0]) + "&" + url.QueryEscape(parameterString)
	signingKey := url.QueryEscape(consumerSecret)+"&"
	signature := calculateSignature(signatureBase, signingKey)

	return "OAuth oauth_callback=\"" + url.QueryEscape(values.Get("oauth_callback")) +
		"\", oauth_consumer_key=\"" + url.QueryEscape(values.Get("oauth_consumer_key")) +
		"\", oauth_nonce=\"" + url.QueryEscape(values.Get("oauth_nonce")) +
		"\", oauth_signature=\"" + url.QueryEscape(signature) +
		"\", oauth_signature_method=\"" + url.QueryEscape(values.Get("oauth_signature_method")) +
		"\", oauth_timestamp=\"" + url.QueryEscape(values.Get("oauth_timestamp")) +
		"\", oauth_version=\"" + url.QueryEscape(values.Get("oauth_version")) + "\"", values.Get("oauth_nonce")
}

func calculateSignature(base, key string) string {
	hash := hmac.New(sha1.New, []byte(key))
	hash.Write([]byte(base))
	signature := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}