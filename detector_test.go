package wraperr_test

import (
	"testing"

	"github.com/izumin5210/wraperr"
)

func TestDetector_CheckPackages(t *testing.T) {
	detector := wraperr.NewDetector()
	err := detector.CheckPackages([]string{"github.com/izumin5210/wraperr/testdata/detector/simple"})

	if err != nil {
		t.Errorf("returned %v", err)
	}
}
