package util

import (
	"Backend-Review/model"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

func GenerateSalt() string {
	character := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 255)
	for i := range b {
		rand.Seed(time.Now().UnixMicro())
		b[i] = character[rand.Intn(len(character))]
	}
	return string(b)
}

func Hashing(strs ...string) string {
	totalstr := ""
	for _, v := range strs {
		totalstr = totalstr + v
	}
	hash := sha256.Sum256([]byte(totalstr))
	return hex.EncodeToString(hash[:])
}

func GenerateToken(userid string, timestamp int64) string {

	character := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 255)
	for i := range b {
		b[i] = character[rand.Intn(len(character))]
	}
	t := &model.Token{
		UserID:    uuid.FromStringOrNil(userid),
		CreatedAt: timestamp,
	}
	hash := sha256.Sum256([]byte(userid + strconv.FormatInt(timestamp, 10) + string(b)))

	token := hash[:]

	t.Token = hex.EncodeToString(token)

	return t.Token
}
