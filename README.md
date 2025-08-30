# JSON-RPC Calculator

A simple calculator server implementing JSON-RPC 2.0 specification.

## Usage

```bash
go run .
```

Server runs on `http://localhost:8090`

## Examples

**Request:**
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"add","params":{"a":15,"b":25},"id":1}' \
  http://localhost:8090/
```

**Response:**
```json
{"jsonrpc":"2.0","result":40,"id":1}
```

**Notification (no response):**
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"log","params":{"message":"Hello"}}' \
  http://localhost:8090/
```

## Methods

- `add` - Addition
- `subtract` - Subtraction  
- `multiply` - Multiplication
- `divide` - Division
- `log` - Log message (notification only)