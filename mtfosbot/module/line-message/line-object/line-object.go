package lineobj

// EventObject -
type EventObject struct {
	Source     *SourceObject          `json:"source" cc:"source"`
	Type       string                 `json:"type" cc:"type"`
	Timestamp  int64                  `json:"timestamp" cc:"timestamp"`
	ReplyToken string                 `json:"replyToken" cc:"replyToken"`
	Message    map[string]interface{} `json:"message" cc:"message"`
}

// SourceObject -
type SourceObject struct {
	Type    string `json:"type" cc:"type"`
	UserID  string `json:"userId" cc:"userId"`
	GroupID string `json:"groupId" cc:"groupId"`
}
