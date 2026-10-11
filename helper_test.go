package config_test

import (
	"errors"
	"strings"
	"testing"
)

// Testing helpers... Thank you, Go standard library, for doing the absolute bare minimum. It is so "simple" to have to repeat basic stuff in every library ever written.

func assertEqual[T comparable](t *testing.T, expected, actual T, msg ...string) {
	if expected == actual {
		return
	}

	t.Helper()
	t.Errorf("expected %v, got %v; %s", expected, actual, strings.Join(msg, ""))
}

func assertEqualSlices[T comparable](t *testing.T, expected, actual []T, msg ...string) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf("slice length unexpected; expected %d, got %d; %s", len(expected), len(actual), strings.Join(msg, ""))
		return
	}

	for i := range len(expected) {
		if expected[i] != actual[i] {
			t.Errorf("slice mismatch at element %d: expected [%v], got [%v]; %s", i, expected[i], actual[i], strings.Join(msg, ""))
			return
		}
	}
}

func assertEqualMaps[K comparable, V comparable](t *testing.T, expected, actual map[K]V, msg ...string) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf("map length unexpected; expected %d, got %d; %s", len(expected), len(actual), strings.Join(msg, ""))
		return
	}

	for k, v := range expected {
		if v2, ok := actual[k]; !ok {
			t.Errorf("missing key %v in actual map; %s", k, strings.Join(msg, ""))
			return
		} else if v != v2 {
			t.Errorf("value mismatch for key %v; expected %v, got %v; %s", k, v, v2, strings.Join(msg, ""))
		}
	}
}

func assertError(t *testing.T, err error, msg ...string) {
	if err != nil {
		return
	}

	t.Helper()
	t.Errorf("expected error, but got nil; %s", strings.Join(msg, ""))
}

func assertNoError(t *testing.T, err error, msg ...string) {
	if err == nil {
		return
	}

	t.Helper()
	t.Errorf("received unexpected error: %v; %s", err, strings.Join(msg, ""))
}

func assertErrorIs(t *testing.T, err, targetErr error, msg ...string) {
	if errors.Is(err, targetErr) {
		return
	}

	t.Helper()
	t.Errorf("received unexpected error: [%v]; expected [%v]; %s", err, targetErr, strings.Join(msg, ""))
}
