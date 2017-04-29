package awslogin

import "testing"

func TestVersionCheck(t *testing.T) {
	outdate, _, _ := VersionCheck("0.0.1")
	expected := true
	if !outdate {
		t.Errorf("expected %v to eq %v", outdate, expected)
	}

	outdate, _, _ = VersionCheck("100.0.0")
	expected = false
	if outdate {
		t.Errorf("expected %v to eq %v", outdate, expected)
	}
}

func TestCheckArgProfileName(t *testing.T) {
	f := CheckArgProfileName("")
	expected := false
	if f {
		t.Errorf("expected %v to eq %v", f, expected)
	}

	tr := CheckArgProfileName("testProfileName")
	expected = true
	if !tr {
		t.Errorf("expected %v to eq %v", tr, expected)
	}
}
