package sms

import (
   "bytes"
   "encoding/gob"
   "fmt"
   "regexp"
   "unicode/utf8"
)

const (
   MaxRuneSize = 4
   MaxRuneCount = 2048
   MaxBodyLen = MaxRuneSize * MaxRuneCount // == maximum 8kib per message
)

// matcher of russian mobile phone number
var matchPhone = regexp.MustCompile(`79\d{2}\d{7}`)

type SMS struct {
   Phone string
   Body  string
}

func (sms *SMS) PhoneValid() bool {
   if len(sms.Phone) == 0 {
      return false
   }
   return matchPhone.Match([]byte(sms.Phone))
}

func (sms *SMS) BodyValid() bool {
   size := len(sms.Body)
   if size == 0 {
      return false
   }
   if len(sms.Body) >  MaxBodyLen {
      return false
   }

   if utf8.RuneCountInString(sms.Body) > MaxRuneCount {
      return false
   }

   return true
}

func Encode(msg SMS) ([]byte, error) {
   var buf bytes.Buffer

   enc := gob.NewEncoder(&buf)
   err := enc.Encode(&msg)
   if err != nil {
      return nil, fmt.Errorf("sms encode: %w", err)
   }

   return buf.Bytes(), nil
}

func Decode(data []byte) (SMS, error) {
   var msg SMS

   dec := gob.NewDecoder(bytes.NewBuffer(data))
   err := dec.Decode(&msg)
   if err != nil {
      return msg, fmt.Errorf("sms decode: %w", err)
   }

   return msg, nil
}
