package base64

import "encoding/base64"

const (
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

var coder = base64.NewEncoding(base64Table)

func Encode(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))

}

func Decode(src string) (s string, err error) {
	buf, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return
	}
	s = string(buf)
	return
}
