package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dayterr/go_agent_metrics/internal/metric"
	"log"
)

func EncryptMetric(m metric.Metrics, key string) string {
	switch m.MType {
	case "gauge":
		log.Println("value", m.Value)
		src := fmt.Sprintf("%s:gauge:%f", m.ID, m.Value)
		h := hmac.New(sha256.New, []byte(key))
		h.Write([]byte(src))
		return hex.EncodeToString(h.Sum(nil))
	case "counter":
		log.Println("delta", m.Delta)
		src := fmt.Sprintf("%s:counter:%d", m.ID, m.Delta)
		h := hmac.New(sha256.New, []byte(key))
		h.Write([]byte(src))
		return hex.EncodeToString(h.Sum(nil))
	}
	return ""
}

func CheckHash(m metric.Metrics, hash string) bool {
	return hmac.Equal([]byte(m.Hash), []byte(hash))
}
