package ip

import "errors"

type Detector func() (ip string, err error)

var Detectors = make(map[string]Detector)

var ErrNotAvailable = errors.New("detector is not available")

var UseDetector = []string{"HEAD", "BODY"}
