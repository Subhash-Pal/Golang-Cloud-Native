# CONFIG
$url = "http://localhost:8080/login"
$username = "shubh"
$password = "1234"

# LOOP
for ($i = 1; $i -le 3; $i++) {

    Write-Host "`n=== Attempt $i (SUCCESS CASE) ==="

    $body = @{
        username = $username
        password = $password
    } | ConvertTo-Json

    try {
        $response = Invoke-RestMethod `
            -Uri $url `
            -Method POST `
            -Body $body `
            -ContentType "application/json" `
            -ErrorAction Stop

        Write-Host "SUCCESS RESPONSE:"
        $response | ConvertTo-Json -Depth 3
    }
    catch {
        Write-Host "ERROR OCCURRED"

        if ($_.Exception.Response -ne $null) {
            $status = $_.Exception.Response.StatusCode.value__
            Write-Host "Status Code: $status"

            $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
            $respBody = $reader.ReadToEnd()

            Write-Host "Response Body:"
            Write-Host $respBody
        }
        else {
            Write-Host "Request failed completely:"
            Write-Host $_
        }
    }

    Start-Sleep -Seconds 1
}