# distkv-cli

`distkv-cli` is a Go command-line application that works with a distributed key-value store built on NATS JetStream.

This README is written so that anyone can run the project by following the steps in order.

## What you are running

This project has two parts:

1. a CLI application written in Go
2. a NATS server with JetStream enabled, running through Docker

Important:

- the CLI is the client
- NATS JetStream is the backend server
- the backend must be started first

If the server is not running, commands like `put`, `get`, `list`, `watch`, and `health` will fail.

## Before you start

Make sure you have these installed:

- Go 1.26 or newer
- Docker Desktop or Docker Engine with Docker Compose
- PowerShell

## Project folder

Open PowerShell and go to this folder:

```powershell
cd "D:\training_golang\Module 11 (Hour 81-88) Building Advanced CLI Applications\Hour88\distkv-cli"
```

## Quick run order

If you only want the shortest correct order, run these commands one by one:

```powershell
docker compose up -d
docker compose ps
go build ./...
go run . bucket create
go run . put app.config.version v1
go run . get app.config.version
go run . list
go run . health
go run . put app.config.version v2
go run . get app.config.version
docker compose down
```

The detailed explanation for each step is below.

## Step 1: Start the backend server

Run:

```powershell
docker compose up -d
```

What this does:

- starts a NATS container
- enables JetStream
- exposes NATS on port `4222`
- exposes monitoring on port `8222`
- makes the server available at `nats://127.0.0.1:4222`
- makes the monitoring page available at `http://127.0.0.1:8222`

Now check whether it is running:

```powershell
docker compose ps
```

What you should see:

- container name similar to `distkv-cli-nats-1`
- status `Up`
- health `healthy`

How to check the server URL:

- NATS client URL: `nats://127.0.0.1:4222`
- NATS monitoring URL: `http://127.0.0.1:8222`

Open the monitoring URL in the browser:

```powershell
start http://127.0.0.1:8222
```

Or check it from PowerShell:

```powershell
Invoke-WebRequest http://127.0.0.1:8222/healthz
```

Expected result:

- status code `200`
- response showing the server is healthy

If you want to use the helper script instead:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\start-backend.ps1
```

## Step 2: Build the CLI

Run:

```powershell
go build ./...
```

What this does:

- downloads dependencies if needed
- compiles the project
- checks that the code is buildable

If you want to use the helper script instead:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\build-cli.ps1
```

If `go run` is blocked on your machine, create an executable:

```powershell
go build -o distkv.exe .
```

Then replace `go run . ...` with:

```powershell
.\distkv.exe ...
```

Example:

```powershell
.\distkv.exe health
```

## Step 3: Create the KV bucket

Run:

```powershell
go run . bucket create
```

What this does:

- connects to the NATS server
- checks whether the JetStream key-value bucket exists
- creates it if missing
- reuses it if it already exists
- creates the distributed storage area named `distkv` inside JetStream

Expected output:

```text
bucket "distkv" created
```

or:

```text
bucket "distkv" reused
```

Effect on the system:

- a JetStream KV bucket named `distkv` is available on the server
- future `put`, `get`, `list`, and `watch` commands use this bucket by default

How to check it on the server:

```powershell
start http://127.0.0.1:8222
```

Then inspect the JetStream information from the monitoring UI.

You can also query the monitoring API from PowerShell:

```powershell
Invoke-WebRequest "http://127.0.0.1:8222/jsz?streams=true" | Select-Object -ExpandProperty Content
```

Look for a JetStream stream related to the KV bucket. For bucket `distkv`, the backing stream is typically named `KV_distkv`.

## Step 4: Store a value

Run:

```powershell
go run . put app.config.version v1
```

What this does:

- stores the key `app.config.version`
- stores the value `v1`
- creates revision `1`
- writes data into the distributed JetStream KV bucket

Expected output:

```text
stored "app.config.version" at revision 1
```

Effect on the system:

- the bucket now contains one key: `app.config.version`
- the current value becomes `v1`
- the current revision becomes `1`

How to check it:

```powershell
go run . get app.config.version
go run . list
```

## Step 5: Read the value

Run:

```powershell
go run . get app.config.version
```

