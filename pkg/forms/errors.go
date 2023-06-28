package forms

// holds validation error message for forms, a "key" can have multiple values
type errMap map[string][]string

// Add key & value pair for validation to errMap
func (e errMap) Add(key, value string) {
	e[key] = append(e[key], value)
}

func (e errMap) Get(key string) string {
	if value, ok := e[key]; ok {
		return value[0]
	}

	return ""
}
