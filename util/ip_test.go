package util

import (
	"fmt"
	"github.com/iok200/go-ok/log"
	"testing"
)

func TestIp(t *testing.T) {
	ip, err := getIP()
	if err != nil {
		log.Infoln(err)
	}
	log.Infoln(ip.String())
}
