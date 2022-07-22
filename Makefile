
COVER_FILE=cover.out
.PHONY: test
test:
	go test -coverprofile=${COVER_FILE}
	go tool cover -html=${COVER_FILE}
	rm -rf ${COVER_FILE}

VERSION=v0.0.3
git-tag:
	git tag $(VERSION); \
	git push origin $(VERSION);