What this does:

- fetches the current value of the key from JetStream KV
- reads the latest stored revision for that key from the server

Expected output:

```text
app.config.version=v1 (rev=1)
```

Effect on the system:

- this command does not change server data
- it only reads and displays the current stored value

## Step 6: List keys

Run:

```powershell
go run . list
```

What this does:

- lists all keys stored in the bucket
- reads the key names currently present in the JetStream KV bucket

Expected output:

```text
app.config.version
```

Effect on the system:

- this command does not change server data
- it only shows which keys currently exist

## Step 7: Run the health check

Run:

```powershell
go run . health
```

What this does:

- checks the connection to the NATS server
- checks whether JetStream is available
- shows which bucket is being used
- confirms that the CLI can talk to the distributed backend

Expected output:

```text
server=nats://127.0.0.1:4222 jetstream=true bucket=distkv
```

If `jetstream=true` appears, the backend is working correctly.

Effect on the system:

- this command does not write anything
- it is only a connectivity and backend status check

## Step 8: Update the value

Run:

```powershell
go run . put app.config.version v2
```

Expected output:

```text
stored "app.config.version" at revision 2
```

Effect on the system:

- the key `app.config.version` changes from `v1` to `v2`
- JetStream stores a new revision
- the latest revision number becomes `2`

Now read it again:

```powershell
go run . get app.config.version
```

Expected output:

```text
app.config.version=v2 (rev=2)
```

This proves the value changed and the revision increased.

## Step 9: Watch live distributed updates

This step shows the distributed behavior.

Open two PowerShell windows.

In the first PowerShell window, run:

```powershell
go run . watch "app.>"
```

This command keeps running and waits for updates.

You can also use:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\watch-demo.ps1
```

In the second PowerShell window, run:

```powershell
go run . put app.config.version v3
```

In the first window, you should see output like:

```text
app.config.version rev=3 op=put value="v3"
```

What this means:

- one client wrote data
- another client immediately saw the change
- the distributed event flow is working

Effect on the system:

- the watcher command itself does not write data
- it subscribes to live key updates from the server
- every future matching `put`, `delete`, or update event is streamed to the console

## Example workflow

Run these commands in order:

```powershell
go run . bucket create
go run . put app.config.version v1
go run . get app.config.version
go run . list
go run . watch "app.>"
```

If you want separate PowerShell files for each example command, use the scripts below.

## Example workflow scripts

Each of these `.ps1` files runs one command from the example workflow:

- [scripts/example-bucket-create.ps1](</D:/training_golang/Module 11 (Hour 81-88) Building Advanced CLI Applications/Hour88/distkv-cli/scripts/example-bucket-create.ps1>): runs `go run . bucket create`
- [scripts/example-put-v1.ps1](</D:/training_golang/Module 11 (Hour 81-88) Building Advanced CLI Applications/Hour88/distkv-cli/scripts/example-put-v1.ps1>): runs `go run . put app.config.version v1`
- [scripts/example-get-version.ps1](</D:/training_golang/Module 11 (Hour 81-88) Building Advanced CLI Applications/Hour88/distkv-cli/scripts/example-get-version.ps1>): runs `go run . get app.config.version`
- [scripts/example-list.ps1](</D:/training_golang/Module 11 (Hour 81-88) Building Advanced CLI Applications/Hour88/distkv-cli/scripts/example-list.ps1>): runs `go run . list`
- [scripts/example-watch-app.ps1](</D:/training_golang/Module 11 (Hour 81-88) Building Advanced CLI Applications/Hour88/distkv-cli/scripts/example-watch-app.ps1>): runs `go run . watch "app.>"`

Run them like this:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\example-bucket-create.ps1
powershell -ExecutionPolicy Bypass -File .\scripts\example-put-v1.ps1
powershell -ExecutionPolicy Bypass -File .\scripts\example-get-version.ps1
powershell -ExecutionPolicy Bypass -File .\scripts\example-list.ps1
powershell -ExecutionPolicy Bypass -File .\scripts\example-watch-app.ps1
```

## Step 10: Run the practical automatically

