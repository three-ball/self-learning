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
- The `internal` directory will contain various ancillary packages used by our API. It will contain the code for interacting with our database, doing data validation, sending emails and so on. Basically, any code which isn’t application-specific and can potentially be reused will live in here. Our Go code under cmd/api will import the packages in the internal directory (but never the other way around).
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
- If we set ConnMaxLifetime to one hour, for example, it means that all connections will be marked as ‘expired’ one hour after they were first created, and cannot be reused after they’ve expired.
  - This doesn’t guarantee that a connection will exist in the pool for a whole hour; it’s possible that a connection will become unusable for some reason and be automatically closed before then.
  - A connection can still be in use more than one hour after being created — it just cannot start to be reused after that time.
  - **This isn’t an idle timeout**. The connection will expire one hour after it was first created — not one hour after it last became idle.
  - Once every second Go runs a background cleanup operation to remove expired connections from the pool.

### SetConnMaxIdleTime()

- This works in a very similar way to ConnMaxLifetime , except it sets the maximum length of time that a connection can be idle for before it is marked as expired. By default there’s no limit.

## Handling Partial Updates

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
- For most applications, on most endpoints, it’s probably not. The exceptions are probably applications which frequently run close to saturation point of their resources, or for specific endpoints which execute slow running or very computationally expensive SQL queries. In those cases, canceling queries aggressively when a client disappears may have a meaningful positive impact and make it worth the trade-off.

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

- This SQL query is designed so that each of the filters behaves like it is ‘optional’. For example, the condition `(LOWER(title) = LOWER($1) OR $1 = '')` will evaluate as `true` if the placeholder parameter `$1` is a case-insensitive match for the movie title or the placeholder parameter equals `''`.
- So this filter condition will essentially be ‘skipped’ when movie title being searched for is the empty string "".

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