# go-tasks-api

### Development & Usage

This project is fully containerized and uses docker compose for local development and testing.

#### Build

Build all required Docker images (application and migrations):
```
make build
```

#### Run the application


```
make boot
```
This will build the image if needed and start the API using Docker Compose.


#### Run database migrations

Apply database migrations using the migration container:

```
make run-migration
```

#### Tasks schema

| Field       | Type      | Constraints                         | Description                                 |
| ----------- | --------- | ----------------------------------- | ------------------------------------------- |
| id          | UUID      | PRIMARY KEY                         | Unique identifier of the task               |
| title       | TEXT      | NOT NULL                            | Task title                                  |
| description | TEXT      | NOT NULL, DEFAULT ''                | Task description                            |
| status      | TEXT      | NOT NULL, DEFAULT 'todo'            | Task status (todo, done) |
| created_at  | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp                          |
| updated_at  | TIMESTAMP | DEFAULT NULL                        | Last update timestamp                       |
| is_active   | BOOLEAN   | NOT NULL, DEFAULT TRUE              | Task active flag                            |


#### Testing

Run linting and unit tests:
```
make test
```

Tests are executed inside Docker to ensure a consistent and reproducible environment.

### API Endpoints

| Method | Endpoint             | Description       |
| -----: | -------------------- | ----------------- |
|   POST | `/api/v1/tasks`      | Create a new task |
|    GET | `/api/v1/tasks`      | List all tasks    |
|    GET | `/api/v1/tasks/{id}` | Get task by ID    |
|    PUT | `/api/v1/tasks/{id}` | Update task by ID |
| DELETE | `/api/v1/tasks/{id}` | Delete task by ID |
