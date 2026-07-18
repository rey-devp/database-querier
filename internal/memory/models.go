package memory

// Input — sesuai Input Contract (generalis)
type Task struct {
	ID          string `json:"_id"`
	UserRequest string `json:"user_request"`
}

// Output — sesuai Output Contract (success)
type AgentResponse struct {
	AgentName string `json:"agent_name"` // fixed: "database_querier"
	Result    string `json:"result"`
}

// Output — sesuai Output Contract (error)
type ErrorResponse struct {
	AgentName string `json:"agent_name"` // fixed: "database_querier"
	Result    string `json:"result"`
}
