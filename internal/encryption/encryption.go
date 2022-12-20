package encryption

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"reflect"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

func readKey(filepath string) any {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal().Err(err).Msg("reading file error")
	}

	block, _ := pem.Decode(bytes)
	if block == nil {
		log.Fatal().Err(err).Msg("decoding PEM error")
	}

	var key any

	switch block.Type {
	case "RSA_PRIVATE_KEY":
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			log.Fatal().Err(err).Msg("parsing rsa_private_key error")
		}

	case "PUBLIC_KEY":
		key, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			log.Fatal().Err(err).Msg("parsing public_key error")
		}

	default:
		log.Fatal().Err(err).Msg("file formatting error")
	}

	return key
}

func (e *Encryptor) EncryptMessage(b []byte) ([]byte, error) {
	hash := sha512.New()

	publicKey, ok := e.keyRead.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("type assertion error")
	}

	msgEnc, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, b, nil)
	if err != nil {
		return nil, errors.New("message encrypting error")
	}

	return msgEnc, nil
}

func (e *Encryptor) DecryptMessage(b []byte) ([]byte, error) {
	hash := sha512.New()

	privateKey, ok := e.keyRead.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("type assertion error")
	}

	msgDecr, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, b, nil)
	if err != nil {
		return nil, errors.New("message decrypting error")
	}

	return msgDecr, nil
}

func (rtwe RoundTripperWithEncryption) RoundTrip(req *http.Request) (*http.Response, error) {
	if (reflect.DeepEqual(rtwe, RoundTripperWithEncryption{})) {
		return rtwe.next.RoundTrip(req)
	}

	var b bytes.Buffer

	switch req.Method {
	case http.MethodPost:
		_, err := b.ReadFrom(req.Body)
		if err != nil {
			return nil, errors.New("reading request error")
		}

		buf, err := rtwe.enc.EncryptMessage(b.Bytes())
		if err != nil {
			return nil, errors.New("encrypting message error")
		}

		req.Body = io.NopCloser(bytes.NewReader(buf))
		req.ContentLength = int64(len(buf))
		return rtwe.next.RoundTrip(req)

	default:
		rtwe.next.RoundTrip(req)
	}
	return rtwe.next.RoundTrip(req)
}