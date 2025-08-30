package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

// JSONRPCServer handles JSON-RPC requests
type JSONRPCServer struct {
	calculator *Calculator
}

// NewJSONRPCServer creates a new JSON-RPC server
func NewJSONRPCServer() *JSONRPCServer {
	return &JSONRPCServer{
		calculator: &Calculator{},
	}
}

// HandleRequest processes a JSON-RPC request and returns a response
func (s *JSONRPCServer) HandleRequest(data []byte) ([]byte, error) {
	log.Printf("Received request: %s", string(data))

	// Parse the incoming message
	message, err := ParseMessage(data)
	if err != nil {
		// Parse error - we can't know the ID, so use null
		errorResp := CreateErrorResponse(err.(*JSONRPCError), nil)
		return json.Marshal(errorResp)
	}

	// Handle batch vs single requests
	switch msg := message.(type) {
	case []interface{}:
		// Batch request
		return s.handleBatchRequest(msg)
	case JSONRPCRequest:
		// Single request
		response := s.handleSingleRequest(msg)
		return json.Marshal(response)
	case JSONRPCNotification:
		// Single notification - no response
		s.handleNotification(msg)
		return nil, nil // No response for notifications
	default:
		// This shouldn't happen if parsing worked correctly
		errorResp := CreateErrorResponse(&JSONRPCError{
			Code:    InvalidRequest,
			Message: "Invalid Request",
			Data:    "Unknown message type",
		}, nil)
		return json.Marshal(errorResp)
	}
}

// handleBatchRequest processes a batch of requests/notifications
func (s *JSONRPCServer) handleBatchRequest(messages []interface{}) ([]byte, error) {
	var responses []JSONRPCResponse

	for _, msg := range messages {
		switch m := msg.(type) {
		case JSONRPCRequest:
			// Request - add response to batch
			response := s.handleSingleRequest(m)
			responses = append(responses, response)
		case JSONRPCNotification:
			// Notification - handle but don't add to responses
			s.handleNotification(m)
		}
	}

	// If no responses (all were notifications), return empty
	if len(responses) == 0 {
		return nil, nil
	}

	return json.Marshal(responses)
}

// handleSingleRequest processes a single JSON-RPC request
func (s *JSONRPCServer) handleSingleRequest(req JSONRPCRequest) JSONRPCResponse {
	// Route the method call
	result, err := s.callMethod(req.Method, req.Params)
	if err != nil {
		// Check if it's already a JSON-RPC error
		if jsonrpcErr, ok := err.(*JSONRPCError); ok {
			return CreateErrorResponse(jsonrpcErr, req.ID)
		}

		// Convert regular error to JSON-RPC error
		jsonrpcErr := &JSONRPCError{
			Code:    InternalError,
			Message: "Internal error",
			Data:    err.Error(),
		}
		return CreateErrorResponse(jsonrpcErr, req.ID)
	}

	return CreateSuccessResponse(result, req.ID)
}

// handleNotification processes a notification (no response)
func (s *JSONRPCServer) handleNotification(notif JSONRPCNotification) {
	log.Printf("Handling notification: %s", notif.Method)

	// Call method but ignore any result/error since it's a notification
	_, err := s.callMethod(notif.Method, notif.Params)
	if err != nil {
		log.Printf("Notification error (ignored): %v", err)
	}
}

// callMethod dispatches method calls to the calculator
func (s *JSONRPCServer) callMethod(method string, params interface{}) (interface{}, error) {
	switch method {
	case "add":
		return s.callCalculatorMethod("Add", params)
	case "subtract":
		return s.callCalculatorMethod("Subtract", params)
	case "multiply":
		return s.callCalculatorMethod("Multiply", params)
	case "divide":
		return s.callCalculatorMethod("Divide", params)
	case "getInfo":
		return s.calculator.GetInfo()
	case "log":
		return s.callNotificationMethod("Log", params)
	default:
		return nil, &JSONRPCError{
			Code:    MethodNotFound,
			Message: "Method not found",
			Data:    fmt.Sprintf("Method '%s' is not available", method),
		}
	}
}

// callCalculatorMethod calls a calculator method that expects CalculatorParams
func (s *JSONRPCServer) callCalculatorMethod(methodName string, params interface{}) (interface{}, error) {
	// Parse parameters
	var calcParams CalculatorParams
	if params != nil {
		paramBytes, err := json.Marshal(params)
		if err != nil {
			return nil, &JSONRPCError{
				Code:    InvalidParams,
				Message: "Invalid params",
				Data:    "Cannot marshal parameters",
			}
		}

		if err := json.Unmarshal(paramBytes, &calcParams); err != nil {
			return nil, &JSONRPCError{
				Code:    InvalidParams,
				Message: "Invalid params",
				Data:    "Expected parameters: {\"a\": number, \"b\": number}",
			}
		}
	} else {
		return nil, &JSONRPCError{
			Code:    InvalidParams,
			Message: "Invalid params",
			Data:    "Parameters required: {\"a\": number, \"b\": number}",
		}
	}

	// Use reflection to call the method
	calcValue := reflect.ValueOf(s.calculator)
	method := calcValue.MethodByName(methodName)
	if !method.IsValid() {
		return nil, &JSONRPCError{
			Code:    InternalError,
			Message: "Internal error",
			Data:    fmt.Sprintf("Method %s not found on calculator", methodName),
		}
	}

	// Call the method
	results := method.Call([]reflect.Value{reflect.ValueOf(calcParams)})

	// Handle results (expecting result, error pattern)
	if len(results) != 2 {
		return nil, &JSONRPCError{
			Code:    InternalError,
			Message: "Internal error",
			Data:    "Unexpected return value count",
		}
	}

	// Check for error
	if !results[1].IsNil() {
		err := results[1].Interface().(error)
		return nil, err
	}

	// Return the result
	return results[0].Interface(), nil
}

// callNotificationMethod calls a method for notifications (no return value expected)
func (s *JSONRPCServer) callNotificationMethod(methodName string, params interface{}) (interface{}, error) {
	switch methodName {
	case "Log":
		var logParams LogParams
		if params != nil {
			paramBytes, err := json.Marshal(params)
			if err != nil {
				return nil, &JSONRPCError{
					Code:    InvalidParams,
					Message: "Invalid params",
					Data:    "Cannot marshal parameters",
				}
			}

			if err := json.Unmarshal(paramBytes, &logParams); err != nil {
				return nil, &JSONRPCError{
					Code:    InvalidParams,
					Message: "Invalid params",
					Data:    "Expected parameters: {\"message\": string}",
				}
			}
		} else {
			return nil, &JSONRPCError{
				Code:    InvalidParams,
				Message: "Invalid params",
				Data:    "Parameters required: {\"message\": string}",
			}
		}

		s.calculator.Log(logParams)
		return nil, nil // No return value for notifications
	}

	return nil, &JSONRPCError{
		Code:    MethodNotFound,
		Message: "Method not found",
		Data:    fmt.Sprintf("Notification method '%s' is not available", methodName),
	}
}

