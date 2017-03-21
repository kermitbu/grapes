all: test testrace

deps:
	go get -d -v github.com/kermitbu/grapes/grpc/...

updatedeps:
	go get -d -v -u -f github.com/kermitbu/grapes/grpc/...

testdeps:
	go get -d -v -t github.com/kermitbu/grapes/grpc/...

updatetestdeps:
	go get -d -v -t -u -f github.com/kermitbu/grapes/grpc/...

build: deps
	go build github.com/kermitbu/grapes/grpc/...

proto:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	go get -u -v github.com/golang/protobuf/protoc-gen-go
	# use $$dir as the root for all proto files in the same directory
	for dir in $$(git ls-files '*.proto' | xargs -n1 dirname | uniq); do \
		protoc -I $$dir --go_out=plugins=grpc:$$dir $$dir/*.proto; \
	done

test: testdeps
	go test -v -cpu 1,4 github.com/kermitbu/grapes/grpc/...

testrace: testdeps
	go test -v -race -cpu 1,4 github.com/kermitbu/grapes/grpc/...

clean:
	go clean -i github.com/kermitbu/grapes/grpc/...

coverage: testdeps
	./coverage.sh --coveralls

.PHONY: \
	all \
	deps \
	updatedeps \
	testdeps \
	updatetestdeps \
	build \
	proto \
	test \
	testrace \
	clean \
	coverage
