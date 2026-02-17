$ErrorActionPreference = "Stop"

Write-Output "1. Login..."
$loginBody = Get-Content login.json -Raw
$loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method Post -Body $loginBody -ContentType "application/json"

if ($loginResponse.code -ne 200) {
    Write-Error "Login failed"
    exit 1
}

$token = $loginResponse.data.token
Write-Output "Login Success. Token length: $($token.Length)"

Write-Output "2. Upsert Reading Progress..."
$syncBody = Get-Content sync_test.json -Raw
$syncResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/sync/reading" -Method Post -Body $syncBody -ContentType "application/json" -Headers @{Authorization="Bearer $token"}

Write-Output "Upsert Response: $($syncResponse | ConvertTo-Json -Depth 5)"

Write-Output "3. Get Reading History..."
$historyResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/sync/reading" -Method Get -ContentType "application/json" -Headers @{Authorization="Bearer $token"}

Write-Output "History Response:"
Write-Output ($historyResponse | ConvertTo-Json -Depth 5)

if ($historyResponse.data[0].surah_id -eq 1 -and $historyResponse.data[0].ayah_number -eq 5) {
    Write-Output "TEST PASSED: Data verified."
} else {
    Write-Error "TEST FAILED: Data mismatch."
}
