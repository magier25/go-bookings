package forms

type errors map[string][]string

// Add adds an error message to a given field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error message
func (e errors) Get(field string) string {
	es, ok := e[field]
	if !ok {
		return ""
	}

	return es[0]
}
