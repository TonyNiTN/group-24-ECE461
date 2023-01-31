package error

import "fmt"

// GraphQLError represents a custom error class for GraphQL errors
type GraphQLError struct {
	Message string
	Path    string
}

func (e *GraphQLError) Error() string {
	return fmt.Sprintf("GraphQL error: %s (path: %s)", e.Message, e.Path)
}

// NewGraphQLError creates a new instance of GraphQLError
func NewGraphQLError(message, path string) *GraphQLError {
	return &GraphQLError{
		Message: message,
		Path:    path,
	}
}
