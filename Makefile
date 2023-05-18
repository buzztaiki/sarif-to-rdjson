default: check build update-testcases test

check:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

build:
	go build ./...

update-testcases:
	(cd testcases/tflint/simple && tflint -f sarif --force > sarif.json)
	(cd testcases/tflint/syntax-error && tflint -f sarif --force > sarif.json || true)
	(cd testcases/tfsec && tfsec -f sarif --soft-fail > sarif.json)

test:
	go test ./...
