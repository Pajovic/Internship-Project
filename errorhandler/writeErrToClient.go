package errorhandler

import "net/http"

// WriteErrToClient is used to return err to client
func WriteErrToClient(w http.ResponseWriter, err error) {
	errMsg, code := GetErrorMsg(err)
	w.WriteHeader(code)
	w.Write([]byte(errMsg))
}
