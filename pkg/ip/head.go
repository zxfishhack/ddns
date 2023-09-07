package ip

import (
	"net/http"
	"os"
)

var ipHead *http.Client
var (
	head     = "X-Client-IP"
	headAddr = ""
)

func init() {
	ipHead = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	Detectors["HEAD"] = getIpFromHead
	v := os.Getenv("IP_HEAD_FIELD")
	if v != "" {
		head = v
	}
	v = os.Getenv("IP_HEAD_URL")
	if v != "" {
		headAddr = v
	}
}

func getIpFromHead() (ip string, err error) {
	if headAddr == "" {
		err = ErrNotAvailable
		return
	}
	var resp *http.Response
	resp, err = ipHead.Head(headAddr)
	if err == nil {
		ip = resp.Header.Get(head)
	}
	return
}
