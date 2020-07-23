package util

import (
	"fmt"
	"testing"
)

func TestIp(t *testing.T) {
	ip, err := getIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ip.String())
}
