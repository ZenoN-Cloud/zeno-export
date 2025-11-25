# Deployment Guide

## Quick Start

### Local Development
```bash
# Build WASM module
make build-wasm

# Serve locally
make serve
# Open http://localhost:8080
```

### Docker (Single Service)
```bash
# Build and run
make docker-run
# Open http://localhost:8080
```

### Docker Compose (Full Stack)
```bash
# Start both zeno-engine and zeno-export
make docker-compose-up
# zeno-export: http://localhost:8080
# zeno-engine: http://localhost:8081
```

## Production Deployment

### Environment Variables
- `NGINX_HOST` - Server hostname (default: localhost)
- `NGINX_PORT` - Server port (default: 80)

### Health Check
- Endpoint: `/health`
- Response: `200 OK`

### CORS Configuration
The nginx configuration includes proper CORS headers for WASM:
- `Cross-Origin-Embedder-Policy: require-corp`
- `Cross-Origin-Opener-Policy: same-origin`

## Integration

### With zeno-engine
```javascript
// 1. Normalize CSV
const normalized = await zenoEngine.normalizeCSV(csvData, '{}');

// 2. Generate Excel
const excelBytes = await zenoExport.toExcel(JSON.stringify(normalized));

// 3. Download
const blob = new Blob([excelBytes], { 
    type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' 
});
const url = URL.createObjectURL(blob);
const a = document.createElement('a');
a.href = url;
a.download = 'transactions.xlsx';
a.click();
```

## Monitoring

### Logs
```bash
# Docker logs
docker logs zeno-export

# Docker Compose logs
docker-compose logs -f zeno-export
```

### Metrics
- WASM module size: ~8.4MB
- Memory usage: <50MB
- Cold start: <1s