package validator

type Validator struct {
	Errors map[string]string
}

func New() Validator {
	return Validator{Errors: make(map[string]string)}
}

func (v Validator) Check(ok bool, key, message string) {
	_, exists := v.Errors[key]
	if !ok && !exists {
		v.Errors[key] = message
	}
}

func (v Validator) IsValid() bool {
	return len(v.Errors) == 0
}
