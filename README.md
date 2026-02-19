# Staff Service

This Go implementation is the minimum viable service for a company to manage their staff records. It currently uses a CSV file (`staff.csv`) as its data source. To better simulate a production database, strict validation and consistent error handling are enforced. The service is built using Gin, a lightweight Go framework designed for REST APIs and microservices. It was implemented with the goal in mind for this to be easily maintained, updated and expanded while not losing focus on scalability and resiliency.

---

## Description
The `Staff Service` loads staff records (`ID`, `Name`, `Role`, and `Quote`) from a data source and stores them in the repository allowing for the specific records to be retrieved using its `ID`.

---

## Features

- Simple health check endpoint (`GET /health`)
- Retrieve staff members by ID (`GET /staff/:id`)
- Input validation for staff IDs
- Clear error handling and HTTP status codes
- Well-formed JSON responses
- Clean architecture separating handlers, repository, and models

---

## How to run
From project root run:
`go run cmd/app/main.go`
To confirm the service is running, access it using a browser or curl:
`http://localhost:8080/health` or run ```curl http://localhost:8080/health```
Upon receiving an OK status, access the desired endpoint using a browser or curl:
`http://localhost:8080/staff/{id}`or run ```curl http://localhost:8080/staff/{id}```
Replace `id` with an valid ID from `staff.csv`

---

## Technology Stack

| Layer | Technology |
|-------|------------|
| Language | Go (go 1.25.7) |
| HTTP Server | Gin |
| Logging | Zap |
| Data | CSV (temporary) |
| Architecture | Handlers, Interfaces, Repositories |

---

## Framework choice

### Why use Gin?

Gin was chosen to enable rapid development and provide clear, well-established patterns for routing, middleware, error handling, and JSON responses. The service is intentionally structured so that core business logic, data access, and validation are not coupled to Gin. Handlers depend on interfaces rather than concrete implementations, making it straightforward to replace Gin or migrate to a different HTTP stack in the future if greater control or flexibility is desired. Although Gin provides a `Run()` helper, this service uses an explicit `http.Server` setup to allow for graceful shutdown. This is not necessary for the current in-memory repository, but ensures the service can evolve cleanly to a database-backed implementation without reworking the application lifecycle.

## Design Considerations

### Why a health check endpoint?
While the only endpoint explicitly requested was `GET /staff/:id`, a health check is included.... A health check is important to be included in any production service so that there is a way to see if the service is up and running or down. The health check is in its own package due to the potential additions that will be added to this service in the future.

Although the only requested endpoint was GET /staff/:id, the goal was to see a production quality service which should have a health check. It is included to verify service availability and has been implemented in its own package to accommodate future additions or the addition of metrics.

### ID Validation
- The original data source contained only 8-character lowercase hexadecimal IDs. This service's validation was designed to be more flexible, allowing any alphanumeric ID up to 64 characters. This accommodates future changes in ID format (e.g., migrating to UUIDs) without modifying validation logic; only existing records would need updates.

- Each ID must be unique. While this constraint will ultimately be enforced at the database level, it is currently checked and enforced in-memory while using the current data source, CSV.

- Special characters are disallowed to reduce the risk of SQL injection and other malicious inputs, anticipating the future database  migration and its SQL queries for this endpoint or those to come.

### Dependency Injection
Dependency injection enables:
- Easy mocking and testing
- Seamless migration from CSV to a database-backed repository

The StaffRepository interface allows StaffHandler to accept any repository implementing GetStaffMemberByID(). To switch repositories, only the constructor argument needs updating. All unit tests are placed in separate packages to further enforce clean testing boundaries and realistically simulate client usage.

---

## Future Work

- Replace CSV storage with a database-backed repository
    - When database-backed repository switch is being implemented `context.Context` should be added to the handler and repository interface
- Add authentication and authorization
- Improve observability (metrics, tracing)
- Run the service behind a reverse proxy (i.e. Nginx) instead of exposing it directly

---

## Endpoints
### Request/Response

- GET /health
```
{
  "status": "OK"
}
```

- GET /staff/:id

`abc123` is an ID and in the CSV file:
GET /staff/abc123
Response:
```
{
  "id": "abc123",
  "name": "Jane Doe",
  "role": "Engineer",
  "quote": "Keep it simple."
}
```

`123abc` is a valid ID but not in the CSV file
GET /staff/123abc
Response:
```
{
  "error": "staff member not found"
}
```

`987zyx!` is not a valid ID
GET /staff/987zyx!
Response:
```
{
  "error": "id contains invalid characters"
}
```

LLM conversation:
https://chatgpt.com/share/698b6e4d-74f8-800e-858c-a54a0abad5ba