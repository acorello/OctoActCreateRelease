.PHONY: build.linux
build.linux:
	mkdir -p _tmp/built
	GOOS=linux go build -ldflags="-s -w" -trimpath -o _tmp/built/octoact_create_release.linux