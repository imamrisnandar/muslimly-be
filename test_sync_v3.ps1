$ErrorActionPreference = "Stop"

Write-Output "1. Login..."
$loginBody = Get-Content login.json -Raw
$loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
$token = $loginResponse.data.token
Write-Output "Token Obtained."

Write-Output "2. Bulk Insert Activities..."
$activityBody = Get-Content activity_test.json -Raw
$activityResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/sync/activity" -Method Post -Body $activityBody -ContentType "application/json" -Headers @{Authorization = "Bearer $token" }
Write-Output "Response: $($activityResponse.message)"

if ($activityResponse.code -eq 200) {
    Write-Output "TEST PASSED: Activity Log Synced."
}
else {
    Write-Error "TEST FAILED: Code $($activityResponse.code)"
}
