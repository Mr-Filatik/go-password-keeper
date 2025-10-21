[Main documentation page](../README.md)

# Description

A dedicated folder for `GitHub`. Contains all workflows, their configs, and other useful files required by `GitHub`.

## Structure

The structure of this folder is as follows:

* 📂 `.github`
    * 📂 `configs`. Contains configs for `workflow`;
        * 📄 `.golangci.yml`. Contains the configuration for `golangci-lint`;
        * 📄 `.labeler.yml`. Contains the configuration for `labeler`;
    * 📂 `workflows`. Contains all workflows used in the project;
        * 📄 `golangci-lint.yml`. Analyzes code for compliance with certain rules (the rules are located in `.golangci.yml`);
        * 📄 `labeler.yml`. Sets tags for changes in the current Pull Request (rules are located in `.labeler.yml`);
        * 📄 `test-coverage.yml`. Checks code coverage by tests;
    * 📄 `pull_request_template.md`. Special file. Describes a comment template that is added when creating a Pull Request in Markdown format.

## Golangci-Lint Settings

The linter settings are taken from the official [documentation](https://golangci-lint.run/docs/configuration/file/) and mostly include the default settings, but with some modifications. The version in use at the time of writing the documentation is `2.4.0`.

### Changes to linter settings

The following linters are either disabled or have had their settings changed:

* ⚠️ `depguard` - settings changed

  For all `main.go` executables launched from the `cmd` folder, the use of any packages other than those from `internal` that provide a method for launching the application is prohibited. For more details, see the [documentation](https://golangci-lint.run/docs/linters/configuration/#depguard).

  ```yml
    depguard:
      rules:
        main:
          files: # Default: $all
            - "**/cmd/**/main.go"
          allow: # Default: []
            - "github.com/mr-filatik/go-password-keeper/internal/client"
            - "github.com/mr-filatik/go-password-keeper/internal/server"
  ```

* ⚠️ `dupl` - exception added

  An exception has been added to tests that removes the check for code duplication. More details in the [documentation](https://golangci-lint.run/docs/configuration/file/#linters-configuration).

  ```yml
      - path: _test\.go
        linters:
          - dupl
  ```

* ⚠️ `exhaustruct` - settings changed

  Added exceptions for checking entities whose names begin with `fake*`, `Fake*`, `mock*`, and `Mock*`. These are used only in tests. Because the type name includes packages, you need to search among all the `(?i)^(?:.*/)?(?:[^.]+\.)?` prefixes. More details in the [documentation](https://golangci-lint.run/docs/linters/configuration/#exhaustruct).

  ```yml
    exhaustruct:
      exclude: # Default: []
        - '(?i)^(?:.*/)?(?:[^.]+\.)?(mock|Mock)\w*$'
        - '(?i)^(?:.*/)?(?:[^.]+\.)?(fake|Fake)\w*$'
  ```

* ⚠️ `forbidigo` - settings changed

  The default value `^(fmt\\.Print(|f|ln)|print|println)$` has been modified to prohibit use of the `print()` function. More details in the [documentation](https://golangci-lint.run/docs/linters/configuration/#forbidigo).

  ```yml
    forbidigo:
      forbid: # Default: ["^(fmt\\.Print(|f|ln)|print|println)$"]
       - pattern: "^(fmt\\.Print(|f|ln)|print|println)$"
         msg: Do not use fmt.Print, fmt,Printf and fmt.Println for log output.
       - pattern: "^print(ln)?$"
         msg: Do not use print() for log output.
  ```

* ⚠️ `funlen` - exception added

  An exception has been added to tests that removes checks for functions that begin with `getTests`. These functions are used to generate test data. More details in the [documentation](https://golangci-lint.run/docs/configuration/file/#linters-configuration).

  ```yml
      - path: _test\.go
        linters:
          - funlen
        text: "^Function 'getTests"
  ```

* ⚠️ `varnamelen` - settings changed

  Added `tt` to the names of variables used in table-driven tests. Also, `w http.ResponseWriter` and `r *http.Request` used in handlers are described. And an exception for `fs *flag.FlagSet`. More details in the [documentation](https://golangci-lint.run/docs/linters/configuration/#varnamelen).

  ```yml
    varnamelen:
      ignore-names: # Default: []
        - tt # using in table-driven tests
      ignore-decls: # Default: []
        - fs *flag.FlagSet # using in configs
        - r *http.Request # using in http handlers
        - w http.ResponseWriter # using in http handlers
  ```

* ❌ `wcl` - removed

  The linter 'wsl' is deprecated (since v2.2.0) due to: new major version. Replaced by wsl_v5.

Additionally, exclusions have been added for generated files:

* ⚠️ Swagger files

  Some linters have been added to the exclusions. More details in the [documentation](https://golangci-lint.run/docs/configuration/file/#linters-configuration).

  ```yml
      - path: ^docs/swagger/
        linters:
          - gochecknoinits
          - gochecknoglobals
          - godot
  ```
