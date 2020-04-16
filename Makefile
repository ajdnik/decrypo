GITHUB_REPO="ajdnik/decrypo"
VERSION="0.3.2"

changelog:
	git-chglog -c .chglog/changelog/config.yml -o CHANGELOG.md --next-tag ${VERSION} ..${VERSION}

devdeps:
	go get -u github.com/git-chglog/git-chglog/cmd/git-chglog
	go get -u golang.org/x/tools/cmd/cover

tag: changelog
	git add CHANGELOG.md
	git commit -m "chore: updated changelog"
	git add Makefile
	git commit -m "chore: version bumped"
	git push
	git tag ${VERSION}
	git push origin ${VERSION}

test:
	go test -v ./...

testcov:
	go test -coverprofile coverage.out -coverpkg=./... ./...
	go tool cover -html=coverage.out

default: changelog

.PHONY: changelog devdeps tag test testcov
