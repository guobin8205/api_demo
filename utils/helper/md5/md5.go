package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func New(t string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(t))
	cipherStr := md5Ctx.Sum(nil)
	md5String := hex.EncodeToString(cipherStr)
	return md5String
}
