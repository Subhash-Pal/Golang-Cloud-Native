go mod init lab3-gin
go get github.com/gin-gonic/gin
🧩 TASK 1 — Basic Modular CRUD API
📁 Structure
lab3-gin/
 ├── main.go
 ├── routes/
 │    └── routes.go
 ├── handlers/
 │    └── user.go
 └── models/
      └── user.go





1. Create a User (POST)
powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/users" `
-Method Post `
-ContentType "application/json" `
-Body '{"name": "John Doe", "email": "john@example.com"}'
Use code with caution.

2. Get All Users (GET)
powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/users" -Method Get
Use code with caution.

3. Get a Specific User (GET)
powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/users/1" -Method Get
Use code with caution.

4. Update a User (PUT)
powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/users/1" `
-Method Put `
-ContentType "application/json" `
-Body '{"name": "John Updated", "email": "updated@example.com"}'
Use code with caution.

5. Delete a User (DELETE)
powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/users/1" -Method Delete
Use code with caution.



