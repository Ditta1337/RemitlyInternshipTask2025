## Remitly Internship Task 2025

## Prerequisites

Make sure you have Docker installed:

- [Docker](https://docs.docker.com/get-docker/)

## How to run

1. Clone the Repository

    ```sh
    git clone https://github.com/Ditta1337/RemitlyInternshipTask2025.git
    cd RemitlyInternshipTask2025
    ```

2. Run the application

    ```shell
    docker-compose up --build
    ```

   This will:

    - Start a PostgreSQL database
    - Build the Go application inside a container
    - Run database migrations
    - Seed database with data (`internal/db/seed/SWIFT_CODES.tsv`)

## How to run tests

1. Run the following command in the project root directory:

    ```shell
    make test
    ```

## Available endpoints

#### SWIFT Codes

- `POST /v1/swift-codes`
    - Adds new bank to the database
    - Example payload:
    ```json
    {
        "swiftCode": "FAKECODEXXX",
        "address": "Some Address",
        "bankName": "Headquarter bank US",
        "countryISO2": "US",
        "countryName": "United States",
        "isHeadquarter": true
    }
    ```

- `GET /v1/swift-codes/{swift-code}`
    - Retrieves details of a bank by its SWIFT code
    - Returns one of two structures:
        - when retrieved bank is a headquarter:
        ```json
        {
            "address": "string",
            "bankName": "string",
            "countryISO2": "string",
            "countryName": "string",
            "isHeadquarter": true,
            "swiftCode": "string",
            "branches": [
                {
                    "address": "string",
                    "bankName": "string",
                    "countryISO2": "string",
                    "isHeadquarter": true,
                    "swiftCode": "string"
                },
                {
                    "address": "string",
                    "bankName": "string",
                    "countryISO2": "string",
                    "isHeadquarter": true,
                    "swiftCode": "string"
                }, ...
            ]
        }
        ```
        - when retrieved bank is a branch:
        ```json
        {
            "address": "string",
            "bankName": "string",
            "countryISO2": "string",
            "countryName": "string",
            "isHeadquarter": false,
            "swiftCode": "string"
        }
        ```

- `GET /v1/swift-codes/country/{countryISO2code}`
    - Retrieves a list of banks for a given ISO2 country code
    - Returns this structure:
    ```json
    {
        "countryISO2": "string",
        "countryName": "string",
        "swiftCodes": [
            {
                "address": "string",
                "bankName": "string",
                "countryISO2": "string",
                "isHeadquarter": true,
                "swiftCode": "string"
            },
            {
                "address": "string",
                "bankName": "string",
                "countryISO2": "string",
                "isHeadquarter": true,
                "swiftCode": "string"
            }, ...
        ]
    }
    ```

- `DELETE /v1/swift-codes/{swift-code}`
    - Removes bank from the database

#### (ADDITIONAL) Swagger
- `GET /v1/swagger/*`
    - Swagger documentation for the API

## Environment variables
The app uses the following environment variables (defined in docker-compose.yaml):

| Variable          | Default Value                                           | Description                                     |
|------------------|---------------------------------------------------------|-------------------------------------------------|
| `ADDR`          | `:8080`                                                 | API listen port                                 |
| `EXTERNAL_URL`  | `http://localhost:8080`                                 | Public-facing API URL                           |
| `DB_ADDR`       | `postgres://admin:remitly2025@db/swift?sslmode=disable` | Database connection string                      |
| `DB_MAX_OPEN_CONNS` | `30`                                                    | Max open DB connections                         |
| `DB_MAX_IDLE_CONNS` | `30`                                                    | Max idle DB connections                         |
| `DB_MAX_IDLE_TIME`  | `15m`                                                   | Max idle time for DB connections                |
| `ENV`          | `production`                                            | App environment (`development` or `production`) |
| `API_VERSION`  | `v1`                                                    | API version                                     |