If you want the main practical commands executed in order, use:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\run-practical.ps1
```

This script runs:

- bucket creation
- first write
- first read
- list
- health check
- second write
- final read

## Step 11: Stop the backend

When you are done, stop the NATS container:

```powershell
docker compose down
```

Or use:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\stop-backend.ps1
```

If you also want to remove the saved JetStream data:

```powershell
docker compose down -v
```

Or:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\stop-backend.ps1 -RemoveVolume
```

## Helper scripts

These scripts are available if you want task-specific shortcuts:

- [scripts/start-backend.ps1](</D:/training_golang/Module 11 (Hour 81-88) Building Advanced CLI Applications/Hour88/distkv-cli/scripts/start-backend.ps1>): start Docker backend
- [scripts/build-cli.ps1](</D:/training_golang/Module 11 (Hour 81-88) Building Advanced CLI Applications/Hour88/distkv-cli/scripts/build-cli.ps1>): build the Go project
- [scripts/run-practical.ps1](</D:/training_golang/Module 11 (Hour 81-88) Building Advanced CLI Applications/Hour88/distkv-cli/scripts/run-practical.ps1>): run the main CLI workflow
- [scripts/watch-demo.ps1](</D:/training_golang/Module 11 (Hour 81-88) Building Advanced CLI Applications/Hour88/distkv-cli/scripts/watch-demo.ps1>): run the watcher
- [scripts/stop-backend.ps1](</D:/training_golang/Module 11 (Hour 81-88) Building Advanced CLI Applications/Hour88/distkv-cli/scripts/stop-backend.ps1>): stop Docker backend
- [scripts/example-bucket-create.ps1](</D:/training_golang/Module 11 (Hour 81-88) Building Advanced CLI Applications/Hour88/distkv-cli/scripts/example-bucket-create.ps1>): example create command
- [scripts/example-put-v1.ps1](</D:/training_golang/Module 11 (Hour 81-88) Building Advanced CLI Applications/Hour88/distkv-cli/scripts/example-put-v1.ps1>): example put command
- [scripts/example-get-version.ps1](</D:/training_golang/Module 11 (Hour 81-88) Building Advanced CLI Applications/Hour88/distkv-cli/scripts/example-get-version.ps1>): example get command
- [scripts/example-list.ps1](</D:/training_golang/Module 11 (Hour 81-88) Building Advanced CLI Applications/Hour88/distkv-cli/scripts/example-list.ps1>): example list command
- [scripts/example-watch-app.ps1](</D:/training_golang/Module 11 (Hour 81-88) Building Advanced CLI Applications/Hour88/distkv-cli/scripts/example-watch-app.ps1>): example watch command

## Common commands reference

Create bucket:

```powershell
go run . bucket create
```

Get bucket info:

```powershell
go run . bucket info
```

Put a value:

```powershell
go run . put my.key hello
```

Get a value:

```powershell
go run . get my.key
```

Delete a key:

```powershell
go run . delete my.key
```

List keys:

```powershell
go run . list
```

Watch changes:

```powershell
go run . watch ">"
```

Health check:

```powershell
go run . health
```

## Command effects summary

This section explains what each main command changes in the system.

`go run . bucket create`

- creates the JetStream KV bucket if missing
- reuses the existing bucket if it already exists
- prepares the backend storage used by the CLI

`go run . put app.config.version v1`

- writes the key `app.config.version`
- stores the value `v1`
- creates revision `1`

`go run . get app.config.version`

- reads the current value from the server
- does not modify any data

`go run . list`

- shows all keys in the current bucket
- does not modify any data

`go run . watch "app.>"`

- opens a live subscription for matching keys
- prints updates as they happen
- does not modify any data

`go run . health`

- checks whether the CLI can reach the server
- confirms JetStream is enabled
- does not modify any data

## How to check on the server side

The CLI talks to this NATS server by default:

- NATS server URL: `nats://127.0.0.1:4222`
- Monitoring URL: `http://127.0.0.1:8222`

### Check container status

```powershell
docker compose ps
```

### Check server health endpoint

```powershell
Invoke-WebRequest http://127.0.0.1:8222/healthz
```

### Check JetStream information

