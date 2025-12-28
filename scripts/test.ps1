Param(
    [string]$CoverOut = "coverage.out",
    [string]$HtmlOut  = "coverage.html"
)

Write-Host "Running go tests with coverage..."
go test ./... -coverprofile=$CoverOut
if ($LASTEXITCODE -ne 0) { Write-Error "go test failed"; exit $LASTEXITCODE }

Write-Host "Coverage summary:"
go tool cover -func=$CoverOut

Write-Host "Generating HTML coverage report: $HtmlOut"
go tool cover -html=$CoverOut -o $HtmlOut

if ($env:CI -eq $null) {
    Try {
        Start-Process $HtmlOut -ErrorAction SilentlyContinue
    } Catch {
        Write-Host "Couldn't open HTML report automatically. File: $HtmlOut"
    }
}

Write-Host "Done."
