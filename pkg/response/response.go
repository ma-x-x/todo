// Package response provides standard API response structures
package response

// Response represents a standard API response format
// @Description Standard API response structure
type Response struct {
    // HTTP status code
    Code    int         `json:"code" example:"200"`
    
    // Response message
    Message string      `json:"message" example:"Success"`
    
    // Response payload
    Data    interface{} `json:"data,omitempty"`
    
    // Request trace ID for debugging
    TraceID string      `json:"traceId,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// NewResponse creates a new Response instance
func NewResponse(code int, message string, data interface{}) *Response {
    return &Response{
        Code:    code,
        Message: message,
        Data:    data,
    }
}

// WithTraceID adds a trace ID to the response
func (r *Response) WithTraceID(traceID string) *Response {
    r.TraceID = traceID
    return r
}

// Success creates a success response
func Success(data interface{}) *Response {
    return NewResponse(200, "Success", data)
}

// Error creates an error response
func Error(code int, message string) *Response {
    return NewResponse(code, message, nil)
} 