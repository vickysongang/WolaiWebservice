// encrypt_password
package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"io"

	"github.com/satori/go.uuid"
)

func GenerateSalt() string {
	return uuid.NewV4().String()
}

func EncryptPassword(password, salt string) string {
	m := md5.New()
	io.WriteString(m, password)
	pwd := hex.EncodeToString(m.Sum(nil))
	io.WriteString(m, salt+pwd+salt)
	pwd = hex.EncodeToString(m.Sum(nil))
	return pwd
}
