package main

import (
	"net/http"
	"net/url"
	"testing"
)

func getRequestWithoutQParam() *http.Request {
	u := url.URL{}
	u.RawQuery = "p=1&p=2"
	r := http.Request{}
	r.URL = &u
	return &r
}

func TestAssignConfigDefaultValues_withoutValues(t *testing.T) {
	blankConfig := serverConfig{"", 0, "", "", "", 0}
	c := assignServerConfigDefaultValues(blankConfig)

	if c.DataHost != "localhost" {
		t.Log("DataHost was not set to value localhost")
		t.Fail()
	}

	if c.DataPort != 3306 {
		t.Log("DataPort was not set to value 3306")
		t.Fail()
	}

	if c.Port != 8080 {
		t.Log("Port was not set to value 8080")
		t.Fail()
	}
}

func TestAssignConfigDefaultValues_withValues(t *testing.T) {
	blankConfig := serverConfig{"127.0.0.1", 3434, "", "", "", 8888}
	c := assignServerConfigDefaultValues(blankConfig)

	if c.DataHost != "127.0.0.1" {
		t.Log("DataHost was changed from value 127.0.0.1")
		t.Fail()
	}

	if c.DataPort != 3434 {
		t.Log("DataPort was changed from value 3434")
		t.Fail()
	}

	if c.Port != 8888 {
		t.Log("Port was changed from value 8888")
		t.Fail()
	}
}

func TestQueryAndParams_withoutQParam(t *testing.T) {
	r := getRequestWithoutQParam()
	_, _, err := getQueryAndParams(r)

	if err == nil {
		t.Log("err should be != nil when no q parameter is passed")
		t.Fail()
	}
}

func TestIsSelectQuery_isSelect(t *testing.T) {
	if !isSelectQuery("SELECT * FROM book") {
		t.Log("SELECT at the start of a string should return true")
		t.Fail()
	}
}

func TestIsSelectQuery_isNotSelect(t *testing.T) {
	if isSelectQuery("DELETE FROM book") {
		t.Log("not SELECT at the start of a string should return false")
		t.Fail()
	}
}

func TestGetTypedInterface_int(t *testing.T) {
	i := getTypedInterface("13")
	if i != 13 {
		t.Log("string '13' should become int 13")
		t.Fail()
	}
}

func TestGetTypedInterface_float(t *testing.T) {
	f := getTypedInterface("7.89")
	if f != 7.89 {
		t.Log("string '7.89' should become float 7.89")
		t.Fail()
	}
}

func TestGetTypedInterface_bool(t *testing.T) {
	b := getTypedInterface("t")
	if b != true {
		t.Log("string 't' should become bool true")
		t.Fail()
	}
}

func TestGetTypedInterface_string(t *testing.T) {
	s := getTypedInterface("butter")
	if s != "butter" {
		t.Log("string 'butter' should become string 'butter'")
		t.Fail()
	}
}
