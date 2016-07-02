package mysql

import (
	"testing"
)

const (
	address  = ""
	username = ""
	password = ""
	timeout  = 30 * 1000000
)

func TestGetMysqlConnector(t *testing.T) {
	_, err := GetMysqlConnector(address, username, password, timeout)
	if err != nil {
		t.Error(err)
	}
	if err != nil {
		t.Error(err)
	}

}
