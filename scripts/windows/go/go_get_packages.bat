@echo off

setlocal enabledelayedexpansion

pushd ..\..\..\ || (
  echo [ERR] Cannot cd ..\..\..\
  pause
  exit /b 1
)

set "PACKAGES_ENV_FILE=scripts\env\.packages.env"

if not exist "%PACKAGES_ENV_FILE%" (
  echo [ERR] File %PACKAGES_ENV_FILE% not found
  pause
  exit /b 1
)

echo [INF] File %PACKAGES_ENV_FILE% found

echo [INF] Start getting packages

for /f "usebackq tokens=* delims=" %%L in ("%PACKAGES_ENV_FILE%") do (
  set "LINE=%%L"

  if not "!LINE!"=="" if not "!LINE:~0,1!"=="#" (
    echo [INF] Run: go get !LINE!
    go get "!LINE!"
    if errorlevel 1 (
      echo [ERR] Failed: "!LINE!"
      popd
      pause
      exit /b 1
    )
  )
)

popd

echo [INF] Done

endlocal

@echo on

pause
