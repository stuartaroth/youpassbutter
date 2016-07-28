package main

import (
	"net/http"
	"net/url"
	"testing"
)

func getRequestWithQParam() *http.Request {
	u := url.URL{}
	u.RawQuery = "q=query&p=parameter1&p=parameter2"
	r := http.Request{}
	r.URL = &u
	return &r
}

func getRequestWithoutQParam() *http.Request {
	u := url.URL{}
	u.RawQuery = "p=1&p=2"
	r := http.Request{}
	r.URL = &u
	return &r
}

func TestAssignConfigDefaultValues_withoutValues(t *testing.T) {
	blankConfig := Config{"", 0, "", "", "", 0}
	c := AssignConfigDefaultValues(blankConfig)

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
	blankConfig := Config{"127.0.0.1", 3434, "", "", "", 8888}
	c := AssignConfigDefaultValues(blankConfig)

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
	_, _, err := GetQueryAndParams(r)

	if err == nil {
		t.Log("err should be != nil when no q parameter is passed")
		t.Fail()
	}
}

func TestQueryAndParams_withQParam(t *testing.T) {
	r := getRequestWithQParam()
	query, params, err := GetQueryAndParams(r)

	if err != nil {
		t.Log("err should be == nil when q parameter is passed")
		t.Fail()
	}

	if query != "query" {
		t.Log("query should == query when it is passed as the q parameter")
		t.Fail()
	}

	if params[0] != "parameter1" || params[1] != "parameter2" {
		t.Log("params should have parameter1 and parameter2 as items in slice")
		t.Fail()
	}
}

func TestIsSelectQuery_isSelect(t *testing.T) {
	if !IsSelectQuery("SELECT * FROM book") {
		t.Log("SELECT at the start of a string should return true")
		t.Fail()
	}
}

func TestIsSelectQuery_isNotSelect(t *testing.T) {
	if IsSelectQuery("DELETE FROM book") {
		t.Log("not SELECT at the start of a string should return false")
		t.Fail()
	}
}

func TestGetTypedInterface_int(t *testing.T) {
	i := GetTypedInterface("13")
	if i != 13 {
		t.Log("string '13' should become int 13")
		t.Fail()
	}
}

func TestGetTypedInterface_float(t *testing.T) {
	f := GetTypedInterface("7.89")
	if f != 7.89 {
		t.Log("string '7.89' should become float 7.89")
		t.Fail()
	}
}

func TestGetTypedInterface_bool(t *testing.T) {
	b := GetTypedInterface("t")
	if b != true {
		t.Log("string 't' should become bool true")
		t.Fail()
	}
}

func TestGetTypedInterface_string(t *testing.T) {
	s := GetTypedInterface("butter")
	if s != "butter" {
		t.Log("string 'butter' should become string 'butter'")
		t.Fail()
	}
}
