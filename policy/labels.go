package policy

type Labels struct {
	labels []string `yaml:"labels"`
}

func (l Labels) conditionMet(event GitLabAdaptor) bool {
	if len(l.labels) == 0 {
		return true
	}
	webhookLabels, err := event.Labels()
	if err != nil {
		return false
	}
	var valid bool
	for _, policyLabel := range l.labels {
		valid = existsInSlice(webhookLabels, policyLabel)
		if !valid {
			break
		}
	}
	return valid
}

// existsInSlice checks for the presence of a string in a slice
func existsInSlice(a []string, b string) bool {
	for _, item := range a {
		if item == b {
			return true
		}
	}
	return false
}
