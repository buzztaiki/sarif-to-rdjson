TFLINT = go run github.com/terraform-linters/tflint@latest
TFSEC = go run github.com/aquasecurity/tfsec/cmd/tfsec@latest
ANSIBLE_LINT = python -mpipx run ansible-lint

TESTRD_REPORTER = local


default: check build test

check:
	go vet .
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

build:
	go build .

update_testcases:
	(cd testcases/tflint/simple && $(TFLINT) --init && $(TFLINT) -f sarif --force > sarif.json)
	(cd testcases/tflint/syntax-error && $(TFLINT) -f sarif --force > sarif.json || true)
	(cd testcases/tfsec && $(TFSEC) -f sarif --soft-fail > sarif.json)
	(cd testcases/ansible-lint && $(ANSIBLE_LINT) -f sarif -q | jq > sarif.json || true)
	(cd testcases/codeql && rm -rf .codeql && codeql database create -l go -q .codeql &&  codeql database analyze --format sarif-latest -o sarif.json -q .codeql)

test:
	go test .

ci: check build test

test_reviewdog: build
	$(MAKE) _test_reviewdog testcase=testcases/tflint/simple name=tflint
	$(MAKE) _test_reviewdog testcase=testcases/tflint/syntax-error name=tflint
	$(MAKE) _test_reviewdog testcase=testcases/tfsec name=tfsec
	$(MAKE) _test_reviewdog testcase=testcases/ansible-lint name=ansible-lint
	$(MAKE) _test_reviewdog testcase=testcases/codeql name=codeql

_test_reviewdog:
	cat $(testcase)/sarif.json | ./sarif-to-rdjson \
		| jq '.diagnostics = (.diagnostics | map(.location.path = "$(testcase)/" + .location.path))' \
	  	| reviewdog -f rdjson --filter-mode nofilter --reporter $(TESTRD_REPORTER)
