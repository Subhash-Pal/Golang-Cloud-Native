

Hour 8 - 
Mock Test: Build a Rate-Limited Server
In this mock test, you will build a rate-limited HTTP server in Go. The server will limit the number of requests a client can make within a specific time window (e.g., 10 requests per minute). This exercise will also include theoretical concepts and debugging techniques.

Part 1: Theory

What is Rate Limiting?
Rate limiting is a technique used to control the number of requests a client can make to a server within a given time period. It helps:
Prevent abuse or overuse of resources.
Protect the server from denial-of-service (DoS) attacks.
Ensure fair usage for all clients.
Key Concepts
Token Bucket Algorithm:
A client has a "bucket" of tokens, where each token allows one request.
Tokens are replenished at a fixed rate (e.g., 1 token per second).
If the bucket is empty, further requests are denied until tokens are replenished.
Leaky Bucket Algorithm:
Requests are queued and processed at a fixed rate.
Excess requests overflow and are rejected.
Fixed Window vs. Sliding Window:
Fixed Window: Limits requests within discrete time intervals (e.g., per minute).
Sliding Window: Smoothly limits requests across overlapping time windows for better accuracy.
Headers for Rate Limiting:
X-RateLimit-Limit: Total allowed requests in the current window.
X-RateLimit-Remaining: Remaining requests in the current window.
X-RateLimit-Reset: Time when the rate limit resets (in seconds).
Middleware-Based Rate Limiting:
Middleware intercepts requests, checks the rate limit, and either allows or denies the request.
