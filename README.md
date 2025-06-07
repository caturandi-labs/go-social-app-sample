# Golang Social App Backend

A Simple Monolithic social blog API application built with Golang 

## Tech Stack Used

1. Go Chi v5 (Router)
2. Pq lib
3. Swaggo (api docs)
4. Docker for containerization
5. etc.

## Project Structure

I'm using simple layered architecture in this project with details in below:

- Infrastructure code / Delivery code related is located in the `cmd/api`.
- Internal code such as DB connection, env configuration, data store layer is located inside `internal` directory.

&copy;  caturandi-labs 2025 - MIT License