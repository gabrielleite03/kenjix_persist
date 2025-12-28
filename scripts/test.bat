@echo off
setlocal

set "COVER_OUT=coverage.out"
set "HTML_OUT=coverage.html"

echo Running go tests with coverage...
go test ./... -coverprofile=%COVER_OUT%
if errorlevel 1 (
  echo go test failed
  endlocal
  exit /b 1
)

echo Coverage summary:
go tool cover -func=%COVER_OUT%

echo Generating HTML report: %HTML_OUT%
go tool cover -html=%COVER_OUT% -o %HTML_OUT%

if exist %HTML_OUT% (
  start "" "%HTML_OUT%" 2>nul
)

endlocal
exit /b 0
