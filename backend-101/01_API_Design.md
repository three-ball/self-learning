# 01. REST API Design

## 1.1. Mindset

- Scalable.
- Consistency.
- Inspect every single aspact of the API.
- No one fits all (Trade-offs).
- Some Resources:
    - [Stripe API Documentation](https://docs.stripe.com/api)
    - [GitLab API Documentation](https://docs.gitlab.com/api/rest/)

## 2.1 HTTP Methods

> REST APIs are built on top of HTTP, which provides a sest of methods to interact with resources. Each method has specific characteristics that determine how it should be used.

- Properties:
  - Idempotency: A same request can be repeated multiple times without changing the result.
  - Safety: No side effects on the server.

- Methods:

| Method | Description | Idempotent | Safe |
|--------|-------------|------------|------|
| GET    | Retrieve data | Yes        | Yes  |
| HEAD   | Retrieve headers | Yes        | Yes  |
| OPTIONS| Describe communication options | Yes        | Yes  |
| TRACE  | Echo back the received request | Yes        | Yes  |
| PUT    | Update resource completely | Yes        | No   |
| DELETE | Remove resource | Yes        | No   |
| POST   | Create resource | No         | No   |
| PATCH  | Update resource partially | No         | No   |

## 2.2. Conventions

> Just some common conventions to follow when designing REST APIs.

- Use Nouns instead of Verbs:
  - Use `/users` instead of `/getUsers`.
  - Use `/users/123` instead of `/getUserById/123`.
- Use Plural Nouns:
  - Use `/users` instead of `/user`.
- Use Nesting for Hierarchical Resources:
  - Use `/users/123/orders` instead of `/orders?userId=123`.
- Versioning:
  - Use `/v1/users` instead of `/users/v1`.
- Slug-case for URLs:
  - https://vienct-blog.com/slug-case-vs-camel-case-in-urls/
- Snake_case for resquest, response, and query parameters:
  - Use `first_name` instead of `firstName`.

## 2.3. Pagination

> Pagination is a technique to limit the number of results returned by an API endpoint, which is especially useful for large datasets.

- Page, size parameters:
  - Example: `/users?page=1&size=10`.
  - Usecase: Management portal, admin dashboard, etc.
  - **Must document**: Page starts at 1 or 0, the default page size, and the maximum page size.
- Offset, limit parameters:
  - Example: `/users?offset=0&limit=10`.
  - Usecase: an infinite scroll, news feed, etc.

***Question:***
```sql
SELECT * FROM users OFFSET 100 LIMIT 10`
```

This common SQL lead to the problem: `Resource skipping`.
- If a resource is deleted, the next page will skip that resource.
- Example:
- Page 1: `1, 2, 3, 4, 5`
- Page 2: `6, 7, 8, 9, 10`
- If resource `1` and `2` is deleted:
    - `6, 7` move to page 1: `3, 4, 5, 6, 7`
    - Page 2: `8, 9, 10, 11, 12`
- Solution: Use a cursor-based pagination (not suitable for random IDs).

```SQL
SELECT * FROM users WHERE id > [last_id] ORDER BY id ASC LIMIT 10
```

- Some other tricks: Performance optimization using `deferred joining`:

```sql
SELECT * FROM 
(
    SELECT * FROM users WHERE id > [last_id] ORDER BY id ASC LIMIT 10
) AS paginated_users
JOIN orders ON paginated_users.id = orders.user_id
```

## 2.4. Sorting

> Sorting allows clients to specify the order in which results should be returned.

- Note: We **must** white list of sortable fields. Indexing is important for performance.
- Examples:

```plaintext
/users?sort=created_at:asc,name:desc
/users?sort=+created_at,-name
```

## 2.5. Some Usecases

### 2.5.1. Export heavy data

**Using polling**:
  - Client sends a request to export data.
  - Server processes the request and returns a job ID.
  - Client polls the server using the job ID to check if the export is ready.
  - Once ready, the client can download the exported data.
- Advantages:
  - Simple to implement.
  - No need for WebSockets or long-polling.
- Disadvantages:
  - Polling can be inefficient, especially if the export takes a long time.
  - Clients may need to implement retry logic for failed requests.

```sh
CLIENT                                             SERVER
  |                                                 |
  | POST /api/exports                               |
  |─────────────────────────────────────────────>   |
  |                                                 | (Create job)
  |    { job_id: "exp_123", status: "pending" }     |
  |<────────────────────────────────────────        |
  |                                                 |
  | GET /api/exports/exp_123 (Poll #1)              |
  |─────────────────────────────────────────────>   |
  |                                                 | |████░░░░░░| 40%
  |    { status: "processing", progress: 40% }      |
  |<────────────────────────────────────────        |
  |                                                 |
  | (Wait 5 seconds)                                |
  |                                                 |
  | GET /api/exports/exp_123 (Poll #2)              |
  |─────────────────────────────────────────────>   |
  |                                                 | |███████░░░| 70%
  |    { status: "processing", progress: 70% }      |
  |<────────────────────────────────────────        |
  |                                                 |
  | (Wait 5 seconds)                                |
  |                                                 |
  | GET /api/exports/exp_123 (Poll #3)              |
  |─────────────────────────────────────────────>   |
  |                                                 | |██████████| 100%
  |    { status: "completed", download_url: "..." } |
  |<────────────────────────────────────────        |
  |                                                 |
  | GET /api/exports/exp_123/download               |
  |─────────────────────────────────────────────>   |
  |                                                 |
```

**Using Callbacks**:
  - Client sends a request to export data.
  - Server processes the request and calls a callback URL provided by the client when the export is ready.
- Advantages:
  - More efficient than polling, as the server notifies the client when the export is ready.
  - Reduces unnecessary requests to the server.
- Disadvantages:
  - Requires the client to expose a callback URL, which may not be feasible in all scenarios.
  - More complex to implement, as it involves handling callbacks and ensuring security.
  - Client and Server will couple together.

```sh
    CLIENT                                             SERVER
    |                                                   |
    | POST /api/exports                                 |
    |─────────────────────────────────────────────>     |
    |                                                   | (Create job)
    |                                                   |
    |    ACK                                            |
    |<────────────────────────────────────────          |
    |                                                   |
    | POST /api/callbacks/export_completed              |
    |<────────────────────────────────────────          |
    |                                                   |
    |                                                   |
    |                                                   |
    |                                                   |
    | GET /api/exports/exp_123/download                 |
    |─────────────────────────────────────────────>     |
    |                                                   |