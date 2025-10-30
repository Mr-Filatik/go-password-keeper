@echo off

setlocal enabledelayedexpansion

pushd ..\..\..\ || (
  echo [ERR] Cannot cd ..\..\..\
  pause
  exit /b 1
)

set "UTILS_ENV_FILE=scripts\env\.utils.env"

if not exist "%UTILS_ENV_FILE%" (
  echo [ERR] File %UTILS_ENV_FILE% not found
  pause
  exit /b 1
)

echo [INF] File %UTILS_ENV_FILE% found

echo [INF] Start install utils

for /f "usebackq tokens=* delims=" %%L in ("%UTILS_ENV_FILE%") do (
  set "LINE=%%L"

  if not "!LINE!"=="" if not "!LINE:~0,1!"=="#" (
    echo [INF] Run: go install !LINE!
    go install "!LINE!"
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
