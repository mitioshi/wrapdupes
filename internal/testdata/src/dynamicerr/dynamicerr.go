// Package dynamicerr illustrates a method call on a local variable in the return statement
// This is a valid case
package dynamicerr

func getMessages() ([]string, error) {
	return []string{}, nil
}
func getValidMessages() (*[]string, string) {
	messages, err := getMessages()
	if err != nil {
		return nil, err.Error()
	}
	return &messages, ""
}
