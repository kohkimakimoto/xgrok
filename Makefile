.PHONY: default dev dist packaging packaging_destroy fmt test testv deps deps_update

default: dev

dev:
	@bash -c $(CURDIR)/_build/dev.sh

dist:
	@bash -c $(CURDIR)/_build/dist.sh

packaging:
	@bash -c $(CURDIR)/_build/packaging.sh

packaging_destroy:
	@sh -c "cd $(CURDIR)/_build/packaging/rpm && vagrant destroy -f"

fmt:
	go fmt $$(go list ./... | grep -v vendor)

deps:
	gom install

deps_update:
	rm Gomfile.lock; rm -rf vendor; gom install && gom lock

build_assets:
	./vendor/bin/go-bindata -nomemcopy -pkg=assets \
		-o=xgrok/server/assets/assets.go \
		assets/server/...
	./vendor/bin/go-bindata -nomemcopy -pkg=assets \
		-o=xgrok/client/assets/assets.go \
		assets/client/...
