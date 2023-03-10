package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestDownloadFileByURL(t *testing.T) {
	t.Setenv(envVarGitHubWorkspace, t.TempDir())
	testFileContents := "someFileContents"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, testFileContents)
	}))
	defer ts.Close()

	pathToSavedFile, err := downloadFileByURL(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	contents, err := os.ReadFile(pathToSavedFile)
	if err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(string(contents)) != testFileContents {
		t.Fatal("File contents did not match expected test file contents")
	}
}

func TestGetRequiredEnv(t *testing.T) {
	testCases := []struct {
		desc             string
		envVarNameToTest string
		shouldSucceed    bool
	}{
		{
			desc:             "Returns env vars that are set",
			envVarNameToTest: "foo",
			shouldSucceed:    false,
		},
		{
			desc:             "Errors when an env var is not set",
			envVarNameToTest: "bar",
			shouldSucceed:    true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.shouldSucceed {
				t.Setenv(tC.envVarNameToTest, "someNonemptyValue")
			}

			val, err := getRequiredEnv(tC.envVarNameToTest)

			if tC.shouldSucceed && (err != nil) {
				t.Fatalf("expected '%s' to exist but error was returned: %v", tC.envVarNameToTest, err)
			} else if !tC.shouldSucceed && (val != "") {
				t.Fatalf("expected '%s' to error but a value was returned: %s", tC.envVarNameToTest, val)
			}
		})
	}
}
