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

### Handling Partial Updates

- Difference between:
  - A client providing a key/value pair which has a zero-value value — like `{"title": ""}` — in which case we want to return a validation error.
  - A client not providing a key/value pair in their JSON at all — in which case we want to ‘skip’ updating the field but not send a validation error.
- Using `pointer` for partial update.

```Go
	type input struct {
		Title   *string       `json:"title"`   // This will be nil if there is no corresponding key in the JSON.
		Year    *int32        `json:"year"`    // Likewise...
		Runtime *data.Runtime `json:"runtime"` // Likewise...
		Genres  []string      `json:"genres"`  // We don't need to change this because slices already have the zero-value nil.
	}
```
### Optimistic Concurrency Control

- The simplest way is using `version` number of database records.
- 2 user `GET` to get the record with version `N`. Update is only executed if the version number in database is still `N`.
- The version can be integer or uuid (using `uuid_generate_v4()` in PostgreSQL).

```SQL
UPDATE movies
SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
WHERE id = $5 AND version = $6
RETURNING version

-- Or it can be:
UPDATE movies
SET title = $1, year = $2, runtime = $3, genres = $4, version = uuid_generate_v4()
WHERE id = $5 AND version = $6
RETURNING version
```
### Managing SQL Query Timeouts

- Using `context` in Golang such as: `ExecContext()` or `QueryRowContext()`, then return `500 Internal Server Error`.
- We need to do:

```Go
// Use the context.WithTimeout() function to create a context.Context which carries a
// 3-second timeout deadline. Note that we're using the empty context.Background()
// as the 'parent' context.
  ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres), //adatper pq.Array
		&movie.Version,
	)
```
- We could have created a context with a timeout in *our handlers* using the request context as the *parent-context*.
- It will be look like:

```Go

func (app *application) exampleHandler(w http.ResponseWriter, r *http.Request) {

	// Create a context.Context with a one-second timeout deadline and which has the
	// request context as the 'parent'.
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	// Pass the context on to the Get() method.
	example, err := app.models.Example.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

}
```

- It will be consider as the good practice, but it will add a lot of additional complexity. On the other hand, we will receive the same error message in two different scenarios:
  - When the database is "real" timeout.
  - Whene client close `http` connection.
- The solution is: Using `ctx.Err()`.
- First, wrapping the error in DB layer:

```Go
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres), //adatper pq.Array
		&movie.Version,
	)

	if err != nil {
		switch {
		case err.Error() == "pq: canceling statement due to user request":
			// Wrap the error with ctx.Err().
			return nil, fmt.Errorf("%v: %w", err, ctx.Err())
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
```
- Then, in our handler, we should use:

```Go
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	example, err := app.models.Example.Get(ctx, id)
	if err != nil {
		switch {
		// If the error is equal to or wraps context.Canceled, then return without taking
		// any further action.
		case errors.Is(err, context.Canceled):
			return
		case errors.Is(err, data.ErrNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
```

- ***Tradeoff:*** using the request context as the parent context for database timeouts **adds quite a lot of behavioral complexity** and introduces nuances that **you and anyone else working on the codebase needs to be aware of**.
- For most applications, on most endpoints, it’s probably not. The exceptions are probably applications which frequently run close to saturation point of their resources, or for specific endpoints which execute slow running or very computationally expensive SQL queries. In those cases, canceling queries aggressively when a client disappears may have a meaningful positive impact and make it worth the trade-off.