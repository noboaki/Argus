package domain

type Labels map[string]string

func (l *Labels) Merge(other Labels) {
	if *l == nil {
		*l = make(Labels)
	}
	for k, v := range other {
		(*l)[k] = v
	}
}
