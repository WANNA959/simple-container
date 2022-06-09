package utils

import (
	"fmt"
	"testing"
)

func TestGetChildPids(t *testing.T) {
	pids := GetChildPids(59767)
	fmt.Println(pids)
}
