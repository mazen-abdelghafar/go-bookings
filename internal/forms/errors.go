package forms

type errors map[string][]string

// adds an error msg for a given form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// returns the first error msg for a given field
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}