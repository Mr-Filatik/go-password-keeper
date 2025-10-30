[Main documentation page](../README.md)

# Description

This folder contains startup and build scripts, and auxiliary CLI utilities for working with the project during development.

## Structure

The structure of this folder is as follows:

* 📂 `scripts`
    * 📂 `env`. Contains common variables for all scripts;
        * 📄 `.packages.env`. Contains all packages and their versions used in the module;
        * 📄 `.utils.env`. Contains all utilities and their versions used in working with the project;
    * 📂 `windows`. Contains scripts for the Windows operating system;
        * 📂 `app`. Contains scripts for running and building applications;
        * 📂 `git`. Contains scripts for working with the `git` commands;
            * 📄`git_clean_interactive.bat`. Allows you to delete temporary files interactively;
        * 📂 `go`. Contains scripts for working with the `Go` language utilities;
            * 📄 `go_fmt.bat`. Formats the code;
            * 📄 `go_get_packages.bat`. Loads the required packages as module dependencies (packages and their versions are specified in `.packages.env`);
            * 📄 `go_install_utils.bat`. Installs utilities (utilities and their versions are specified in `.utils.env`);
            * 📄 `go_lint.bat`. Runs a linter on the code;
            * 📄 `go_lint_cache_clean.bat`. Clears the linter cache;
            * 📄 `go_mod_init.bat`. Initializes the module;
            * 📄 `go_mod_tidy.bat`. Updates the module dependencies;
            * 📄 `go_test.bat`. Runs all module tests;
            * 📄 `go_test_coverage.bat`. Calculates code coverage of modules and saves it to the `coverage.out` file;
            * 📄 `go_test_coverage_func.bat`. Reads data from the `coverage.out` file and displays the coverage for each function and the percentage of test coverage in the entire module;
            * 📄 `go_test_coverage_html.bat`. Reads data from the `coverage.out` file and displays which sections of code are covered by tests via a web interface.
