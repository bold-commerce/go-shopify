package goshopify

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestWithVersion(t *testing.T) {
	c := MustNewClient(app, "fooshop", "abcd", WithVersion(testApiVersion))
	expected := fmt.Sprintf("admin/api/%s", testApiVersion)
	if c.pathPrefix != expected {
		t.Errorf("WithVersion client.pathPrefix = %s, expected %s", c.pathPrefix, expected)
	}
}

func TestWithVersionNoVersion(t *testing.T) {
	c := MustNewClient(app, "fooshop", "abcd", WithVersion(""))
	expected := "admin"
	if c.pathPrefix != expected {
		t.Errorf("WithVersion client.pathPrefix = %s, expected %s", c.pathPrefix, expected)
	}
}

func TestWithoutVersionInInitiation(t *testing.T) {
	c := MustNewClient(app, "fooshop", "abcd")
	expected := "admin"
	if c.pathPrefix != expected {
		t.Errorf("WithVersion client.pathPrefix = %s, expected %s", c.pathPrefix, expected)
	}
}

func TestWithVersionInvalidVersion(t *testing.T) {
	c := MustNewClient(app, "fooshop", "abcd", WithVersion("9999-99b"))
	expected := "admin"
	if c.pathPrefix != expected {
		t.Errorf("WithVersion client.pathPrefix = %s, expected %s", c.pathPrefix, expected)
	}
}

func TestWithRetry(t *testing.T) {
	c := MustNewClient(app, "fooshop", "abcd", WithRetry(5))
	expected := 5
	if c.retries != expected {
		t.Errorf("WithRetry client.retries = %d, expected %d", c.retries, expected)
	}
}

func TestWithLogger(t *testing.T) {
	logger := &LeveledLogger{Level: LevelDebug}
	c := MustNewClient(app, "fooshop", "abcd", WithLogger(logger))

	if c.log != logger {
		t.Errorf("WithLogger expected logs to match %v != %v", c.log, logger)
	}
}

func TestWithUnstableVersion(t *testing.T) {
	c := MustNewClient(app, "fooshop", "abcd", WithVersion(UnstableApiVersion))
	expected := fmt.Sprintf("admin/api/%s", UnstableApiVersion)
	if c.pathPrefix != expected {
		t.Errorf("WithVersion client.pathPrefix = %s, expected %s", c.pathPrefix, expected)
	}
}

func TestWithHTTPClient(t *testing.T) {
	c := MustNewClient(app, "fooshop", "abcd", WithHTTPClient(&http.Client{Timeout: 30 * time.Second}))
	expected := 30 * time.Second

	if c.Client.Timeout.String() != expected.String() {
		t.Errorf("WithVersion client.Client = %s, expected %s", c.Client.Timeout, expected)
	}
}
