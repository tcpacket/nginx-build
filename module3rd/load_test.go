package module3rd

import (
	"testing"
)

func TestModules3rd(t *testing.T) {

	modules3rdConf := "../config/modules.cfg.example"
	modules3rd, err := Load(modules3rdConf)
	if err != nil {
		t.Fatalf("Failed to load %s", modules3rdConf)
	}

	for _, m := range modules3rd {
		var want Module3rd
		switch m.Name {
		case "ngx_http_hello_world":
			want.Name = "ngx_http_hello_world"
			want.Form = "git"
			want.Url = "https://github.com/cubicdaiya/ngx_http_hello_world"
			want.Dynamic = false
		default:
			t.Fatalf("unexpected module: %v", m)
		}

		if m != want {
			t.Fatalf("got: %v, want: %v", m, want)
		}
	}
}
