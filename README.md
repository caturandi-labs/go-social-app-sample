# Golang Social App Backend

A simple monolithic social blog API application built with Golang.

## Tech Stack Used

1. Go Chi v5 (Router)
2. Pq lib
3. Swaggo (api docs)
4. Docker for containerization
5. etc.

## Project Structure

This project uses a simple layered architecture. The details are as follows:

- Application layer / Delivery-related code is located in the `cmd/api` directory.
- Domain and Infrastructure layers are located inside the `internal` directory.

&copy; by caturandi-labs 2025 - MIT License