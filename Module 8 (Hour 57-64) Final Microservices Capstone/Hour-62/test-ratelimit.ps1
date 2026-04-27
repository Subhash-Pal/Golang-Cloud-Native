for ($i = 1; $i -le 7; $i++) {
    Write-Host "Attempt $i"

    $body = @{
        username = "shubh"
        password = "wrongpass"
    } | ConvertTo-Json

    try {
        $response = Invoke-RestMethod -Uri "http://localhost:8080/login" `
            -Method POST `
            -Body $body `
            -ContentType "application/json"

        $response | ConvertTo-Json -Depth 3
    }
    catch {
        if ($_.Exception.Response) {
            $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
            $respBody = $reader.ReadToEnd()
            Write-Host "Error Response: $respBody"
        } else {
            Write-Host "Request failed: $_"
        }
    }

    Start-Sleep -Seconds 1
}