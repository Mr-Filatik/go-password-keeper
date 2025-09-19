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
