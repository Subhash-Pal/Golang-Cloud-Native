Certainly! Below is a well-structured `README.md` file for your project that includes details about the application, how to build and run it, and notes on debugging the port conflict issue.

---

# Go Docker Application

This repository contains a simple Go-based HTTP server containerized using Docker. The application serves as a learning exercise for Docker fundamentals, multi-stage builds, and Docker Compose.

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Project Structure](#project-structure)
3. [Building and Running the Application](#building-and-running-the-application)
4. [Debugging Port Conflicts](#debugging-port-conflicts)
5. [Contributing](#contributing)
6. [License](#license)

---

## Prerequisites

Before you begin, ensure you have the following installed on your system:
- **Docker Desktop**: Install from [here](https://www.docker.com/products/docker-desktop/).
- **Go (Optional)**: Required only if you want to test the application locally without Docker.
- **WSL 2 (Windows Users)**: Ensure WSL 2 is enabled and configured for Docker Desktop.

---

## Project Structure

The project has the following structure:

```
.
├── main.go          # The Go application source code.
├── Dockerfile       # Defines the Docker image for the application.
├── docker-compose.yml (optional) # Multi-container setup (if applicable).
└── README.md        # This file.
```

---

## Building and Running the Application

### Step 1: Build the Docker Image

To build the Docker image, run the following command in the project root directory:

```bash
docker build -t go-docker-app .
```

This will create a Docker image named `go-docker-app`.

### Step 2: Run the Docker Container

Run the container using the following command:

```bash
docker run -d -p 9090:8080 go-docker-app
```

- The `-d` flag runs the container in detached mode.
- The `-p 9090:8080` maps the container's port `8080` to the host's port `9090`.

### Step 3: Access the Application

Once the container is running, access the application in your browser:

```
http://localhost:9090
```

You should see the message:

```
Hello, Docker!
```

---

## Debugging Port Conflicts

During development, you may encounter errors like:

```
Bind for 0.0.0.0:8080 failed: port is already allocated
```

This happens when port `8080` is already in use by another process or service (e.g., Docker Desktop itself). Here’s how to resolve it:

### Option 1: Use a Different Host Port

Instead of mapping the container's port `8080` to the host's port `8080`, use a different host port like `9090`:

```bash
docker run -d -p 9090:8080 go-docker-app
```

Access the application at:

```
http://localhost:9090
```

### Option 2: Change Docker Desktop Settings

If Docker Desktop is using port `8080` internally, you can reconfigure it:

1. Open Docker Desktop settings.
2. Navigate to the networking section.
3. Change the default HTTP port from `8080` to another value (e.g., `9090`).
4. Restart Docker Desktop.

### Option 3: Stop Conflicting Processes

Identify and stop processes using port `8080`:

```powershell
netstat -ano | findstr :8080
Stop-Process -Id <PID> -Force
```

Replace `<PID>` with the actual Process ID from the `netstat` output.

---

## Contributing

Contributions are welcome! If you’d like to improve this project, follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -m "Add your feature"`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Open a pull request.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---