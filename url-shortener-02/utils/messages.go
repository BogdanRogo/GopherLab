package utils

import "encoding/json"

// messages contains user returned messages
// user response structure is a generic structure
// when the JSON marshal is made, if a key from that value is nil
// the structure key is not marshaled when using `omitempty`

type userResponse struct {
	URL   string `json:"url,omitempty"`
	Error string `json:"error,omitempty"`
}

func (um *userResponse) newUserResponse() ([]byte, error) {
	message, err := json.Marshal(um)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// ReturnError returns a JSON marsheled user response
func ReturnError(errorMessage string) string {
	userMessage := userResponse{}
	userMessage.Error = errorMessage
	message, _ := userMessage.newUserResponse()
	// errors are important handle them !
	return string(message)
}

// ReturnURL returns a JSON marsheled user response
func ReturnURL(URL string) []byte {
	userMessage := userResponse{}
	userMessage.URL = URL
	message, _ := userMessage.newUserResponse()
	// errors are important handle them !
	return message
}
