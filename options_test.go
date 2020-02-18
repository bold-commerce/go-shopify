package goshopify

import (
	"fmt"
	"testing"
)

func TestWithVersion(t *testing.T) {
	c := NewClient(app, "fooshop", "abcd", WithVersion(testApiVersion))
	expected := fmt.Sprintf("admin/api/%s", testApiVersion)
	if c.pathPrefix != expected {
		t.Errorf("WithVersion client.pathPrefix = %s, expected %s", c.pathPrefix, expected)
	}
}

func TestWithVersionNoVersion(t *testing.T) {
	c := NewClient(app, "fooshop", "abcd", WithVersion(""))
	expected := "admin"
	if c.pathPrefix != expected {
		t.Errorf("WithVersion client.pathPrefix = %s, expected %s", c.pathPrefix, expected)
	}
}

func TestWithoutVersionInInitiation(t *testing.T) {
	c := NewClient(app, "fooshop", "abcd")
	expected := "admin"
	if c.pathPrefix != expected {
		t.Errorf("WithVersion client.pathPrefix = %s, expected %s", c.pathPrefix, expected)
	}
}

func TestWithVersionInvalidVersion(t *testing.T) {
	c := NewClient(app, "fooshop", "abcd", WithVersion("9999-99b"))
	expected := "admin"
	if c.pathPrefix != expected {
		t.Errorf("WithVersion client.pathPrefix = %s, expected %s", c.pathPrefix, expected)
	}
}

func TestWithRetry(t *testing.T) {
	c := NewClient(app, "fooshop", "abcd", WithRetry(5))
	expected := 5
	if c.retries != expected {
		t.Errorf("WithRetry client.retries = %d, expected %d", c.retries, expected)
	}
}
