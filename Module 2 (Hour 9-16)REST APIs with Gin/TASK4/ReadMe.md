1. How to Run the Server
Initialize your module (if you haven't already):
powershell
go mod init registration-app
go get github.com/gin-gonic/gin
Use code with caution.

Start the application:
Save your code as main.go and run:
powershell
go run main.go
Use code with caution.

The server will start listening on http://localhost:8080. 
The Go Programming Language
The Go Programming Language
 +3
2. How to Test (Success Case)
Use this PowerShell command to send valid data. Note that age must be 18 or higher to pass the gte=18 rule. 
Naukri.com
Naukri.com
powershell
$body = @{ 
    name = "John Doe"
    email = "john@example.com"
    age = 25 
} | ConvertTo-Json

Invoke-RestMethod -Method Post -Uri http://localhost:8080/register -ContentType "application/json" -Body $body


Use code with caution.

Expected Response: {"message":"Valid input"} 

3. How to Test (Validation Failure)
To see how the binding tags work, send invalid data (e.g., an age under 18 or an invalid email format). 

powershell

# 1. Create data that violates multiple rules:
# - 'name' is missing (required)
# - 'email' is in an invalid format (email)
# - 'age' is 15 (gte=18)
$badData = @{ 
    email = "not-an-email"
    age   = 15 
} | ConvertTo-Json

# 2. Run the test with error handling to avoid the "red text" crash
try {
    Invoke-RestMethod -Method Post -Uri http://localhost:8080/register -ContentType "application/json" -Body $badData
} catch {
    # Extract the actual JSON error sent by your Go server
    $respStream = $_.Exception.Response.GetResponseStream()
    $reader = New-Object System.IO.StreamReader($respStream)
    Write-Host "Server Response (Bad Input):" -ForegroundColor Yellow
    $reader.ReadToEnd()
}



What this test confirms
When you run this, your Go server will catch these issues via the c.ShouldBindJSON(&input) line. The response will look like this: 
Medium
Medium
Name: Key: 'RegisterInput.Name' Error:Field validation for 'Name' failed on the 'required' tag.
Email: Error:Field validation for 'Email' failed on the 'email' tag.
Age: Error:Field validation for 'Age' failed on the 'gte' tag. 
Medium
Medium
 +2
Summary of common "Bad Input" scenarios:
Missing Field: Send a JSON object without the "name" key.
Invalid Email: Send "email": "john.doe" (missing the @ or domain).
Underage: Send "age": 17 or lower.
Malformed JSON: Send a string that isn't valid JSON (like missing a brace {).