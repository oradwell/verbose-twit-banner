package main

import (
	"errors"
	"regexp"
	"testing"
)

func TestGetRandIndexInArray(t *testing.T) {
	type randIndexTest struct {
		list []string
	}

	tests := []randIndexTest{
		{[]string{"pen"}},
		{[]string{"pom", "sam"}},
		{[]string{}},
	}
	for _, tt := range tests {
		list := tt.list
		rInd, err := getRandIndexInArray(list)
		if len(list) == 0 && err == nil {
			t.Errorf("getRandIndexInArray(%#v) = %d want error (%s)", list, rInd, err)
		}

		if len(list) != 0 && (rInd >= len(list) || rInd < 0) {
			t.Errorf("getRandIndexInArray(%#v) = %d want between %d and %d", list, rInd, 0, len(list))
		}
	}
}

func TestGetJpegPathInDirectory(t *testing.T) {
	type getJpegPathInDirectoryTest struct {
		directory  string
		err        error
		imageRegex string
	}

	tests := map[string]getJpegPathInDirectoryTest{
		"non-existent dir": {"bob", errEmptyDir, ""},
		"images dir":       {"images", nil, ".*\\.jpg"},
	}

	for name, tt := range tests {
		want := regexp.MustCompile(tt.imageRegex)
		t.Run(name, func(t *testing.T) {
			path, err := GetJpegPathInDirectory(tt.directory)
			if tt.err != nil && !errors.Is(err, tt.err) {
				t.Errorf("GetJpegPathInDirectory(%s) = %s with error %q want error %q", tt.directory, path, tt.err, err)
			} else if !want.MatchString(path) {
				t.Errorf("GetJpegPathInDirectory(%s) = %s want match for %#q", tt.directory, path, tt.imageRegex)
			}
		})
	}
}
