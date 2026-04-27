
This code is a practical demonstration of Token-Based Authentication. It shows how modern apps keep a user logged in without asking for a password every few minutes.
Here is the breakdown of what this example teaches you:


#1. The Two-Token Strategy
Instead of just one token, this uses two. This is a standard security practice:
Access Token (Short-lived): In your code, it lasts only 5 minutes. This is what the user sends to the server to prove who they are. If it gets stolen, it's only useful for a very short time.
Refresh Token (Long-lived): In your code, it lasts 24 hours. Its only job is to get a new Access Token when the old one expires. This stays "hidden" and isn't sent with every request.

2. How the "Handshake" Works

The flow of this code mimics a real-world app:
Login: You provide credentials (the code just assumes "admin"). The server gives you both tokens.
Accessing Data: You use the Access Token.
Expiration: After 5 minutes, the Access Token stops working.
The Refresh: Instead of making the user type their password again, the app automatically sends the Refresh Token to the /refresh endpoint to get a fresh 5-minute Access Token.

3. JWT Structure (JSON Web Tokens)
The code shows how a JWT is built. Every token created here contains:
Payload (Claims): The username and the exp (expiration timestamp).
Signature: The secret variable is used to "sign" the token. This ensures that if a user tries to change their username from "user" to "admin" inside the token, the signature will break and the server will reject it.

4. Statelessness
Notice there is no database being used to store these tokens. The server doesn't need to "remember" you. Every time you send a token, the server just checks the signature against its secret. If the math checks out, you are logged in. This makes the server very fast and easy to scale.

5. Practical Web Development with Gin
It shows how to:
Set up a Web Server in Go.
Create POST routes.
Handle JSON input (BindJSON) and JSON output (c.JSON).

What’s missing? In a real app, you would add a Middleware—a "security guard" function that checks the Access Token before allowing someone to visit a page like /profile or /dashboard.


Use code with caution.

How to test this in PowerShell:
Restart the server: Press Ctrl+C in your terminal and run go run main.go again.

Login to get a new token:
powershell
$loginResponse = Invoke-RestMethod -Method Post -Uri http://localhost:8080/login

$loginResponse # This shows you the access and refresh tokens
Use code with caution.

Refresh using the stored variable:
powershell
$body = @{ refresh = $loginResponse.refresh } | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri http://localhost:8080/refresh -ContentType "application


How to test this in PowerShell:
Restart the server: Press Ctrl+C in your terminal and run go run main.go again.
Login to get a new token:

bash
```
powershell

$loginResponse = Invoke-RestMethod -Method Post -Uri http://localhost:8080/login

$loginResponse # This shows you the access and refresh tokens

Use code with caution.

Refresh using the stored variable:

powershell

$body = @{ refresh = $loginResponse.refresh } | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri http://localhost:8080/refresh -ContentType "application/json" -Body $body

Use code with caution.

```
