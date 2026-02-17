$ErrorActionPreference = "Stop"

Write-Output "1. Login..."
$loginBody = Get-Content login.json -Raw
$loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
$token = $loginResponse.data.token
Write-Output "Token Obtained."

Write-Output "2. Upsert Settings..."
$settingsBody = Get-Content settings_test.json -Raw
$settingsResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/sync/settings" -Method Post -Body $settingsBody -ContentType "application/json" -Headers @{Authorization = "Bearer $token" }
Write-Output "Upsert Settings: $($settingsResponse.message)"

Write-Output "3. Get Settings..."
$getSettings = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/sync/settings" -Method Get -ContentType "application/json" -Headers @{Authorization = "Bearer $token" }
Write-Output "Current Settings:"
Write-Output ($getSettings | ConvertTo-Json -Depth 5)

Write-Output "4. Upsert Reading V2..."
$readingBody = Get-Content reading_test_v2.json -Raw
$readingResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/sync/reading" -Method Post -Body $readingBody -ContentType "application/json" -Headers @{Authorization = "Bearer $token" }
Write-Output "Upsert Reading: $($readingResponse.message)"

Write-Output "5. Get Reading History..."
$history = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/sync/reading" -Method Get -ContentType "application/json" -Headers @{Authorization = "Bearer $token" }
Write-Output "Most Recent History:"
Write-Output ($history.data[0] | ConvertTo-Json -Depth 5)

if ($history.data[0].page_number -eq 42 -and $getSettings.data.Length -ge 2) {
    Write-Output "TEST PASSED: Settings and Reading V2 verified."
}
else {
    Write-Error "TEST FAILED: Data mismatch."
}
