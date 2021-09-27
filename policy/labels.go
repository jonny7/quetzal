package policy

// Labels represent the required labels policy condition
type Labels struct {
	labels []string `yaml:"labels"`
}

// ForbiddenLabels represent any label that should exclude the webhook
// if present
type ForbiddenLabels struct {
	forbiddenLabels []string `yaml:"forbiddenLabels"`
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

func (fl ForbiddenLabels) conditionMet(event GitLabAdaptor) bool {
	if len(fl.forbiddenLabels) == 0 {
		return true
	}

	webhookLabels, err := event.Labels()
	if err != nil {
		return false
	}

	var valid = true
	for _, forbiddenPolicyLabel := range fl.forbiddenLabels {
		reject := existsInSlice(webhookLabels, forbiddenPolicyLabel)
		if reject {
			valid = false
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
