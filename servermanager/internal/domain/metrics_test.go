package domain

import "testing"

func TestParseManagedName(t *testing.T) {
	const id = "8b9f5e9e-9a3a-4b7e-9a59-1d2f3a4b5c6d"

	cases := []struct {
		name     string
		in       string
		wantOK   bool
		wantKind ContainerKind
	}{
		{"app container", "stt-app-" + id, true, KindApp},
		{"db container", "stt-db-" + id, true, KindDB},
		{"docker list leading slash", "/stt-app-" + id, true, KindApp},
		{"foreign container", "postgres", false, KindApp},
		{"prefix only", "stt-app-", false, KindApp},
		{"non-uuid remainder", "stt-app-not-a-uuid", false, KindApp},
		{"db volume name", "stt-db-data-" + id, false, KindApp},
		{"substring match from the name filter", "my-stt-app-" + id, false, KindApp},
		{"trailing garbage", "stt-app-" + id + "-old", false, KindApp},
		{"32-hex uuid form", "stt-app-8b9f5e9e9a3a4b7e9a591d2f3a4b5c6d", false, KindApp},
		{"network name", "stt-net-" + id, false, KindApp},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := ParseManagedName(tc.in)
			if ok != tc.wantOK {
				t.Fatalf("ParseManagedName(%q) ok = %v, want %v", tc.in, ok, tc.wantOK)
			}
			if !tc.wantOK {
				return
			}
			if got.Kind != tc.wantKind || got.OwnerID != id {
				t.Errorf("ParseManagedName(%q) = %+v, want kind %q owner %s", tc.in, got, tc.wantKind, id)
			}
			if got.Name == "" || got.Name[0] == '/' {
				t.Errorf("Name = %q, want the slash-stripped container name", got.Name)
			}
		})
	}
}
