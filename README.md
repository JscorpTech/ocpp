# OCPP Server

Go tilida yozilgan OCPP (Open Charge Point Protocol) 1.x serveri. Elektr avtomobil zaryadlash stantsiyalari bilan aloqa qilish uchun ishlatiladi.

## Xususiyatlar

- ✅ OCPP 1.x protokolini qo'llab-quvvatlash
- ✅ WebSocket orqali real-time aloqa
- ✅ Redis orqali event management va remote commands
- ✅ Transaction boshqaruvi
- ✅ Heartbeat va health monitoring
- ✅ Meter values va charging session tracking

## Arxitektura

```
.
├── cmd/
│   └── main.go           # Dastur kirish nuqtasi
├── internal/
│   ├── client/           # HTTP clientlar
│   ├── config/           # Konfiguratsiya
│   ├── domain/           # Domain modellar va eventlar
│   ├── ocpp/            # OCPP server va handlerlar
│   └── services/         # Redis va boshqa servislar
└── pkg/                  # Public kutubxonalar
```

## O'rnatish

### Talablar

- Go 1.24+
- Redis 6.0+

### Dependencies o'rnatish

```bash
go mod download
```

## Konfiguratsiya

`.env` fayl yarating:

```env
BASE_URL=http://your-backend-api:8000
ADDR=:10800
```

**Environment Variables:**

- `BASE_URL` - Backend API URL (majburiy)
- `ADDR` - Server manzil (default: `:10800`)

## Ishga tushirish

### Development

```bash
go run ./cmd/main.go
```

### Production

```bash
# Build
make build

# Run
./bin/ocpp
```

## Redis Integration

### Events Queue

Server barcha eventlarni `events` qatoriga yuboradi:

```json
{
  "event": "start_transaction",
  "data": {
    "charger": "charger-001",
    "conn": 1,
    "tag": "RFID-12345",
    "meter_start": 0
  }
}
```

**Event turlari:**

- `health` - Heartbeat
- `change_connector_status` - Konektor holati o'zgarishi
- `start_transaction` - Zaryadlash boshlanishi
- `stop_transaction` - Zaryadlash tugashi
- `meter_value` - Elektr o'lchov ma'lumotlari

### Remote Commands

Backend `commands` qatoriga komandalar yuborishi mumkin:

```json
{
  "CpID": "charger-001",
  "data": [2, "unique-id", "RemoteStartTransaction", {"connectorId": 1, "idTag": "RFID-12345"}]
}
```

## OCPP Handlers

Server quyidagi OCPP xabarlarini qabul qiladi:

| Handler | Tavsif |
|---------|--------|
| `BootNotification` | Charger ulanishi |
| `StatusNotification` | Konektor holati |
| `Authorize` | RFID avtorizatsiya |
| `Heartbeat` | Health check |
| `StartTransaction` | Zaryadlash boshlanishi |
| `StopTransaction` | Zaryadlash tugashi |
| `MeterValues` | Elektr o'lchov ma'lumotlari |

## Testing

### Barcha testlar

```bash
make test
```

### Unit testlar

```bash
make test-unit
```

### Coverage

```bash
make test-coverage
```

Coverage hisobotni `coverage.html` faylida ko'rish mumkin.

## API Integration

Server backend API bilan `BASE_URL` orqali integratsiya qiladi:

**Endpoint:** `GET /api/transaction/tag/{tag}/`

Transaction ma'lumotlarini olish uchun ishlatiladi.

## Development

### Code formatting

```bash
make fmt
```

### Linting

```bash
make lint
```

### Dependencies update

```bash
make deps
```

## Production Deployment

1. Build docker image:
```bash
docker build -t ocpp-server .
```

2. Run with docker-compose:
```yaml
version: '3.8'
services:
  ocpp:
    image: ocpp-server
    ports:
      - "10800:10800"
    environment:
      - BASE_URL=http://backend:8000
      - ADDR=:10800
    depends_on:
      - redis
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
```

## Monitoring

Server zap logger orqali barcha eventlarni log qiladi:

- Info level: Normal operatsiyalar
- Error level: Xatolar va exception'lar

## Contributing

1. Fork qiling
2. Feature branch yarating (`git checkout -b feature/amazing`)
3. Commit qiling (`git commit -m 'Add amazing feature'`)
4. Push qiling (`git push origin feature/amazing`)
5. Pull Request oching

## License

MIT License

## Support

Savollar uchun: [Issues](https://github.com/JscorpTech/ocpp/issues)
