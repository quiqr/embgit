package main
import (
  "encoding/base64"
)

func toBase64(b []byte) string {
  return base64.StdEncoding.EncodeToString(b)
}

