# Project Backlog

## Setup and Configuration
- [x] **Initialize Go Module**: Set up the Go project structure with modules for better dependency management.
- [x] **Define load balancer**: Choose between ngnix or traefik as load balancer.
  - Outcome: I Decided to use nginx 
- [x] **Database Selection and Integration**: Choose a performant database, integrate with Go using direct SQL, considering consistency and performance for operations.
  - postgresql will be the database.


## Transactions Endpoint
- [x] **POST /clientes/[id]/transacoes Endpoint Implementation**: Implement the endpoint to handle credit and debit transactions, ensuring all fields are validated as per the requirements.
  - Validate `id` as an integer representing the customer ID.
  - Ensure `valor` is a positive integer for transaction amount in cents.
  - Restrict `tipo` to "c" for credit or "d" for debit.
  - Limit `descricao` length to 1-10 characters.
  - Implement business logic to prevent debit transactions from exceeding the customer's available limit.
  - Return appropriate HTTP status codes (200 for success, 422 for invalid debit operation, and 404 for non-existent customer ID).

- [x] **Transaction Business Logic**: Develop the logic to update the customer's balance and limit upon a transaction, and enforce rules for debit transactions.

- [ ] **Unit and Integration Tests for Transactions**: Ensure all edge cases and success scenarios are covered.

## Statement Endpoint
- [x] **GET /clientes/[id]/extrato Endpoint Implementation**: Create the endpoint to fetch a customer's statement, including balance, limit, and the last transactions.
  - Implement logic to format and return the statement data correctly.
  - Ensure the list of transactions is returned in descending order by date.
  - Return HTTP 404 for non-existent customer IDs.

- [ ] **Unit and Integration Tests for Statement**: Test the endpoint for accuracy, performance, and error handling.

## Initial Customer Setup
- [x] **Database Seed Script**: Write a script to pre-populate the database with the specified customer records, ensuring correct IDs, limits, and initial balances are set.
  - Explicitly exclude ID 6 to test for the non-existence case.

- [ ] **Verification Tests for Initial Setup**: Confirm that the seeding process correctly initializes the database with the required customer data.

Remember to prioritize these tasks based on dependencies, starting with the initial setup, followed by core endpoint implementation, and lastly, testing and optimization.

