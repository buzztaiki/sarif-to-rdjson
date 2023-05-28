TFLINT = go run github.com/terraform-linters/tflint@latest
TFSEC = go run github.com/aquasecurity/tfsec/cmd/tfsec@latest
ANSIBLE_LINT = python -mpipx run ansible-lint

default: check build test

check:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

build:
	go build ./...

update-testcases:
	(cd testcases/tflint/simple && $(TFLINT) --init && $(TFLINT) -f sarif --force > sarif.json)
	(cd testcases/tflint/syntax-error && $(TFLINT) -f sarif --force > sarif.json || true)
	(cd testcases/tfsec && $(TFSEC) -f sarif --soft-fail > sarif.json)
	(cd testcases/ansible-lint && $(ANSIBLE_LINT) -f sarif -q | jq > sarif.json || true)
	(cd testcases/codeql && rm -rf .codeql && codeql database create -l go -q .codeql &&  codeql database analyze --format sarif-latest -o sarif.json -q .codeql)


	go test ./...

test:
	go test ./...

ci: check build test
