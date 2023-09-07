package ip

import (
	"io"
	"net/http"
	"os"
)

var bodyAddr = "https://api.ipify.org"

func init() {
	v := os.Getenv("IP_WEB_ADDR")
	if v != "" {
		bodyAddr = v
	}
	Detectors["BODY"] = getIpFromWeb
}

func getIpFromWeb() (ip string, err error) {
	if bodyAddr == "" {
		err = ErrNotAvailable
		return
	}
	var resp *http.Response
	resp, err = http.Get(bodyAddr)
	if err == nil {
		var b []byte
		b, err = io.ReadAll(resp.Body)
		ip = string(b)
	}
	return
}
