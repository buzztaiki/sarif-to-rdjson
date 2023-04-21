- [SARIF standardization #1175](https://github.com/reviewdog/reviewdog/issues/1175)
- TODO test tflint error
- TODO add source from tool.driver
- https://github.com/terraform-linters/tflint/blob/master/formatter/sarif.go
- ansible-lint pattern
  ```
    {
      "tool": {
        "driver": {
          "name": "ansible-lint",
          "version": "6.14.5.dev0",
          "informationUri": "https://github.com/ansible/ansible-lint",
          "rules": [
            {
              "id": "no-changed-when",
              "name": "no-changed-when",
              "shortDescription": {
                "text": "Commands should not change things if nothing needs doing."
              },
              "defaultConfiguration": {
                "level": "error"
              },
              "help": {
                "text": ""
              },
              "helpUri": "https://ansible-lint.readthedocs.io/rules/no-changed-when/",
              "properties": {
                "tags": [
                  "command-shell",
                  "idempotency"
                ]
              }
            }
          ]
        }
      },
      "results": [
        {
          "ruleId": "no-changed-when",
          "message": {
            "text": "Task/Handler: XXX"
          },
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "roles/arc/tasks/main.yaml",
                  "uriBaseId": "SRCROOT"
                },
                "region": {
                  "startLine": 24
                }
              }
            }
          ]
        },
  ```
