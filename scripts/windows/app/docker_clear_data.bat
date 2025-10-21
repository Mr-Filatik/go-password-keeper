cd ..\..\..\

@echo off

IF EXIST "deploy\docker\data" (
    echo Deleting the deploy\docker\data folder...
    rmdir /s /q "deploy\docker\data"
    echo The folder was deleted successfully.
) ELSE (
    echo Folder deploy\docker\data not found.
)

@echo on

pause
