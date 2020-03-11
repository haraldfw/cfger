#!make

SRC_DIRS := .

version:
	@echo $(VERSION)

install-tools:
	./build/install-tools.sh

env:
	env

fmt:
	gofmt -w $(SRC_DIRS)

test: build-dirs
	./build/test.sh $(SRC_DIRS)

lint: lint-all

lint-all:
	revive -config revive.toml -formatter friendly -exclude vendor/... ./...

lint-update:
	./build/lint-update.sh

mods: mod
mod:
	GOSUMDB=off ./build/mod.sh

watch:
	reflex --start-service=true -r '\.go$$' make

watch-tests: watch-test
watch-test:
	reflex --start-service=true -r '\.go$$' make test

watch-all:
	reflex --start-service=true -r '((\.md)|(_title)|(\.go))$$' -R '/resources.go$$' make test
bin-clean:
	rm -rf .go bin
