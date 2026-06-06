package wasmplugin

import (
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

func TestUserInfoMsgpackRoundTrip(t *testing.T) {
	expected := UserInfo{ID: 42, FullName: "Иван Иванов"}
	data, err := msgpack.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	var actual UserInfo
	if err := msgpack.Unmarshal(data, &actual); err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}

func TestGetUserInfoRejectsInvalidID(t *testing.T) {
	_, err := GetUserInfo(0)
	if err == nil {
		t.Fatal("expected error for zero user ID")
	}
}
