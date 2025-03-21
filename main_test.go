package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestHandleHello(t *testing.T) {
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	handleHello(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected Status OK; got %v", res.Status)
	}

	expected := "Hello, Go!"
	actual := rec.Body.String()
	if actual != expected {
		t.Errorf("Expected body %q; got %q", expected, actual)
	}
}

func TestGetPeople(t *testing.T) {
	req, err := http.NewRequest("GET", "/people", nil)
	if err != nil {
		t.Fatalf("Could not creare request: %v", err)
	}

	rec := httptest.NewRecorder()
	getPeople(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected Status OK; got %v", res.Status)
	}

	expected := []Person{}
	var actual []Person
	if err := json.NewDecoder(rec.Body).Decode(&actual); err != nil {
		t.Fatalf("Could not decode response body: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected body %v; got %v", expected, actual)
	}
}

func TestCreatePerson(t *testing.T) {
	person := Person{Name: "Tiji", Age: 30, IsLearningGo: true, Skills: map[string]int{"Go": 1, "Python": 2, "Java": 3}}
	var bodyReader *strings.Reader
	if err := json.NewDecoder(bodyReader).Decode(&person); err != nil {
		t.Fatalf("Could not decode response body: %v", err)
	}

	req, err := http.NewRequest("POST", "/create", bodyReader)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	createPerson(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected Status OK; got %v", res.Status)
	}
}
