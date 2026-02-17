$ErrorActionPreference = "Stop"

Write-Output "1. Login..."
$loginBody = Get-Content login.json -Raw
$loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
$token = $loginResponse.data.token
Write-Output "Token Obtained."

Write-Output "2. Register Device..."
$deviceBody = Get-Content device_test.json -Raw
$deviceResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/notifications/register" -Method Post -Body $deviceBody -ContentType "application/json" -Headers @{Authorization="Bearer $token"}
Write-Output "Response: $($deviceResponse.message)"

if ($deviceResponse.code -eq 200) {
    Write-Output "TEST PASSED: Device Registered."
} else {
    Write-Error "TEST FAILED: Code $($deviceResponse.code)"
}
