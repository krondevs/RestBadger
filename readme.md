# RestBadger

A lightweight NoSQL database server that provides a REST API interface for BadgerDB operations. Built with Go and Gin, offering SQL-like commands through HTTP endpoints with JSON responses, multi-database support, encryption, and TTL capabilities.

## Features

- 🚀 **Fast REST API** - Built with Gin framework
- 🔑 **Key-Value Storage** - Powered by BadgerDB
- 📝 **SQL-like Commands** - Familiar syntax (SELECT, INSERT, UPDATE, DELETE, LIKE)
- 🗄️ **Multi-Database Support** - Manage multiple isolated databases
- 🔒 **AES-256 Encryption** - Optional database encryption with master keys
- ⏰ **TTL Support** - Automatic data expiration
- 🔄 **Backup & Restore** - Built-in database backup functionality
- 🛡️ **API Key Authentication** - Secure access control
- 🌐 **CORS Enabled** - Ready for web applications
- 🔍 **Prefix Queries** - Search with pagination support
- 🗜️ **Database Compression** - Manual garbage collection
- ⚡ **Auto Recovery** - Smart database recovery mechanisms
- 📦 **Single Binary** - Easy deployment

## Quick Start

### Installation

```bash
git clone https://github.com/yourusername/restbadger
cd restbadger
go mod tidy
go run main.go
```

### Configuration

On first run, `dbconfig.json` will be created automatically:

```json
{
 "apikey": "gorms",
 "dbport": "3308"
}
```

## API Usage

All requests are made via POST to `http://localhost:3308/data`

### Request Format

```json
{
  "query": "COMMAND",
  "apikey": "your_api_key",
  "key": "record_key",
  "values": [optional_array_of_values],
  "db": "database_name",
  "ttl": 3600,
  "encrypt": "32-byte-encryption-key-here-12345"
}
```

### Response Format

```json
{
  "status": "success|error",
  "message": "description",
  "result": []
}
```

## Commands

### INSERT
Insert new data with a unique key:

```json
{
  "query": "INSERT",
  "apikey": "gorms",
  "key": "user_123",
  "values": ["John", "Doe", 25, "john@email.com"],
  "db": "users"
}
```

### INSERT with TTL
Insert data with automatic expiration:

```json
{
  "query": "INSERT",
  "apikey": "gorms",
  "key": "session_abc123",
  "values": ["user_456", "token_xyz", "active"],
  "db": "sessions",
  "ttl": 3600
}
```

### INSERT with Encryption
Insert data in an encrypted database:

```json
{
  "query": "INSERT",
  "apikey": "gorms",
  "key": "sensitive_data_123",
  "values": ["confidential", "information"],
  "db": "secure_db",
  "encrypt": "my-32-byte-encryption-key-here12"
}
```

### SELECT
Retrieve data by key:

```json
{
  "query": "SELECT",
  "apikey": "gorms",
  "key": "user_123",
  "db": "users"
}
```

### UPDATE
Update existing data:

```json
{
  "query": "UPDATE",
  "apikey": "gorms",
  "key": "user_123",
  "values": ["John", "Smith", 26, "johnsmith@email.com"],
  "db": "users"
}
```

### DELETE
Delete data by key:

```json
{
  "query": "DELETE",
  "apikey": "gorms",
  "key": "user_123",
  "db": "users"
}
```

### LIKE
Search by prefix with pagination:

```json
{
  "query": "LIKE",
  "apikey": "gorms",
  "key": "user_",
  "values": [10],
  "db": "users"
}
```

### BACKUP
Create database backup with timestamp:

```json
{
  "query": "BACKUP",
  "apikey": "gorms",
  "db": "users"
}
```

### RESTORE
Restore from backup:

```json
{
  "query": "RESTORE",
  "apikey": "gorms",
  "key": "backups/users_2024-12-01_10_30_00.bak",
  "db": "users_restored"
}
```

### COMPRESS
Run database garbage collection:

```json
{
  "query": "COMPRESS",
  "apikey": "gorms",
  "db": "users"
}
```

## Multi-Database Support

RestBadger supports multiple isolated databases. Each database is identified by the `db` field in your requests:

```json
// Users database
{
  "query": "INSERT",
  "key": "user_123",
  "values": ["John", "Doe"],
  "db": "users",
  "apikey": "gorms"
}

// Products database
{
  "query": "INSERT", 
  "key": "product_456",
  "values": ["Laptop", 999.99],
  "db": "products",
  "apikey": "gorms"
}

// Sessions database with TTL
{
  "query": "INSERT",
  "key": "session_789",
  "values": ["active_session_data"],
  "db": "sessions",
  "ttl": 1800,
  "apikey": "gorms"
}
```

## Encryption

RestBadger supports AES-256 encryption for sensitive data. Use a 32-byte key for encryption:

