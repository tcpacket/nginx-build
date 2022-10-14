export GO111MODULE=on

waf-builder: *.go builder/*.go command/*.go configure/*.go module3rd/*.go openresty/*.go util/*.go
	go build -ldflags "-X main.NginxBuildVersion=`git rev-list HEAD -n1`" -o $@

build-example: waf-builder
	./waf-builder -c config/configure.example -m config/modules.cfg.example -d work -clear

check:
	go test ./...

fmt:
	go fmt ./...

clean:
	rm -rf waf-builder
