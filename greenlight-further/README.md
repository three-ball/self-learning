# Greenlight - Let's Go Further

> This project is based on "Let's Go Further - Alex Edwards" ebook. Just a simple HTTP Server for serving film database.

## Project structure

```bash
.
├── bin
├── cmd
│   └── api
│       └── main.go
├── go.mod
├── internal
├── Makefile
├── migrations
├── README.md
└── remote

7 directories, 4 files
```

- The `bin` directory will contain our compiled application binaries, ready for deployment to a production server.
- The `cmd/api` directory will contain the application-specific code for our Greenlight API application. This will include the code for running the server, reading and writing HTTP requests, and managing authentication.
- The `internal` directory will contain various ancillary packages used by our API. It will contain the code for interacting with our database, doing data validation, sending emails and so on. Basically, any code which isn’t application-specific and can potentially be reused will live in here. Our Go code under cmd/api will import the packages in the internal directory (but never the other way around).
- The `migrations` directory will contain the SQL migration files for our database.
- The `remote` directory will contain the configuration files and setup scripts for our production server.
- The `Makefile` will contain recipes for automating common administrative tasks — like auditing our Go code, building binaries, and executing database migrations

## Notes

### Configuring the Database Connection Pool

> How does the `sql.DB` connection pool work?

- Pool conatins two types of connections: `in-use` connections (the connection when you're using it to perform database task) and `idle` connections.
- When you instruct Go to perform a database task, it will first check if any idle connections are available in the pool. If one is available, then Go will reuse this existing connection and mark it as in-use for the duration of the task. If there are no idle connections in the pool when you need one, then Go will create a new additional connection.

#### SetMaxOpenConns()

- Max connection limit (`in-use` + `idle`). Default is unlimited.
- The higher SetMaxOpenConss, *the more database queries can be performed in concurrently.*
- By default, PGSQL has hard limit of 100 open connections (can be modify in `postgresql.conf`), if the connection application create more than PGSQL allow, it will return the error ***"sorry, too many clients already"***.
- If the MaxOpenConns limit is reached, and all connections are in-use, then any further database tasks will be forced to wait until a connection becomes free and marked as idle.

#### SetMaxIdleConns()

- Max idle connections limit in the pool. Default is 2.
- In theory, allowing a higher number of idle connections in the pool will improve performance because it makes ***it less likely that a new connection needs to be established from scratch***.
- `MaxIdleConns` limit should always be less than or equal to `MaxOpenConns`.

#### SetConnMaxLifetime()

- The maximum length of time the a connection can be reused for. Default is no limit life time.
- If we set ConnMaxLifetime to one hour, for example, it means that all connections will be marked as ‘expired’ one hour after they were first created, and cannot be reused after they’ve expired.
  - This doesn’t guarantee that a connection will exist in the pool for a whole hour; it’s possible that a connection will become unusable for some reason and be automatically closed before then.
  - A connection can still be in use more than one hour after being created — it just cannot start to be reused after that time.
  - **This isn’t an idle timeout**. The connection will expire one hour after it was first created — not one hour after it last became idle.
  - Once every second Go runs a background cleanup operation to remove expired connections from the pool.

#### SetConnMaxIdleTime()

- This works in a very similar way to ConnMaxLifetime , except it sets the maximum length of time that a connection can be idle for before it is marked as expired. By default there’s no limit.

