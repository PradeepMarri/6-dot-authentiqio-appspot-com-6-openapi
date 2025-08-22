package models

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
)

type Tool struct {
	Definition mcp.Tool
	Handler    func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

// Claims represents the Claims schema from the OpenAPI specification
type Claims struct {
	Scope string `json:"scope"` // claim scope
	Sub string `json:"sub"` // UUID
	TypeField string `json:"type,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// Error represents the Error schema from the OpenAPI specification
type Error struct {
	Detail string `json:"detail,omitempty"`
	ErrorField int `json:"error"`
	Title string `json:"title,omitempty"`
	TypeField string `json:"type,omitempty"` // unique uri for this error
}

// PushToken represents the PushToken schema from the OpenAPI specification
type PushToken struct {
	Exp int `json:"exp,omitempty"`
	Iat int `json:"iat,omitempty"`
	Iss string `json:"iss"` // issuer (URI)
	Nbf int `json:"nbf,omitempty"`
	Sub string `json:"sub"` // UUID and public signing key
	Aud string `json:"aud"` // audience (URI)
}

// AuthentiqID represents the AuthentiqID schema from the OpenAPI specification
type AuthentiqID struct {
	Devtoken string `json:"devtoken,omitempty"` // device token for push messages
	Sub string `json:"sub"` // UUID and public signing key
}
