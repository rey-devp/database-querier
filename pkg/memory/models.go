package memory

// Request dari Orchestrator (sesuai Guide Banana Dev Team)
type OrchestratorRequest struct {
	TaskID    string          `json:"task_id"`
	AgentType string          `json:"agent_type"`
	Payload   RequestPayload  `json:"payload"`
	Metadata  RequestMetadata `json:"metadata"`
}

type RequestPayload struct {
	URL     string `json:"url"`
	Keyword string `json:"keyword"`
	RawText string `json:"raw_text"` // Input utama kita
}

type RequestMetadata struct {
	Sender    string `json:"sender"`
	Timestamp int64  `json:"timestamp"`
}

// Internal task (untuk in-memory store)
type Task struct {
	ID          string
	UserRequest string
}

// Response sukses (HTTP 200)
type SuccessResponse struct {
	Status  string       `json:"status"` // "success"
	TaskID  string       `json:"task_id"`
	Data    ResponseData `json:"data"`
	Message string       `json:"message"`
}

type ResponseData struct {
	Result  string  `json:"result"`
	FileURL *string `json:"file_url"` // null untuk agent teks
}

// Response error (HTTP 400/500)
type ErrorResponse struct {
	Status  string      `json:"status"` // "error"
	TaskID  string      `json:"task_id"`
	Data    interface{} `json:"data"`   // null
	Message string      `json:"message"`
}
