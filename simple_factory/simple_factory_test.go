package simple_factory

import (
	"testing"
)

func TestNewApi(t *testing.T) {
	api := NewApi("huawei")
	say := api.say("HuaWei")
	if say != "HuaWei" {
		t.Fatal("test is fail")
	}
}
