package domain

import (
	"errors"
	"testing"
)

func TestDBImageFor(t *testing.T) {
	img, err := DBImageFor("postgres", "16")
	if err != nil {
		t.Fatalf("DBImageFor(postgres, 16) returned error: %v", err)
	}
	if img.Ref != "postgres:16" || img.DataDir != "/var/lib/postgresql/data" {
		t.Errorf("DBImageFor(postgres, 16) = %+v", img)
	}

	rejected := []struct{ dbType, version string }{
		{"mysql", "8.0"},
		{"mongodb", "7"},
		{"postgres", "15"},
		{"postgres", "latest"},
		{"postgres", ""},
		{"", "16"},
		{"postgres", "16-alpine"},
	}
	for _, tc := range rejected {
		if _, err := DBImageFor(tc.dbType, tc.version); !errors.Is(err, ErrInvalid) {
			t.Errorf("DBImageFor(%q, %q) = %v, want ErrInvalid", tc.dbType, tc.version, err)
		}
	}
}

func TestValidateDBIdent(t *testing.T) {
	for _, ok := range []string{"app", "my_db", "_x", "a1234567890"} {
		if err := ValidateDBIdent("db_name", ok); err != nil {
			t.Errorf("ValidateDBIdent(%q) = %v, want nil", ok, err)
		}
	}
	for _, bad := range []string{"", "1abc", "App", "a-b", "a b", "a;drop", "café", string(make([]byte, 64))} {
		if err := ValidateDBIdent("db_name", bad); !errors.Is(err, ErrInvalid) {
			t.Errorf("ValidateDBIdent(%q) = %v, want ErrInvalid", bad, err)
		}
	}
}

func TestValidateDBPassword(t *testing.T) {
	if err := ValidateDBPassword("0123456789abcdef0123456789abcdef"); err != nil {
		t.Errorf("hex password rejected: %v", err)
	}
	for _, bad := range []string{"", "short", "has spaces padpadpadpad", "newline\npadpadpadpadpad", "tab\tpadpadpadpadpad"} {
		if err := ValidateDBPassword(bad); !errors.Is(err, ErrInvalid) {
			t.Errorf("ValidateDBPassword(%q) = %v, want ErrInvalid", bad, err)
		}
	}
}

func TestDBNames(t *testing.T) {
	if got := DBNetworkName("x"); got != "stt-dbnet-x" {
		t.Errorf("DBNetworkName = %q", got)
	}
	if got := DBVolumeName("x"); got != "stt-db-data-x" {
		t.Errorf("DBVolumeName = %q", got)
	}
	// The db network prefix must never match the app-network sweep prefix.
	if got := AppNetworkName("x"); got == DBNetworkName("x") {
		t.Error("app and db network names collide")
	}
	dbSpec := ContainerSpec{AppID: "x", Kind: KindDB}
	if dbSpec.ContainerName() != "stt-db-x" || dbSpec.NetworkName() != "stt-dbnet-x" {
		t.Errorf("db spec names = %q / %q", dbSpec.ContainerName(), dbSpec.NetworkName())
	}
	appSpec := ContainerSpec{AppID: "x"}
	if appSpec.ContainerName() != "stt-app-x" || appSpec.NetworkName() != "stt-net-x" {
		t.Errorf("app spec names = %q / %q", appSpec.ContainerName(), appSpec.NetworkName())
	}
}
