```markdown
# RestBadger

A lightweight NoSQL database server that provides a REST API interface for BadgerDB operations. Built with Go and Gin, offering SQL-like commands through HTTP endpoints with JSON responses.

## Features

- 🚀 **Fast REST API** - Built with Gin framework
- 🔑 **Key-Value Storage** - Powered by BadgerDB
- 📝 **SQL-like Commands** - Familiar syntax (SELECT, INSERT, UPDATE, DELETE, LIKE)
- 🔄 **Backup & Restore** - Built-in database backup functionality
- 🛡️ **API Key Authentication** - Secure access control
- 🌐 **CORS Enabled** - Ready for web applications
- 🔍 **Prefix Queries** - Search with pagination support
- 🗜️ **Database Compression** - Automatic garbage collection
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
 "dbport": "3308",
 "directory": "./database"
}
```

## API Usage

All requests are made via POST to `http://localhost:3308/data`

### Request Format

```json
{
  "query": "COMMAND key_or_pattern",
  "apikey": "your_api_key",
  "values": [optional_array_of_values]
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
  "query": "INSERT user_123",
  "apikey": "gorms",
  "values": ["John", "Doe", 25, "john@email.com"]
}
```

### SELECT
Retrieve data by key:

```json
{
  "query": "SELECT user_123",
  "apikey": "gorms"
}
```

### UPDATE
Update existing data:

```json
{
  "query": "UPDATE user_123",
  "apikey": "gorms",
  "values": ["John", "Smith", 26, "johnsmith@email.com"]
}
```

### DELETE
Delete data by key:

```json
{
  "query": "DELETE user_123",
  "apikey": "gorms"
}
```

### LIKE
Search by prefix with pagination:

```json
{
  "query": "LIKE user_",
  "apikey": "gorms",
  "values": [10]
}
```

### BACKUP
Create database backup:

```json
{
  "query": "BACKUP backup_20241201.bak",
  "apikey": "gorms"
}
```

### RESTORE
Restore from backup:

```json
{
  "query": "RESTORE backup_20241201.bak",
  "apikey": "gorms"
}
```

### COMPRESS
Run database garbage collection:

```json
{
  "query": "COMPRESS",
  "apikey": "gorms"
}
```

## Examples

### JavaScript/Fetch

```javascript
async function insertUser(userData) {
  const response = await fetch('http://localhost:3308/data', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      query: 'INSERT user_' + userData.id,
      apikey: 'gorms',
      values: [userData.name, userData.email, userData.age]
    })
  });
  return await response.json();
}
```

### cURL

```bash
curl -X POST http://localhost:3308/data \
  -H "Content-Type: application/json" \
  -d '{
    "query": "SELECT user_123",
    "apikey": "gorms"
  }'
```

## Error Handling

The API returns appropriate HTTP status codes:

- `200` - Success
- `400` - Bad Request (invalid syntax, query too long)
- `401` - Unauthorized (invalid API key)
- `404` - Not Found (key doesn't exist)
- `500` - Internal Server Error

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `apikey` | API authentication key | `"gorms"` |
| `dbport` | Server port | `"3308"` |
| `directory` | Database storage path | `"./database"` |

## Database Recovery

RestBadger includes automatic recovery mechanisms:

1. **Normal Open** - Standard database opening
2. **Bypass Lock Guard** - Recovery from improper shutdown
3. **Read-Only Mode** - Safe access when write-locked
4. **Minimal Settings** - Reduced memory configuration
5. **Backup & Recreate** - Last resort recovery

## Performance Tips

- Use meaningful key prefixes for efficient LIKE queries
- Run COMPRESS periodically to optimize storage
- Keep values array reasonable in size
- Use appropriate pagination limits in LIKE queries

## Use Cases

- **Microservices** - Lightweight data storage
- **Prototyping** - Quick database setup
- **IoT Applications** - Embedded systems storage
- **Cache Layer** - Fast key-value caching
- **Session Storage** - User session management
- **Configuration Storage** - Application settings

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
```