# Ticket Dummy Data Generator

This project generates dummy ticket data for learning about database indexing and performance optimization.

## Files

- `tickets.sql` - Database schema with comprehensive indexing strategies
- `main.go` - Go application to generate and insert 2 million random ticket records
- `go.mod` - Go module dependencies

## Prerequisites

1. **PostgreSQL Database**
   - Install PostgreSQL on your system
   - Create a database named `tickets_db`
   - Update connection parameters in `main.go` if needed

2. **Go Environment**
   - Go 1.21 or later

## Setup Instructions

### 1. Database Setup

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE tickets_db;

# Connect to the database
\c tickets_db

# Run the schema
\i tickets.sql
```

### 2. Go Application Setup

```bash
# Initialize Go modules and download dependencies
go mod tidy

# Run the application
go run main.go
```

## Database Schema

The `tickets` table includes the following fields designed for indexing practice:

### Core Fields
- `id` - Primary key (BIGSERIAL)
- `ticket_number` - Unique identifier (VARCHAR)
- `title` - Ticket title (VARCHAR)
- `description` - Detailed description (TEXT)
- `status` - Current status (VARCHAR)
- `priority` - Priority level (VARCHAR)
- `category` - Ticket category (VARCHAR)

### User Fields
- `user_id` - Customer ID (BIGINT)
- `assigned_to` - Agent ID (BIGINT)

### Timestamps
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp
- `resolved_at` - Resolution timestamp (nullable)
- `due_date` - Due date (nullable)

### Additional Fields
- `tags` - Array of tags (TEXT[])
- `metadata` - JSONB metadata
- `is_escalated` - Escalation flag (BOOLEAN)
- `customer_satisfaction_score` - Rating 1-5 (INTEGER)
- `response_time_hours` - Response time in hours (INTEGER)
- `resolution_time_hours` - Resolution time in hours (INTEGER)

## Indexing Strategy

The schema includes various types of indexes:

### Single Column Indexes
- Status, priority, category, user_id, assigned_to
- Created_at, updated_at for time-based queries

### Composite Indexes
- (status, priority) - Common filtering combination
- (category, status) - Category-specific status queries
- (user_id, status, created_at) - User-specific ticket history
- (assigned_to, status, priority) - Agent workload queries

### Partial Indexes
- Open high-priority tickets
- Unassigned tickets
- Escalated tickets

### Specialized Indexes
- GIN indexes for JSONB metadata and array tags

## Performance Testing Queries

After generating the data, try these queries to test index performance:

```sql
-- Test single column index
EXPLAIN ANALYZE SELECT * FROM tickets WHERE status = 'open';

-- Test composite index
EXPLAIN ANALYZE SELECT * FROM tickets WHERE status = 'open' AND priority = 'high';

-- Test partial index
EXPLAIN ANALYZE SELECT * FROM tickets WHERE status = 'open' AND priority = 'high' ORDER BY created_at;

-- Test JSONB index
EXPLAIN ANALYZE SELECT * FROM tickets WHERE metadata->>'source' = 'web';

-- Test array index
EXPLAIN ANALYZE SELECT * FROM tickets WHERE tags @> ARRAY['bug'];

-- Test time-based queries
EXPLAIN ANALYZE SELECT * FROM tickets WHERE created_at >= '2024-01-01' AND created_at < '2024-02-01';
```

## Customization

You can modify the following in `main.go`:

- **Database connection parameters** (lines 15-21)
- **Total records to generate** (change `totalRecords` constant)
- **Batch size** (change `batchSize` constant)
- **Random data values** (modify the arrays at the top)
- **Ticket generation logic** (modify `generateRandomTicket` function)

## Performance Notes

- The application uses batch inserts with transactions for better performance
- Connection pooling is configured for optimal database connections
- Progress is displayed every 10,000 records
- Average insertion rate: ~5,000-10,000 records/second (depends on hardware)

## Learning Objectives

This setup allows you to:

1. **Practice Index Design** - Understand different index types and their use cases
2. **Query Performance Analysis** - Use EXPLAIN ANALYZE to understand query execution plans
3. **Index Effectiveness** - Compare query performance with and without indexes
4. **Composite Index Optimization** - Learn column ordering in multi-column indexes
5. **Partial Index Usage** - Optimize for specific query patterns
6. **JSONB and Array Indexing** - Work with modern PostgreSQL features

## Troubleshooting

- **Connection Issues**: Verify PostgreSQL is running and connection parameters are correct
- **Memory Issues**: Reduce batch size if you encounter memory problems
- **Slow Insertions**: Check database configuration and available system resources
- **Index Conflicts**: Drop and recreate indexes if you modify the schema