### Generating Encryption Keys

```bash
# Generate a secure 32-byte key (example)
key="my-super-secret-32-byte-key-here"
```

### Using Encrypted Databases

```json
{
  "query": "INSERT",
  "apikey": "gorms",
  "key": "confidential_123",
  "values": ["secret", "data"],
  "db": "encrypted_db",
  "encrypt": "RestBadger-Test-Key-32-bytes-long"
}
```

**Important Notes:**
- The same encryption key must be used for all operations on an encrypted database
- Backups of encrypted databases remain encrypted
- Losing the encryption key means losing access to the data permanently

## TTL (Time To Live)

Data can automatically expire after a specified time:

```json
{
  "query": "INSERT",
  "apikey": "gorms", 
  "key": "temp_data_123",
  "values": ["temporary", "information"],
  "db": "cache",
  "ttl": 300
}
```

TTL is specified in seconds. After expiration, data is automatically removed during garbage collection.

## Examples

### JavaScript/Fetch

```javascript
// Insert user data
async function insertUser(userData) {
  const response = await fetch('http://localhost:3308/data', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      query: 'INSERT',
      apikey: 'gorms',
      key: 'user_' + userData.id,
      values: [userData.name, userData.email, userData.age],
      db: 'users'
    })
  });
  return await response.json();
}

// Insert with TTL (session data)
async function createSession(sessionData) {
  const response = await fetch('http://localhost:3308/data', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      query: 'INSERT',
      apikey: 'gorms',
      key: 'session_' + sessionData.id,
      values: [sessionData.userId, sessionData.token],
      db: 'sessions',
      ttl: 3600 // 1 hour
    })
  });
  return await response.json();
}

// Insert encrypted data
async function insertSecureData(data) {
  const response = await fetch('http://localhost:3308/data', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      query: 'INSERT',
      apikey: 'gorms',
      key: 'secure_' + data.id,
      values: [data.confidential, data.sensitive],
      db: 'vault',
      encrypt: 'your-32-byte-encryption-key-here12'
    })
  });
  return await response.json();
}
```

### cURL

```bash
# Basic insert
curl -X POST http://localhost:3308/data \
  -H "Content-Type: application/json" \
  -d '{
    "query": "INSERT",
    "apikey": "gorms",
    "key": "test_123",
    "values": ["hello", "world"],
    "db": "test_db"
  }'

# Insert with TTL
curl -X POST http://localhost:3308/data \
  -H "Content-Type: application/json" \
  -d '{
    "query": "INSERT",
    "apikey": "gorms", 
    "key": "cache_456",
    "values": ["cached", "data"],
    "db": "cache",
    "ttl": 300
  }'

# Select data
curl -X POST http://localhost:3308/data \
  -H "Content-Type: application/json" \
  -d '{
    "query": "SELECT",
    "apikey": "gorms",
    "key": "test_123",
    "db": "test_db"
  }'
```

## File Structure

RestBadger creates the following directory structure:

```
./
├── dbconfig.json          # Configuration file
├── databases/             # Database files
│   ├── users/            # User database
│   ├── products/         # Products database
│   └── sessions/         # Sessions database
└── backups/              # Automatic backups
    ├── users_2024-12-01_10_30_00.bak
    └── products_2024-12-01_10_31_15.bak
```

## Error Handling

The API returns appropriate HTTP status codes:

- `200` - Success
- `400` - Bad Request (invalid syntax, missing database name)
- `401` - Unauthorized (invalid API key)
- `404` - Not Found (key doesn't exist)
- `500` - Internal Server Error

## Use Cases

- **Microservices** - Lightweight data storage with database isolation
- **Session Management** - TTL-based session storage
- **Caching Layer** - Fast key-value caching with expiration
- **Configuration Storage** - Application settings per environment
- **Secure Storage** - Encrypted sensitive data
- **Multi-tenant Applications** - Isolated databases per tenant
- **IoT Applications** - Embedded systems storage
- **Prototyping** - Quick database setup for development

## Performance Tips

- Use meaningful key prefixes for efficient LIKE queries
- Run COMPRESS periodically to optimize storage
- Use appropriate TTL values to prevent data accumulation
- Consider encryption overhead for performance-critical applications
- Use separate databases for different data types or tenants

## Security Considerations

- **API Keys**: Change the default API key in production
- **Encryption Keys**: Use strong 32-byte keys for encrypted databases
- **Key Management**: Store encryption keys securely, separate from the application
- **Network**: Use HTTPS in production environments
- **Access Control**: Implement additional authentication layers if needed

## Requirements

- Go 1.19 or higher
- BadgerDB v4
- Gin framework

## License

MIT License - see LICENSE file for details

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

For issues and questions, please open an issue on GitHub.