```powershell
Invoke-WebRequest "http://127.0.0.1:8222/jsz?streams=true" | Select-Object -ExpandProperty Content
```

This shows JetStream information and stream details. For the bucket `distkv`, look for the backing stream `KV_distkv`.

### Open NATS monitoring in the browser

```powershell
start http://127.0.0.1:8222
```

### Check whether the CLI can reach the server

```powershell
go run . health
```

Expected output:

```text
server=nats://127.0.0.1:4222 jetstream=true bucket=distkv
```

### Check stored data from the CLI side

```powershell
go run . list
go run . get app.config.version
```

### Check live updates

Console 1:

```powershell
go run . watch "app.>"
```

Console 2:

```powershell
go run . put app.config.version v4
```

## Important note about the monitoring API

If you try this:

```powershell
Invoke-WebRequest http://127.0.0.1:8222/var
```

and get:

```text
404 page not found
```

that is expected in this setup.

Use these endpoints instead:

- `http://127.0.0.1:8222/healthz`
- `http://127.0.0.1:8222/jsz`
- `http://127.0.0.1:8222/jsz?streams=true`
- `http://127.0.0.1:8222/connz`
- `http://127.0.0.1:8222/routez`
- `http://127.0.0.1:8222/subsz`

Examples:

```powershell
Invoke-WebRequest http://127.0.0.1:8222/healthz
Invoke-WebRequest "http://127.0.0.1:8222/jsz?streams=true" | Select-Object -ExpandProperty Content
```

So for this project, do not use `/var` as the primary verification endpoint. Use `/healthz` and `/jsz?streams=true` instead.

## Useful flags with full examples

The CLI supports these global flags:

- `--server`: NATS server URL
- `--bucket`: JetStream KV bucket name
- `--timeout`: request timeout
- `--json`: structured machine-readable output

### `--server`

Use this when your NATS server is running on a different address or port.

Example:

```powershell
go run . --server nats://127.0.0.1:4222 health
```

What it means:

- tells the CLI which NATS server to connect to
- useful when the server is not using the default address

Another example with a different host:

```powershell
go run . --server nats://192.168.1.10:4222 health
```

### `--bucket`

Use this when you want to store data in a different JetStream KV bucket instead of the default `distkv`.

Example:

```powershell
go run . --bucket training bucket create
```

Now write and read data from that bucket:

```powershell
go run . --bucket training put user.name subhash
go run . --bucket training get user.name
```

What it means:

- creates or uses the bucket named `training`
- keeps its data separate from the default `distkv` bucket

### `--timeout`

Use this when you want the command to wait longer or fail faster.

Example:

```powershell
go run . --timeout 10s health
```

Another example:

```powershell
go run . --timeout 2s get app.config.version
```

What it means:

- `10s` means wait up to 10 seconds
- `2s` means wait up to 2 seconds
- useful when network connections are slow or when you want faster failure

### `--json`

Use this when you want machine-readable output instead of plain text.

Example:

```powershell
go run . --json health
```

Example output:

```json
{
  "server": "nats://127.0.0.1:4222",
  "jetstream": true,
  "bucket": "distkv"
}
```

Another example:

```powershell
go run . --json get app.config.version
```

Example output:

```json
{
  "bucket": "distkv",
  "key": "app.config.version",
  "value": "v2",
  "revision": 2,
  "created": "2026-04-27T00:00:00Z"
}
```

### Using multiple flags together

You can combine flags in one command.

Example:

```powershell
go run . --server nats://127.0.0.1:4222 --bucket training --timeout 10s --json health
```

What this command does:

- connects to the server at `127.0.0.1:4222`
- uses the bucket `training`
- waits up to 10 seconds
- prints JSON output

## Default settings

By default the CLI uses:

- server: `nats://127.0.0.1:4222`
- bucket: `distkv`
- timeout: `5s`

You can override them:

```powershell
go run . --server nats://127.0.0.1:4222 --bucket training health
```

## If something does not work

Check whether the backend is running:

```powershell
docker compose ps
```

Check whether the CLI can reach the backend:

```powershell
go run . health
```

Rebuild the project:

```powershell
go build ./...
```

Restart the backend:

```powershell
docker compose down
docker compose up -d
```
