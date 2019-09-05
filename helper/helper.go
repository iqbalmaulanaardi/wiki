package helper

import (
	"encoding/json"
	"net/http"
	"strings"
)

func catch(err error) {
	if err != nil {
		panic(err)
	}
}
// respondwithError return error message
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondwithJSON(w, code, map[string]string{"message": msg})
}
// respondwithJSON write json response format
func RespondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
func StringToByte(str string) string {
	stringSlice := []string{str}
	stringByte := "\x00" + strings.Join(stringSlice, "\x20\x00") // x20 = space and x00 = null
	return stringByte
}