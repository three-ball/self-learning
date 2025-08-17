# Authenticate & Authorize Users

> Authentication verifies who you are, while authorization determines what you can do.

## Introduction

- Authentication and authorization are fundamental security concepts in backend development. While they are often mentioned together, **they serve distinct purposes in the security framework of an application.**
- In simple terms, authentication is the process of verifying **who a user** is, while authorization is the process of **verifying what they have access to**.



## Authentication
### Definition

- **Authentication** is the process of verifying the identity of a user, system, or entity. It answers the question: "Who are you?". Authentication confirms that users are who they claim to be before granting them access to protected resources.
- Typically, users prove they are who they say they are by entering a password, using a token, or through biometric verification.

### Types of Authentication

#### Basic Authentication

- **Definition**: A simple authentication scheme built into the HTTP protocol, where the client provides a username and password.
- When employing Basic Authentication, users include an encoded string in the `Authorization` header of each request they make. The string is used by the request's recipient to verify users' identity and rights to access a resource.
- When using Basic authentication, the client combines the username and password with a colon (`:`) and encodes it using Base64. The resulting string is then sent in the HTTP header.
- **Example**: If the username is `vienct3` and the password is `1234`, the encoded string would be `dmllbmN0MzoxMjM0`.

```shell
Authorization: Basic dmllbmN0MzoxMjM0
```

- **Pros**: Simple to implement, widely supported, built into HTTP.
- **Cons**:
    - **Credentials are sent in every request**, which can be insecure if not used over HTTPS.
    - **No built-in mechanism for token expiration or revocation**.
    - Not efficient for modern applications that require stateless authentication. **Every request must include credentials, which can lead to performance issues**: Backend servers must validate credentials on every request, which can be resource-intensive.

#### Token-Based Authentication

- **Definition**: Authentication using tokens that are generated after initial verification and used for subsequent requests.

##### Stateful Token

- Server stores the session token in memory or a database. After first authentication, the server generates a session token and sends it to the client. The client includes this token in subsequent requests.
- **How it works:**
    1. User logs in with credentials.
    2. Server validates credentials and creates a session, then stores the session token.
    3. Server sends a session token to the client.
    4. Client includes the session token in subsequent requests.
    5. Server checks the session token against its stored sessions to authenticate the user.

- **Pros**:
    - **Revocation is possible**: If a user logs out or the session expires, the server can invalidate the session token by removing it from its storage.
    - **Full control over session management**: The server can manage session lifetimes, renew tokens, and handle user sessions more flexibly.
    - Token itsefl contains minimal information, reducing the amount of data sent with each request and sensitive data exposure.
- **Cons**:
    - **Requires server-side storage for session tokens**, which can lead to scalability issues. Requires looking up the session token in a database or memory store on every request, which can add latency.
    - **Session state management can be complex**, especially in distributed systems where multiple servers need to access the same session store.

```mermaid
sequenceDiagram
    participant Client
    participant Server
    participant SessionStore as Session Store<br/>(Memory/Database)
    
    Note over Client, SessionStore: Initial Authentication
    Client->>Server: 1. Login with credentials (username/password)
    Server->>Server: 2. Validate credentials
    Server->>SessionStore: 3. Create and store session token
    SessionStore-->>Server: Session token stored
    Server->>Client: 4. Send session token
    
    Note over Client, SessionStore: Subsequent Requests
    Client->>Server: 5. Request with session token
    Server->>SessionStore: 6. Check token validity
    SessionStore-->>Server: Token validation result
    
    alt Token is valid
        Server->>Client: Authenticated response
    else Token is invalid/expired
        Server->>Client: 401 Unauthorized
    end
    
    Note over Client, SessionStore: Session Revocation
    Client->>Server: Logout request
    Server->>SessionStore: Remove/invalidate session token
    SessionStore-->>Server: Token removed
    Server->>Client: Logout confirmation
```

##### Stateless Token

- Self-contained tokens where all necessary information is encoded within the token itself.
- After successful authentication, the application generates a token with all necessary data, signs it with a public key, and sends it back to a client. There is a standard for token generation, it is **JWT (JSON Web Token)**.

```shell
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30
```

- **How it works:**
    1. User logs in with credentials.
    2. Server validates credentials and signs a JWT containing user information and claims with a secret key.
    3. Server sends the JWT to the client.
    4. Client includes the JWT in subsequent requests.
    5. Server verifies the JWT signature and extracts user information.

```mermaid
sequenceDiagram
    participant Client
    participant Server
    
    Note over Client, Server: Initial Authentication
    Client->>Server: 1. Login with credentials (username/password)
    Server->>Server: 2. Validate credentials
    Server->>Server: 3. Generate JWT with user claims and sign with secret key
    Server->>Client: 4. Send JWT token
    
    Note over Client, Server: Subsequent Requests
    Client->>Server: 5. Request with JWT in Authorization header
    Server->>Server: 6. Verify JWT signature and extract claims
    
    alt JWT is valid and not expired
        Server->>Client: Authenticated response with user data
    else JWT is invalid/expired/tampered
        Server->>Client: 401 Unauthorized
    end
    
    Note over Client, Server: No Server-Side Session Storage Required
    Note over Server: JWT contains all necessary user information<br/>No database lookup needed for authentication
```


- Some information about JWT: [JSON Web Token](https://auth0.com/docs/secure/tokens/json-web-tokens)
    - The information contained within the JSON object can be verified and trusted because it is digitally signed using a secret or public/private key pair.
    - A well-formed JWT consists of three parts seperated by dots (`.`): `Header`, `Payload`, and `Signature`.
        - **Header**: Contains metadata about the token, such as the signing algorithm used (e.g., HMAC SHA256).
        - **Payload**: The payload contains statements about the entity (typically, the user) and additional entity attributes, which are called `claims`. In this example, our entity is a user.
        - **Signature**: Used to validate that the token is trustworthy and has not been tampered with. When you use a JWT, you must check its signature before storing and using it.
    - Token claims can include: [JSON Web Token Claims](https://www.iana.org/assignments/jwt/jwt.xhtml#claims)
        - `iss` (Issuer): The entity that issued the token.
        - `sub` (Subject): The subject of the token, typically the user ID.
        - `aud` (Audience): Recipient for which the JWT is intended
        - `exp` (Expiration Time): The time after which the token is no longer valid.
        - `iat` (Issued At): The time at which the token was issued.
        - `nbf` (Not Before): Time before which the JWT must not be accepted for processing.
        - `jti` (JWT ID): A unique identifier for the token.

- **Pros**:
    - **Stateless**: No need for server-side session storage, reducing server load and improving scalability.
    - **Reduced database lookups**: Since the token contains all necessary information, the server does not need to query a database for user information on every request.
    - **Can work across distributed systems**: Tokens can be easily shared between different services or microservices without requiring a central session store.
- **Cons**:
    - **Token revocation can be complex**: If a user logs out or if a token needs to be invalidated, the server cannot simply remove it from a session store. Instead, it may require additional mechanisms like token blacklisting or short-lived tokens with refresh tokens.
    - **Token size can be larger**: Since the token contains all user information, it can be larger than a simple session ID, which may impact performance if not managed properly.

- **Hyrid Approach**: 
    - Short-lived JWTs for API authentication.
    - Server-side tracking of token revocation in a blocklist.
    - Refresh token patterns where a stateful refresh token is used to get new stateless access tokens.

## Authorization
### Definition
### Types of Authorization

## Implementation in Go
### Basic Authentication
### JWT (JSON Web Tokens)
### OAuth 2.0
### Role-Based Access Control (RBAC)

## Usecase & Best Practices

## References