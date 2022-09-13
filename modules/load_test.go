package modules

import (
	"testing"

	"github.com/kortschak/utter"
)

func TestModLoad(t *testing.T) {
	target := "../config/modules.json.example"
	mods, err := Load(target)
	if err != nil {
		t.Fatalf("Failed to load %s: %v", target, err)
	}

	for _, m := range mods {
		var want Module
		switch m.Name {
		case "ngx_http_hello_world":
			want.Name = "ngx_http_hello_world"
			want.Form = "git"
			want.URL = "https://github.com/cubicdaiya/ngx_http_hello_world"
			want.Dynamic = false
		default:
			t.Fatalf("unexpected module: %v", m)
		}

		if m != want {
			t.Fatalf("got: \n%v\n---\n want: \n%v", utter.Sdump(m), utter.Sdump(want))
		}
	}
}

func TestModules3rdWithNJS(t *testing.T) {

	modules3rdConf := "../config/modules.json.njs"
	modules3rd, err := Load(modules3rdConf)
	if err != nil {
		t.Fatalf("Failed to load %s", modules3rdConf)
	}

	for _, m := range modules3rd {
		var want Module
		switch m.Name {
		case "njs/nginx":
			want.Name = "njs/nginx"
			want.Form = "hg"
			want.URL = "https://hg.nginx.org/njs"
			want.Dynamic = false
			want.Shprov = "./configure && make"
			want.ShprovDir = ".."
		default:
			t.Fatalf("unexpected module: %v", m)
		}

		if m != want {
			t.Fatalf("got: %v, want: %v", m, want)
		}
	}
}
