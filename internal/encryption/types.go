package encryption

import "net/http"

type Encryptor struct {
	keyRead any
}

func NewEncryptor(key any) Encryptor {
	return Encryptor{keyRead: key}
}

type RoundTripperWithEncryption struct {
	next http.RoundTripper
	enc Encryptor
}

func NewRoundTripperWithEncryption(enc Encryptor) *RoundTripperWithEncryption {
	return &RoundTripperWithEncryption{next: http.DefaultTransport, enc: enc}
}