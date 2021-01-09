package apimsg

import (
	"testing"
)

func TestGetRes(t *testing.T) {
	res := GetRes("Success", nil)
	if res.Status != 200 {
		t.Error("Status Code not match")
	}
}
