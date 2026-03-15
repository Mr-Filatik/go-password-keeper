cd ..\..\..\

setlocal

set "BIN_DIR=bin"

set "SWAG_CMD=%BIN_DIR%\swag.exe"
set "SWAG_VERSION=v1.16.6"
set "SWAG_PACKAGE=github.com/swaggo/swag/cmd/swag@%SWAG_VERSION%"

mkdir "%BIN_DIR%"

set "GOBIN=%CD%\%BIN_DIR%"

go install %SWAG_PACKAGE%

%SWAG_CMD% init --dir ./internal/server/http --generalInfo server.go --output docs/swagger/server/ --parseInternal

pause

rem    --markdownFiles ./docs/markdown/server ^
rem    --dir ./internal,./docs
rem    -g server/http/server_docs.go ^