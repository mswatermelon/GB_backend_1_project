package url_shorterer

import (
	"crypto/md5"
	"encoding/hex"
)

func ShortifyUrl(url string) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:10])
}
