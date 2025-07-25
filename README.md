# SecuSend Backend

[![Github issues](https://img.shields.io/github/issues/secusend/secusend-backend?logo=github)](https://github.com/secusend/secusend-backend/issues)
[![Github Pull Requests](https://img.shields.io/github/issues-pr/secusend/secusend-backend?logo=github)](https://github.com/secusend/secusend-backend/pulls)
[![Github License](https://img.shields.io/github/license/secusend/secusend-backend)](https://github.com/secusend/secusend-backend/blob/master/LICENSE)

The backend API for SecuSend - a secure, self-destructing note sharing application. This Go-based REST API handles encrypted note storage, retrieval, and automatic cleanup with MongoDB as the database.

## Live Demo

The backend powers the live application at: [https://secusend.eu.org](https://secusend.eu.org)

## Features

- **AES-256 Encryption**: Notes are encrypted using AES-256-GCM with PBKDF2 key derivation
- **Self-Destructing Notes**: Automatic deletion after first read
- **Time-Based Expiration**: Configurable expiration times (1h, 1d, 1w, 1m, 1y)
- **Rate Limiting**: Built-in protection against abuse (3 requests/second per IP)
- **Automatic Cleanup**: Daily cron job removes expired notes
- **Security Headers**: Helmet middleware for enhanced security
- **CORS Support**: Cross-origin resource sharing enabled
- **Request Compression**: Gzip/Brotli compression for optimal performance

## Installation

### Prerequisites

- Go 1.24.3 or higher
- MongoDB instance
- Git

### Local Development

1. **Clone the repository:**
   ```bash
   git clone https://github.com/secusend/secusend-backend.git
   cd secusend-backend
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up environment variables in `.env`:**
   ```bash
   MONGOURI=mongodb://...
   ```

4. **Run the application:**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:3000`

### Development with Air (Hot Reload)

For development with automatic reloading:

```bash
# Install Air
go install github.com/air-verse/air@latest

# Run with hot reload
air
```

## Docker Deployment

### Run with Docker

```bash
docker build -t secusend-backend .

docker run -d \
  --name secusend-backend \
  -p 3000:3000 \
  -e MONGOURI="your-mongodb-connection-string" \
  secusend-backend
```

### Docker Compose

```yaml
version: '3.8'
services:
  secusend-backend:
    build: .
    ports:
      - "3000:3000"
    environment:
      - MONGOURI=mongodb://mongo:27017
    depends_on:
      - mongo
  
  mongo:
    image: mongo:latest
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db

volumes:
  mongo_data:
```

## Security Features

- **AES-256-GCM Encryption**: Industry-standard encryption for password-protected notes
- **PBKDF2 Key Derivation**: Secure password-based key generation with 10,000 iterations
- **Rate Limiting**: 3 requests per second per IP address
- **Input Validation**: Size limits and sanitization
- **Security Headers**: Helmet middleware for XSS, CSRF protection
- **Unique Key Generation**: Cryptographically secure random keys

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Related Projects

- [SecuSend Frontend](https://github.com/secusend/secusend-frontend) - Vue web interface

## License

This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.