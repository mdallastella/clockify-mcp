package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/server"

	"clockify-mcp/internal/clockify"
	"clockify-mcp/internal/tools"
)

func main() {
	// --- Required environment variables ---
	apiKey := os.Getenv("CLOCKIFY_API_KEY")
	if apiKey == "" {
		log.Fatal("CLOCKIFY_API_KEY environment variable is required")
	}

	mcpToken := os.Getenv("MCP_TOKEN")
	if mcpToken == "" {
		log.Fatal("MCP_TOKEN environment variable is required")
	}

	// --- Optional environment variables ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// --- Build Clockify client ---
	clockifyClient := clockify.NewClient(apiKey)

	// --- Build MCP server ---
	mcpServer := server.NewMCPServer(
		"clockify-mcp",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Register all domain tools
	tools.RegisterWorkspaceTools(mcpServer, clockifyClient)
	tools.RegisterTimeEntryTools(mcpServer, clockifyClient)
	tools.RegisterProjectTools(mcpServer, clockifyClient)
	tools.RegisterTaskTools(mcpServer, clockifyClient)
	tools.RegisterTagTools(mcpServer, clockifyClient)

	// --- Build Streamable HTTP transport ---
	httpServer := server.NewStreamableHTTPServer(mcpServer,
		server.WithEndpointPath("/mcp"),
	)

	// --- MCP_TOKEN auth middleware ---
	handler := authMiddleware(mcpToken, httpServer)

	addr := ":" + port
	log.Printf("Clockify MCP server listening on %s/mcp", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// authMiddleware enforces Bearer token authentication on all incoming requests.
// Requests missing or presenting an incorrect Authorization: Bearer <MCP_TOKEN>
// header are rejected with HTTP 401.
func authMiddleware(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			http.Error(w, "invalid Authorization header format, expected 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		providedToken := strings.TrimPrefix(authHeader, prefix)
		if providedToken != token {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

