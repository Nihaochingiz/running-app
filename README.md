## Compose sample application
### Go server with a Nginx proxy and a Postgres database

Project structure:
```
.
├── backend
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── db
│   └── password.txt
├── compose.yaml
├── proxy
│   └── nginx.conf
└── README.md
```

[_compose.yaml_](compose.yaml)
```shell
services:
  backend:
    build:
      context: backend
      target: builder
    ...
  db:
    image: postgres
    ...
  proxy:
    image: nginx
    volumes:
      - type: bind
        source: ./proxy/nginx.conf
        target: /etc/nginx/conf.d/default.conf
        read_only: true
    ports:
      - 80:80
    ...
```
The compose file defines an application with three services `proxy`, `backend` and `db`.
When deploying the application, docker compose maps port 80 of the proxy service container to port 80 of the host as specified in the file.
Make sure port 80 on the host is not already being in use.

## Deploy with docker compose

```shell
$ docker compose up -d
Creating network "nginx-golang-postgres_default" with the default driver
Pulling db (postgres:)...
latest: Pulling from library/postgres
...
Successfully built 5f7c899f9b49
Successfully tagged nginx-golang-postgres_proxy:latest
WARNING: Image for service proxy was built because it did not already exist. To rebuild this image you must use `docker-compose build` or `docker-compose up --build`.
Creating nginx-golang-postgres_db_1 ... done
Creating nginx-golang-postgres_backend_1 ... done
Creating nginx-golang-postgres_proxy_1   ... done
```

## Expected result

Listing containers must show three containers running and the port mapping as below:
```shell
$ docker compose ps
NAME                              COMMAND                  SERVICE             STATUS              PORTS
nginx-golang-postgres-backend-1   "/code/bin/backend"      backend             running
nginx-golang-postgres-db-1        "docker-entrypoint.s…"   db                  running (healthy)   5432/tcp
nginx-golang-postgres-proxy-1     "/docker-entrypoint.…"   proxy               running             0.0.0.0:80->80/tcp
```

After the application starts, navigate to `http://localhost:80` in your web browser or run:
```shell
$ curl http://localhost:8000
```

Stop and remove the containers
```shell
$ docker compose down
```
API Endpoints
1. Get All Running Statistics
Endpoint: GET /running-statistics
Response: Returns a JSON array of all running statistics.
Example:
json
{
  "statistics": [
    {
      "id": 1,
      "date": "2021-12-01",
      "distance": "5km",
      "time": "00:30:00",
      "created_at": "2021-12-01T12:00:00Z"
    }
  ]
}
2. Create Running Statistics
Endpoint: POST /running-statistics
Request Body: JSON object containing the running statistic data.
Example:
json
{
  "date": "2021-12-01",
  "distance": "5km",
  "time": "00:30:00"
}
Response: Returns the ID of the created statistic.
Response Example:
json
{
  "id": 2
}
3. Get a Specific Running Statistic
Endpoint: GET /running-statistics/{id}
Response: Returns the running statistic with the specified ID.
Example:
json
{
  "id": 1,
  "date": "2021-12-01",
  "distance": "5km",
  "time": "00:30:00",
  "created_at": "2021-12-01T12:00:00Z"
}
4. Delete a Running Statistic
Endpoint: DELETE /running-statistics/{id}
Response: Returns a 204 No Content status if the deletion is successful or a 404 Not Found status if no statistic with the given ID exists.
Database Structure
The application uses a PostgreSQL database with a table named running_statistics. The table structure is as follows:

sql
CREATE TABLE running_statistics (
    id SERIAL PRIMARY KEY,
    date VARCHAR(10),
    distance VARCHAR(10),
    time VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
Error Handling
The API returns appropriate HTTP status codes for different error scenarios:

400 Bad Request for client-side errors (e.g., invalid input).
404 Not Found when the requested resource does not exist.
500 Internal Server Error for server-side issues.
Logging
The application logs incoming HTTP requests to the standard output using the Gorilla Handlers package.

Contributing
If you would like to contribute to this project, please fork the repository, create a new branch, and submit a pull request.

License
This project is licensed under the MIT License. See the LICENSE file for more details.
