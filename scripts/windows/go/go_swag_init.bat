cd ..\..\..\

setlocal

set "BIN_DIR=bin"

set "SWAG_CMD=%BIN_DIR%\swag.exe"
set "SWAG_VERSION=v1.16.6"
set "SWAG_PACKAGE=github.com/swaggo/swag/cmd/swag@%SWAG_VERSION%"

mkdir "%BIN_DIR%"

set "GOBIN=%CD%\%BIN_DIR%"

go install %SWAG_PACKAGE%

%SWAG_CMD% init -g cmd/server/main.go -o docs/swagger/server

pause