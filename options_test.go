package goshopify

import (
	"fmt"
	"testing"
	"time"
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

func TestWithUnstableVersion(t *testing.T) {
	c := NewClient(app, "fooshop", "abcd", WithVersion(UnstableApiVersion))
	expected := fmt.Sprintf("admin/api/%s", UnstableApiVersion)
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

func TestWithLogger(t *testing.T) {
	logger := &LeveledLogger{Level: LevelDebug}
	c := NewClient(app, "fooshop", "abcd", WithLogger(logger))

	if c.log != logger {
		t.Errorf("WithLogger expected logs to match %v != %v", c.log, logger)
	}
}

func TestWithSleep(t *testing.T) {
	var called bool
	sleep := func(d time.Duration) { called = true }
	c := NewClient(app, "fooshop", "abcd", WithSleep(sleep))
	c.sleep(time.Duration(1))

	if !called {
		t.Errorf("expected passed function to be called")
	}
}

func TestWithoutSleep(t *testing.T) {
	var called bool
	c := NewClient(app, "fooshop", "abcd")
	c.sleep(time.Duration(1))

	if called {
		t.Errorf("expected called to remain false")
	}
}
