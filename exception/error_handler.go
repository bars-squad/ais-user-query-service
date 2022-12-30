package exception

import (
	"encoding/json"
	"fmt"

	"github.com/bars-squad/ais-user-query-service/responses"
)

// InternalError wrap the error that contains error, and translate to error interface with informatif description.
func InternalError(r responses.Responses) error {
	if r.ErrorProperty() == nil {
		return nil
	}

	errMessage := &responses.ResponsesImpl{
		Error:   r.ErrorProperty(),
		Code:    r.CodeProperty(),
		Data:    r.DataProperty(),
		Message: r.MessageProperty(),
		Status:  r.StatusProperty(),
	}
	/* 	errMessage := map[string]interface{}{
		"code":    r.CodeProperty(),
		"data":    r.DataProperty(),
		"message": r.MessageProperty(),
		"status":  r.StatusProperty(),
	} */

	errMessageBuff, _ := json.Marshal(errMessage)
	return fmt.Errorf(string(errMessageBuff))
}
