package validator

type Validator struct {
	Erorrs map[string]string
}

func New() *Validator {
	return &Validator{Erorrs: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Erorrs) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Erorrs[key]; !exists {
		v.Erorrs[key] = message
	}
}

func (v *Validator) Check(expression bool, key string, message string) {
	if !expression {
		v.AddError(key, message)
	}
}

func (v *Validator) In(value string, keys []string) bool {
	for _, key := range(keys) {
		if key == value {
			return true
		}
	}
	return false
}