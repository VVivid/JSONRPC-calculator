package main

import (
	"fmt"
	"log"
)

// Calculator provides arithmetic operations
type Calculator struct{}

// CalculatorParams represents parameters for binary operations
type CalculatorParams struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

// LogParams represents parameters for log notification
type LogParams struct {
	Message string `json:"message"`
}

// Add performs addition
func (c *Calculator) Add(params CalculatorParams) (float64, error) {
	result := params.A + params.B
	log.Printf("Calculator: %f + %f = %f", params.A, params.B, result)
	return result, nil
}

// Subtract performs subtraction
func (c *Calculator) Subtract(params CalculatorParams) (float64, error) {
	result := params.A - params.B
	log.Printf("Calculator: %f - %f = %f", params.A, params.B, result)
	return result, nil
}

// Multiply performs multiplication
func (c *Calculator) Multiply(params CalculatorParams) (float64, error) {
	result := params.A * params.B
	log.Printf("Calculator: %f * %f = %f", params.A, params.B, result)
	return result, nil
}

// Divide performs division with error handling for divide by zero
func (c *Calculator) Divide(params CalculatorParams) (float64, error) {
	if params.B == 0 {
		return 0, &JSONRPCError{
			Code:    -32000, // Application error
			Message: "Division by zero",
			Data:    fmt.Sprintf("Cannot divide %f by zero", params.A),
		}
	}
	
	result := params.A / params.B
	log.Printf("Calculator: %f / %f = %f", params.A, params.B, result)
	return result, nil
}

// Log handles notification messages (no response)
func (c *Calculator) Log(params LogParams) {
	log.Printf("Calculator Log: %s", params.Message)
	// Note: This is a notification, so we don't return anything
}

// GetInfo returns information about the calculator (demonstrates method without params)
func (c *Calculator) GetInfo() (map[string]interface{}, error) {
	info := map[string]interface{}{
		"name":        "JSON-RPC Calculator",
		"version":     "1.0",
		"methods":     []string{"add", "subtract", "multiply", "divide"},
		"description": "A simple calculator implementing JSON-RPC 2.0",
	}
	
	log.Printf("Calculator: GetInfo called")
	return info, nil
}