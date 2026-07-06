package memory

// Input — sesuai Input Contract
type Task struct {
	ID            string                 `json:"_id"`
	UserRequest   string                 `json:"user_request"`
	SharedContext map[string]interface{} `json:"shared_context"`
}

// Output — sesuai Output Contract (success)
type AgentResponse struct {
	AgentName     string        `json:"agent_name"` // fixed: "database_querier"
	MemoryPayload MemoryPayload `json:"memory_payload"`
}

type MemoryPayload struct {
	Result QueryResult `json:"database_querier_result"`
}

type QueryResult struct {
	Collection string      `json:"collection"`
	Operation  string      `json:"operation"`
	Filter     interface{} `json:"filter"`
	Projection interface{} `json:"projection,omitempty"`
	Documents  interface{} `json:"documents"`
	Total      int         `json:"total"`
}

// Output — sesuai Output Contract (error)
type ErrorResponse struct {
	AgentName string `json:"agent_name"` // fixed: "database_querier"
	Status    string `json:"status"`     // "failed"
	Message   string `json:"message"`
}
