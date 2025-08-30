package main

import (
	"encoding/json"
	"fmt"
)

// JSONRPCVersion represents the JSON-RPC protocol version
const JSONRPCVersion = "2.0"

// JSONRPCMessage represents any JSON-RPC message that can be sent over the wire
type JSONRPCMessage interface {
	GetJSONRPC() string
}

// JSONRPCRequest represents a JSON-RPC request
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      interface{} `json:"id"`
}

func (r JSONRPCRequest) GetJSONRPC() string {
	return r.JSONRPC
}

// JSONRPCNotification represents a JSON-RPC notification (no ID, no response expected)
type JSONRPCNotification struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

func (n JSONRPCNotification) GetJSONRPC() string {
	return n.JSONRPC
}

// JSONRPCResponse represents a JSON-RPC response
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

func (r JSONRPCResponse) GetJSONRPC() string {
	return r.JSONRPC
}

// JSONRPCError represents a JSON-RPC error
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e JSONRPCError) Error() string {
	return fmt.Sprintf("JSON-RPC Error %d: %s", e.Code, e.Message)
}

// Standard JSON-RPC error codes
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

// ParseMessage attempts to parse a JSON-RPC message and determine its type
func ParseMessage(data []byte) (interface{}, error) {
	// First, try to determine if it's a batch request (array)
	if len(data) > 0 && data[0] == '[' {
		var batch []json.RawMessage
		if err := json.Unmarshal(data, &batch); err != nil {
			return nil, &JSONRPCError{
				Code:    ParseError,
				Message: "Parse error",
				Data:    err.Error(),
			}
		}
		
		var messages []interface{}
		for _, raw := range batch {
			msg, err := ParseSingleMessage(raw)
			if err != nil {
				return nil, err
			}
			messages = append(messages, msg)
		}
		return messages, nil
	}
	
	// Single message
	return ParseSingleMessage(data)
}

// ParseSingleMessage parses a single JSON-RPC message
func ParseSingleMessage(data []byte) (interface{}, error) {
	// Parse once into a raw message that preserves ID field
	var raw struct {
		JSONRPC string          `json:"jsonrpc"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params,omitempty"`
		ID      *json.RawMessage `json:"id,omitempty"` // Pointer to detect presence
	}
	
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, &JSONRPCError{
			Code:    ParseError,
			Message: "Parse error",
			Data:    err.Error(),
		}
	}
	
	// Validate common fields
	if raw.JSONRPC != JSONRPCVersion {
		return nil, &JSONRPCError{
			Code:    InvalidRequest,
			Message: "Invalid Request",
			Data:    "jsonrpc field must be '2.0'",
		}
	}
	
	if raw.Method == "" {
		return nil, &JSONRPCError{
			Code:    InvalidRequest,
			Message: "Invalid Request",
			Data:    "method field is required",
		}
	}
	
	// Parse params once (same for both request and notification)
	var params interface{}
	if raw.Params != nil {
		if err := json.Unmarshal(raw.Params, &params); err != nil {
			return nil, &JSONRPCError{
				Code:    InvalidRequest,
				Message: "Invalid Request",
				Data:    "Invalid params field",
			}
		}
	}
	
	// Only difference: check ID at the end to determine type
	if raw.ID != nil {
		// It's a request (has ID, expects response)
		var id interface{}
		if err := json.Unmarshal(*raw.ID, &id); err != nil {
			return nil, &JSONRPCError{
				Code:    InvalidRequest,
				Message: "Invalid Request",
				Data:    "Invalid ID field",
			}
		}
		
		return JSONRPCRequest{
			JSONRPC: raw.JSONRPC,
			Method:  raw.Method,
			Params:  params,
			ID:      id,
		}, nil
	}
	
	// No ID = notification (no response expected)
	return JSONRPCNotification{
		JSONRPC: raw.JSONRPC,
		Method:  raw.Method,
		Params:  params,
	}, nil
}

// CreateSuccessResponse creates a successful JSON-RPC response
func CreateSuccessResponse(result interface{}, id interface{}) JSONRPCResponse {
	return JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		Result:  result,
		ID:      id,
	}
}

// CreateErrorResponse creates an error JSON-RPC response
func CreateErrorResponse(err *JSONRPCError, id interface{}) JSONRPCResponse {
	return JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		Error:   err,
		ID:      id,
	}
}