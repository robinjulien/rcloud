package filemanager

import (
	"strings"
	"testing"
)

func TestSanitizePath(t *testing.T) {
	var paths []string = []string{
		"..",
		"./..",
		"/",
		"../..",
		"/..",
		"~",
		"~/..",
	}

	var sPaths []string = make([]string, len(paths))

	for i, path := range paths {
		sPaths[i] = SanitizePath(path)
	}

	for _, spath := range sPaths {
		//rootre, _ := regexp.MatchString("^/", spath)
		if strings.Contains(spath, "..") || strings.Contains(spath, "~") || len(spath) == 0 || spath[0] == '/' || spath == "" {
			t.Errorf("%q wrong in SanitizePath", spath)
		}
	}

	if SanitizePath("/") != "." {
		t.Errorf("/ =!> . in SanitizePath")
	}

	if SanitizePath(".") != "." {
		t.Errorf(". =!> . in SanitizePath")
	}

	if SanitizePath("") != "." {
		t.Errorf("space =!> . in SanitizePath")
	}
}
