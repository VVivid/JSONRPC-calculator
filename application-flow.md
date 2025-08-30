# JSON-RPC Calculator Application Flow

## Overview Diagram

```mermaid
flowchart TD
    Start([Client Request]) --> Receive[Receive JSON-RPC Request]
    Receive --> Parse[Parse JSON Request]
    Parse --> Validate{Valid JSON-RPC?}
    Validate -->|No| Error[Return Error Response]
    Validate -->|Yes| Identify[Identify Request Type]
    
    Identify --> Type{Request Type?}
    Type -->|Method Call| Match[Match Method Name]
    Type -->|Notification| Notify[Process Notification]
    
    Match --> Method{Method Exists?}
    Method -->|No| MethodError[Method Not Found Error]
    Method -->|Yes| Extract[Extract Parameters]
    
    Extract --> Execute[Execute Calculator Function]
    Execute --> Result[Generate Result]
    Result --> Response[Build JSON-RPC Response]
    
    MethodError --> Response
    Error --> Send[Send Response to Client]
    Response --> Send
    Notify --> End([End])
    Send --> End
    
    style Start fill:#e1f5fe
    style End fill:#e1f5fe
    style Execute fill:#c8e6c9
    style Error fill:#ffcdd2
    style MethodError fill:#ffcdd2
```

## Detailed Flow Steps

### 1. **Receive Request**
The server receives an incoming JSON-RPC request from the client containing method name and parameters.

### 2. **Parse Request**
Parse the JSON payload to extract the request structure and validate it against JSON-RPC 2.0 specification.

### 3. **Identify Request Type**
Determine if the request is:
- A method call (expects a response)
- A notification (no response expected)
- A batch request (multiple operations)

### 4. **Match Method**
Map the requested method name to the corresponding calculator function:
- `add` → Addition function
- `subtract` → Subtraction function
- `multiply` → Multiplication function
- `divide` → Division function

### 5. **Execute Function**
Call the matched calculator function with the provided parameters and handle any computational errors (e.g., division by zero).

### 6. **Send Response**
Format and return the JSON-RPC response with either:
- **Success**: Result with matching request ID
- **Error**: Error object with code and message