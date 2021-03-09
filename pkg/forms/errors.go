package forms

// 407
//
type errors map[string][]string

// struktur func : func (biasnaya struct receiver) nama fungsi(input?) output{}
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
