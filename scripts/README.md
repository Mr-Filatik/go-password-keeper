# Description

This folder contains startup and build scripts, and auxiliary CLI utilities for working with the project during development.

## Structure

The structure of this folder is as follows:

* 📂 `scripts`
    * 📂 `env`. Contains common variables for all scripts;
        * 📄 `packages.env`. Contains all packages and their versions used in the module;
    * 📂 `windows`. Contains scripts for the Windows operating system;
        * 📂 `app`. Contains scripts for running and building applications;
        * 📂 `git`. Contains scripts for working with the `git` commands;
            * 📄`git_clean_interactive.bat`. Allows you to delete temporary files interactively;
        * 📂 `go`. Contains scripts for working with the `Go` language utilities;
            * 📄 `go_fmt.bat`. Formats the code;
            * 📄 `go_get_packages.bat`. Loads the required packages as module dependencies (packages and their versions are specified in `packages.env`);
            * 📄 `go_mod_init.bat`. Initializes the module;
            * 📄 `go_mod_tidy.bat`. Updates the module dependencies.
