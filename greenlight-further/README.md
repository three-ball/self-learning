# Greenlight - Let's Go Further

> This project is based on "Let's Go Further - Alex Edwards" ebook. Just a simple HTTP Server for serving film database.

- [Greenlight - Let's Go Further](#greenlight---lets-go-further)
	- [Project structure](#project-structure)
	- [Configuring the Database Connection Pool](#configuring-the-database-connection-pool)
		- [SetMaxOpenConns()](#setmaxopenconns)
		- [SetMaxIdleConns()](#setmaxidleconns)
		- [SetConnMaxLifetime()](#setconnmaxlifetime)
		- [SetConnMaxIdleTime()](#setconnmaxidletime)
	- [Handling Partial Updates](#handling-partial-updates)
	- [Optimistic Concurrency Control](#optimistic-concurrency-control)
	- [Managing SQL Query Timeouts](#managing-sql-query-timeouts)
	- [Query with filter, full-text search and sorting, paging.](#query-with-filter-full-text-search-and-sorting-paging)
		- [Dynamic Filtering](#dynamic-filtering)
		- [Full-Text Search](#full-text-search)
		- [Sorting lists](#sorting-lists)
		- [Paginating Lists](#paginating-lists)
		- [Returning Pagination Metadata](#returning-pagination-metadata)
	- [Structured JSON Log Entries](#structured-json-log-entries)
	- [Rate Limiting](#rate-limiting)
		- [Global Rate Limiting](#global-rate-limiting)
	- [Graceful Shutdown](#graceful-shutdown)
	- [User Model Setup and Registration](#user-model-setup-and-registration)
	- [Background job](#background-job)
	- [Authentication](#authentication)
		- [Authentication options](#authentication-options)
			- [HTTP Basic Authentication](#http-basic-authentication)
			- [Token authentication](#token-authentication)
				- [Stateful token authentication](#stateful-token-authentication)
				- [Stateless token authentication](#stateless-token-authentication)
			- [API Key authentication](#api-key-authentication)
			- [OAuth 2.0 / OpenID Connect](#oauth-20--openid-connect)
			- [What authentication approach should I use?](#what-authentication-approach-should-i-use)
		- [Authenticating Requests](#authenticating-requests)
			- [Anonymous User](#anonymous-user)
			- [Reading and writing to the request context](#reading-and-writing-to-the-request-context)


## Project structure

```bash
.
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
- The `internal` directory will contain various ancillary packages used by our API. It will contain the code for interacting with our database, doing data validation, sending emails and so on. Basically, any code which isn't application-specific and can potentially be reused will live in here. Our Go code under cmd/api will import the packages in the internal directory (but never the other way around).
- The `migrations` directory will contain the SQL migration files for our database.
- The `remote` directory will contain the configuration files and setup scripts for our production server.
- The `Makefile` will contain recipes for automating common administrative tasks — like auditing our Go code, building binaries, and executing database migrations

## Configuring the Database Connection Pool

> How does the `sql.DB` connection pool work?

- Pool conatins two types of connections: `in-use` connections (the connection when you're using it to perform database task) and `idle` connections.
- When you instruct Go to perform a database task, it will first check if any idle connections are available in the pool. If one is available, then Go will reuse this existing connection and mark it as in-use for the duration of the task. If there are no idle connections in the pool when you need one, then Go will create a new additional connection.

### SetMaxOpenConns()

- Max connection limit (`in-use` + `idle`). Default is unlimited.
- The higher SetMaxOpenConss, *the more database queries can be performed in concurrently.*
- By default, PGSQL has hard limit of 100 open connections (can be modify in `postgresql.conf`), if the connection application create more than PGSQL allow, it will return the error ***"sorry, too many clients already"***.
- If the MaxOpenConns limit is reached, and all connections are in-use, then any further database tasks will be forced to wait until a connection becomes free and marked as idle.

### SetMaxIdleConns()

- Max idle connections limit in the pool. Default is 2.
- In theory, allowing a higher number of idle connections in the pool will improve performance because it makes ***it less likely that a new connection needs to be established from scratch***.
- `MaxIdleConns` limit should always be less than or equal to `MaxOpenConns`.

### SetConnMaxLifetime()

- The maximum length of time the a connection can be reused for. Default is no limit life time.
- If we set ConnMaxLifetime to one hour, for example, it means that all connections will be marked as 'expired' one hour after they were first created, and cannot be reused after they've expired.
  - This doesn't guarantee that a connection will exist in the pool for a whole hour; it's possible that a connection will become unusable for some reason and be automatically closed before then.
  - A connection can still be in use more than one hour after being created — it just cannot start to be reused after that time.
  - **This isn't an idle timeout**. The connection will expire one hour after it was first created — not one hour after it last became idle.
  - Once every second Go runs a background cleanup operation to remove expired connections from the pool.

### SetConnMaxIdleTime()

- This works in a very similar way to ConnMaxLifetime , except it sets the maximum length of time that a connection can be idle for before it is marked as expired. By default there's no limit.

## Handling Partial Updates

- Difference between:
  - A client providing a key/value pair which has a zero-value value — like `{"title": ""}` — in which case we want to return a validation error.
  - A client not providing a key/value pair in their JSON at all — in which case we want to 'skip' updating the field but not send a validation error.
- Using `pointer` for partial update.

```Go
	type input struct {
		Title   *string       `json:"title"`   // This will be nil if there is no corresponding key in the JSON.
		Year    *int32        `json:"year"`    // Likewise...
		Runtime *data.Runtime `json:"runtime"` // Likewise...
		Genres  []string      `json:"genres"`  // We don't need to change this because slices already have the zero-value nil.
	}
```
## Optimistic Concurrency Control

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
## Managing SQL Query Timeouts

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
- For most applications, on most endpoints, it's probably not. The exceptions are probably applications which frequently run close to saturation point of their resources, or for specific endpoints which execute slow running or very computationally expensive SQL queries. In those cases, canceling queries aggressively when a client disappears may have a meaningful positive impact and make it worth the trade-off.

## Query with filter, full-text search and sorting, paging.

> please return the **first 5 records** where the movie **name includes godfather** and the genresinclude **crime and drama**,sorted by **descending release year**

```bash
/v1/movies?title=godfather&genres=crime,drama&page=1&page_size=5&sort=-year
```

### Dynamic Filtering

The hardest part of building a dynamic filtering feature like this is the SQL query to retrieve the data — we need it to work with no filters, filters on both `title` and `genres`, or a filter on only one of them.

```SQL
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE (LOWER(title) = LOWER($1) OR $1 = '')
AND (genres @> $2 OR $2 = '{}')
ORDER BY id
```

- This SQL query is designed so that each of the filters behaves like it is 'optional'. For example, the condition `(LOWER(title) = LOWER($1) OR $1 = '')` will evaluate as `true` if the placeholder parameter `$1` is a case-insensitive match for the movie title or the placeholder parameter equals `''`.
- So this filter condition will essentially be 'skipped' when movie title being searched for is the empty string "".

### Full-Text Search

- Using Postgresql Full-Text search function.

```SQL
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE  (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
AND (genres @> $2 OR $2 = '{}')
ORDER BY id
```
- `to_tsvector('simple',title)` function takes a movie tittle and splits it into *lexemes* with `simple` configuration - which means that the lexemes are just lowercase versions of the words in the tittle. Example: the title `The Break Fast Club` would be split into the lexemes `breakfast`, `club` and `the`.
- `plainto_tsquery('simple', $1)` function takes a search value and turns it into a formatted query term. It uses the `simple` configuration. The search value `The Club` would result in the query term **`the` & `club`**
- The `@@` is the matches operator.
- Create index for text search: [GIN Index](https://www.postgresql.org/docs/current/textsearch-indexes.html).

```SQL
CREATE INDEX IF NOT EXISTS movies_title_idx ON movies USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS movies_genres_idx ON movies USING GIN (genres);
```
- Also useful document about indexing: [Index Types](https://www.postgresql.org/docs/13/indexes-types.html). 

### Sorting lists

- Likewise, in our database multiple movies will have the same year value. If we order based on the year column, then the movies are guaranteed be ordered by year, but the movies for a particular year could appear in any order at any time.
- We just we need to make sure that the ORDER BY clause always includes a primary key column:

```SQL
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE (STRPOS(LOWER(title), LOWER($1)) > 0 OR $1 = '')
AND (genres @> $2 OR $2 = '{}')
ORDER BY year DESC, id ASC
```
- Dynamic Sorting:

```Go
// internal/model/filters.go
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

// internal/model/movies.go
query := fmt.Sprintf(`
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC`, filters.sortColumn(), filters.sortDirection())
```
### Paginating Lists

```SQL
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
AND (genres @> $2 OR $2 = '{}')
ORDER BY %s %s, id ASC
LIMIT 5 OFFSET 10
```
- In the `offset()` method there is the theoretical risk of an `**integer overflow**` as we are multiplying two int values together However, this is mitigated by the validation rules we created in our **ValidateFilters() function, where we enforced maximum values of page_size=100 and page=10000000** (10 million).

### Returning Pagination Metadata

- The response will be better if we add some additional metadata:

```json
{
    "metadata": {
        "current_page": 1,
        "page_size": 20,
        "first_page": 1,
        "last_page": 42,
        "total_records": 832
    },
    "movies": [
        {
            "id": 1,
            "title": "Moana",
            "year": 2015,
            "runtime": "107 mins",
            "genres": [
                "animation",
                "adventure"
            ],
            "version": 1
        }
    ]
}
```

- The challenging part of doing this is generating the total_records figure. We want this to ***reflect the total number of available records given the title and genres filters that are applied*** — ***not the absolute total of records*** in the movies table.
- Using [Window Function](https://www.postgresql.org/docs/current/tutorial-window.html)
- The inclusion of the `count(*) OVER()` expression at the start of the query will result in the filtered record count being included as the first value in each row.

```SQL
SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
FROM movies
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
AND (genres @> $2 OR $2 = '{}')
ORDER BY %s %s, id ASC
LIMIT $3 OFFSET $4
```
- Define metadata:

```Go
// Define a new Metadata struct for holding the pagination metadata.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// The calculateMetadata() function calculates the appropriate pagination metadata
// values given the total number of records, current page, and page size values. Note
// that the last page value is calculated using the math.Ceil() function, which rounds
// up a float to the nearest integer. So, for example, if there were 12 records in total
// and a page size of 5, the last page value would be math.Ceil(12/5) = 3.
func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		// Note that we return an empty Metadata struct if there are no records.
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
```
## Structured JSON Log Entries

| Key       	| Description       |
|----------------|----------------|
| Level  		| A code indicating the severity of the log entry. In this project we will use the following three severity levels, ordered from least to most severe: <br>  - `INFO` (least severe) <br> - `ERROR` <br> - `FATAL` (most severe)  |
| Time  		| The UTC time that the log entry was made with second precision.  |
| Message  		| A string containing the free-text information or error message.  |
| Properties  	| Any additional information relevant to the log entry in string key/value pairs  |
| Trace  		| A stack trace for debugging purposes (optional). |

## Rate Limiting

### Global Rate Limiting

```bash
go get golang.org/x/time/rate@latest
```

> A Limiter controls how frequently events are allowed to happen. It implements a ***“token bucket”*** ofsize `b`, ***initially full and refilled at rate r tokens persecond***.

- We should **limiter per client** (one client bad action can't effect to another).
- Using Map to store client limiter, but **our clients will grow up, so we should add clean up job**.
- Map will be used in concurrent routines, we should you `mutext` to avoid racing condition.
- If your infrastructure is distributed, with **your application running on multiple servers behind a load balancer, then you'll need to use an alternative approach**.
- We can use function of HAProxy, Nginx or cache database like Redis.

```Go
func (app *application) rateLimit(next http.Handler) http.Handler {
	// maximum is 4 requests
	// restore 2 requests per second

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			// Lock the mutex to prevent any rate limiter checks from happening while
			// the cleanup is taking place.
			mu.Lock()
			// Loop through all clients. If they haven't been seen within the last three
			// minutes, delete the corresponding entry from the map.
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			// Importantly, unlock the mutex when the cleanup is complete.
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Extract the client's IP address from the request.
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			// Lock the mutex to prevent this code from being executed concurrently.
			mu.Lock()
			// Check to see if the IP address already exists in the map. If it doesn't, then
			// initialize a new rate limiter and add the IP address and limiter to the map.
			if _, found := clients[ip]; !found {
				clients[ip].limiter = rate.NewLimiter(2, 4)
			}

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}
			mu.Unlock()
			next.ServeHTTP(w, r)
		})
}
```

## Graceful Shutdown

| Key       	| Description    						|  Keyboard shortcut  | Catchable  |
|---------------|---------------------------------------|---------------------|------------|
| SIGINT  		| Interrupt from keyboard 				| Ctrl+C 			  | Yes 	   |
| SIGQUIT  		| Quit from keyboard 					| Ctrl+\ 			  | Yes 	   |
| SIGKILL  		| Kill process (terminate immediately)  | - 				  | No 	       |
| SIGTERM  		| Terminate process in orderly manner 	| - 				  | Yes 	   |

- Implementation:

```Go
func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})
	shutdownError := make(chan error)
	// Start a background goroutine.
	go func() {
		// Create a quit channel which carries os.Signal values.
		quit := make(chan os.Signal, 1)

		// Use signal.Notify() to listen for incoming SIGINT and SIGTERM signals and
		// relay them to the quit channel. Any other signals will not be caught by
		// signal.Notify() and will retain their default behavior.
		

		// Read the signal from the quit channel. This code will block until a signal is
		// received.
		s := <-quit

		// Log a message to say that the signal has been caught. Notice that we also
		// call the String() method on the signal to get the signal name and include it
		// in the log entry properties.
		app.logger.PrintInfo("caught signal", map[string]string{
			"signal": s.String(),
		})

		// Create a context with a 5-second timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		// Call Shutdown() on our server, passing in the context we just made.
		// Shutdown() will return nil if the graceful shutdown was successful, or an
		// error (which may happen because of a problem closing the listeners, or
		// because the shutdown didn't complete before the 5-second context deadline is
		// hit). We relay this return value to the shutdownError channel.
		shutdownError <- srv.Shutdown(ctx)
	}()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}
	// At this point we know that the graceful shutdown completed successfully and we
	// log a "stopped server" message.
	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})
	return nil
}
```

- At first glance this code might seem a bit complex, but at a high-level what it's doing can be summarized very simply: when we receive a `SIGINT` or `SIGTERM` signal, we instruct our server to stop accepting any new HTTP requests, and give any in-flight requests a **'grace period'** of 5 secondsto complete before the application isterminated.

## User Model Setup and Registration

```sql
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1
);
```

- `citext` (case-insensitive text): store text data exactly as it is inputted - without changing the case in anyway. But the comparisons against the data are alway case-insensitive... Including lookups on associated indexes.
- `UNIQUE` combined with `citext` means: t no two rows in the database can have the same email value — even if they have different cases - ***no two usersshould exist with the same email address***.
- `bytea` (binary string) to store one-way hash of the user's password.

## Background job

Some important things we should consider when using background job:

- Note that we don't want to use the `app.serverErrorResponse()` helper to handle any errors in our background goroutine, as that would result in us trying to write a second HTTP response
- The code running in the background goroutine forms a **closure** over the user and app variables. It's important to be aware that these **'closed over'** variables are not scoped to the background goroutine, which means that any changes you make to them will be reflected in the rest of your codebase.
- Background job should be retried n times to increase the probability that emails are successfully sent.

```Go
	for i := 1; i <= 3; i++ {
		err = m.dialer.DialAndSend(msg)
		// If everything worked, return nil.
		if nil == err {
			return nil
		}
		// If it didn't work, sleep for a short time and retry.
		time.Sleep(500 * time.Millisecond)
	}
```

- It's important to bear in mind that any panic which happens in this background goroutine
will not be automatically recovered by our `recoverPanic()` middleware or Go's http.Server.

```Go
// The background() helper accepts an arbitrary function as a parameter.
func (app *application) background(fn func()) {
	// Launch a background goroutine.
	go func() {
		// Recover any panic.
		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()
		// Execute the arbitrary function that we passed as the parameter.
		fn()
	}()
}

// Using helper to recover panic:
	app.background(func() {
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})
```

## Authentication

### Authentication options

- Some of the most common approaches at a high level:
  - Basic authentication
  - Stateful token authentication
  - Stateless token authentication
  - API key authentication
  - OAuth 2.0 / OpenID Connect

#### HTTP Basic Authentication

- The client includes an Authorization header with every request containing their credentials. The credentials need to be in the format `username:password` and ***base-64 encoded***. It will look like in the header:

```bash
 Authorization: Basic YWxpY2VAZXhhbXBsZS5jb206cGE1NXdvcmQ=
```

- It's often useful in the scenario where your API doesn't have "real" user account but you want a quick and easy way to restrict access.
- For APIs with "real" user accounts and — in particular — hashed passwords, it's not such a great fit. **Comparing the password provided by a client against a (slow) hashed password is a deliberately costly operation, and when using HTTP basic authentication you need to do that check for every request**.

#### Token authentication

- The client sends a request to your API containing their credentials. The API verifies that the credentials are correct, generates a bearer token which represents the user, and sends it back to the user. It will look like in the header:

```bash
 Authorization: Bearer <token>
```

- For APIs where user passwords are hashed this approach is better than basic authentication because it means that the slow password check only has to be done periodically — either when creating a token for the first time or after a token has expired.
- The downside is that managing tokens can be complicated for clients — they will need to implement the necessary logic for caching tokens, monitoring and managing token expiry,and periodically generating new tokens.

##### Stateful token authentication

- In a stateful token approach, the value of the token is a high-entropy cryptographically secure random string. **This token — or a fast hash of it — is stored server-side in a database, alongside the user ID and an expiry time for the token**.
- **The big advantage of this is that your API maintains control over the tokens** — it's straightforward to revoke tokens on a per-token or per-user basis by deleting them from the database or marking them as expired.
- Downsides: **you will need to make a database lookup** to check the user's activation status or retrieve additional information about them anyway.

##### Stateless token authentication

- **Stateless tokens encode the user ID and expiry time in the token itself**. The token is cryptographically signed to prevent tampering and (in some cases) encrypted to prevent the contents being read.
- There are a few different technologies that you can use:  a `JWT` (JSON Web Token) is probably the most well-known approach, but `PASETO`, `Branca` and `nacl/secretbox` are viable alternatives too.
- Selling point: **the work to encode and decode the token can be done in memory, and all the information required to identify the user is contained within the token itself**. There's no need to perform a database lookup to find out who a request is coming from.
- Downsides:  they can't easily be revoked once they are issued.
- In an emergency, **you could effectively revoke all tokens by changing the secret used for signing your tokens**. Another workaround is to maintain a blocklist of revoked tokens in a database (although that defeats the 'stateless' aspect of having stateless tokens).

>  **Note:** You should generally avoid storing additional information in a stateless token, such as a user's activation status or permissions, and *using that as the basis for authorization checks*. During the lifetime of the token, **the information encoded** into it will potentially become stale and **out-of-sync with the real data in your system**.

- **Must read**: [Critical vulnerabilities in JSON Web Token libraries](https://curity.io/resources/learn/jwt-best-practices/) and [JWT Security Best Practices](https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/)

#### API Key authentication

-  The idea behind API-key authentication is that a user has a non-expiring secret 'key' associated with their account like:

```bash
 Authorization: Key <key>
```

- On one hand, **this is nice for the client as they can use the same key for every request and they don't need to write code to manage tokens or expiry**. On the other hand, the user now has two long-lived secrets to manage which can potentially compromise their account: their password, and their API key.

#### OAuth 2.0 / OpenID Connect

- With this approach, information about your users (and their passwords) is stored by a **third-party identity provider like Google or Facebook rather than yourself**.
- High level design:
  - Then you want to authenticate a request, you redirect the user to an **`authentication and consent`** form hosted by the identity provider.
  - If the user consents, then the identity provider sends your API an authorization code.
  - Your API then sends the authorization code to another endpoint provided by the identity provider. They verify the authorization code, and if it's valid they will send you a JSON response containing an ID token.
  - This ID token is itself a JWT. You need to validate and decode this JWT to get the actual user information, which includes things like their email address, name, birth date, timezone etc.
  - Now that you know who the user is, you can then implement a stateful or stateless authentication token pattern so that you don't have to go through the whole process for every subsequent request.

#### What authentication approach should I use?

- If your API doesn't have 'real' user accounts with slow password hashes, then HTTP basic authentication can be a good — and often overlooked — fit.
- If you don't want to store user passwords yourself, all your users have accounts with a third-party identity provider that supports OpenID Connect, and your API is the back-end for a website… then use OpenID Connect.
- If you require delegated authentication, such as when your API has a microservice architecture with different services for performing authentication and performing other tasks, then use stateless authentication tokens.
- Otherwise use API keys or stateful authentication tokens. In general:
  - Stateful authentication tokens are a nice fit for APIs that act as the back-end for a website or single-page application, as there is a natural moment when the user logs in where they can be exchanged for user credentials.
  - In contrast, API keys can be better for more 'general purpose' APIs because they're permanent and simpler for developers to use in their applications and scripts.

### Authenticating Requests

- Essentially, once a client has an authentication token we will expect them to include it with all subsequent requests in an Authorization header, like so:

```bash
 Authorization: Bearer IEYZQUBEMPPAKPOAWTPV6YJ6RM
```

-  When we receive these requests, we'll use a new `authenticate()` middleware method to execute the following logic:
   - If the authentication **token is not valid**, then we will send the client a `401 Unauthorized`.
   - If the authentication **token is valid**, we will **look up the user details** and add their details to the request context.
   - If **no Authorization header was provided at all**, then we will add the details for an `anonymous user`.

#### Anonymous User

- Create anonymous user:

```Go
// Define a custom ErrDuplicateEmail error.
var (
	AnonymousUser     = &User{}
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

// Check if a User instance is the AnonymousUser.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}
/*
	data.AnonymousUser.IsAnonymous() // → Returns true

	otherUser := &data.User{}
	otherUser.IsAnonymous()          // → Returns false
*/
```

####  Reading and writing to the request context

- Every http.Request that our application processes has a context.Context embedded in it, which we can use to store key/value pairs containing arbitrary data during the lifetime of the request.
- Any values stored in the request context have the type `interface{}`.
- It's good practice to use your own custom type for the request context keys.

```Go
package main

import (
	"context"
	"greenlight-further/internal/model"
	"net/http"
)

// Define a custom contextKey type, with the underlying type string.
type contextKey string

// Convert the string "user" to a contextKey type and assign it to the userContextKey
// constant. We'll use this constant as the key for getting and setting user information
// in the request context.
const userContextKey = contextKey("user")

// The contextSetUser() method returns a new copy of the request with the provided
// User struct added to the context. Note that we use our userContextKey constant as the
// key.
func (app *application) contextSetUser(r *http.Request, user *model.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// The contextSetUser() retrieves the User struct from the request context. The only
// time that we'll use this helper is when we logically expect there to be User struct
// value in the context, and if it doesn't exist it will firmly be an 'unexpected' error.
// As we discussed earlier in the book, it's OK to panic in those circumstances.
func (app *application) contextGetUser(r *http.Request) *model.User {
	user, ok := r.Context().Value(userContextKey).(*model.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
```

