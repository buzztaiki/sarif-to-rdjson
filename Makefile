default: check build update-testcases test

check:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

build:
	go build ./...

update-testcases:
	(cd testcases/tflint/simple && tflint -f sarif --force > sarif.json)
	(cd testcases/tflint/syntax-error && tflint -f sarif --force > sarif.json || true)

test:
	go test ./...
