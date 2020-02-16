package crawler

import (
	"crypto/md5"
	"fmt"
)

func hashString(data interface{}) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprint(data))))
}
