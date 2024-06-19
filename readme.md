

# Simple Banking Service Example
This project leverages the Golang API Scafold by [Chris Fryer](https://github.com/cfryerdev/golang-api-scaffold) and design patterns from [Kenny McClive's](https://github.com/kmcclive/go-api-pattern) example project

## Technologies
Here is a list of the technologies used in this project:
* We use GIN for http routing.
* We use GORM for our ORM layer.
* We use PG for postgres database access.

## Endpoints
Here is an overall layout of what endpoints come with this architecture:

| Method | Route                      | Description                                  |
| ------ | ----------------------     | -------------------------------------------- |
| N/A    | /swagger/index.html        | Swagger UI                                   |
| GET    | /health/                   | Health check.                                |
| POST   | /auth/login                | Creates a JWT token for access.              |
| GET    | /accounts                  | Gets a list of records.                      |
| GET    | /accounts/:id              | Gets a record by id.                         |
| POST   | /accounts/                 | Creates a record.                            |
| PUT    | /accounts/:id              | Updates a record.                            |
| DELETE | /accounts/:id              | Deletes a record.                            |
| GET    | /transactions              | Gets a list of records.                      |
| GET    | /transactions/:id          | Gets a record by id.                         |
| POST   | /transactions/             | Creates a record.                            |
| PUT    | /transactions/:id          | Updates a record.                            |
| DELETE | /transactions/:id          | Deletes a record.                            |


## Installing dependencies
```bash
go mod vendor
```
## Switching between application types
* open the main.go file and change the runtime import line to include the path to the package with the execute function that you want to run 

## Running Api locally
```bash
go run main.go
```
Visit: `http://localhost:8080/swagger/index.html`

## Running CLI (default) 
* follow the prompts and instructions from the console.  

## Using makefile
to see the available targets and their uses, run the following command in the root directory of the project:
``` bash
make help
```
you can start and stop the postgres container by running make on the startdb and stopdb targets

