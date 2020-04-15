GITHUB_REPO="ajdnik/decrypo"
VERSION="0.3.2"

changelog:
	git-chglog -c .chglog/changelog/config.yml -o CHANGELOG.md --next-tag ${VERSION} ..${VERSION}

devdeps:
	go get -u github.com/git-chglog/git-chglog/cmd/git-chglog

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
	go test -v -coverprofile coverage.out -coverpkg=./... ./... 

default: changelog

.PHONY: changelog devdeps tag test testcov
