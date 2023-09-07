package ip

func GetMyIP() (ip string, err error) {
	for _, detector := range UseDetector {
		if f, ok := Detectors[detector]; ok {
			ip, err = f()
			if err == nil && ip != "" {
				break
			}
		}
	}
	return
}
