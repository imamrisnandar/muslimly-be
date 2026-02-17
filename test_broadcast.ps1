$ErrorActionPreference = "Stop"

Write-Output "1. Login..."
$loginBody = Get-Content login.json -Raw
$loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
$token = $loginResponse.data.token
Write-Output "Token Obtained."

Write-Output "2. Triggering Broadcast..."
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/notifications/test-broadcast" -Method Post -ContentType "application/json" -Headers @{Authorization = "Bearer $token" }
    Write-Output "Response: $($response.message)"
    Write-Output "TEST PASSED: Broadcast sent to registered devices."
}
catch {
    Write-Error "TEST FAILED: $($_.Exception.Message)"
}
