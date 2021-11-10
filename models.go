package main

// Metadata of the retention rule
type RuleMetadata struct {
	// UUID of rule
	ID int `json:"id"`

	// Priority of rule when doing calculating
	Priority int `json:"priority"`

	// Disabled rule
	Disabled bool `json:"disabled"`

	// Action of the rule performs
	// "retain"
	Action string `json:"action" valid:"Required"`

	// Template ID
	Template string `json:"template" valid:"Required"`

	// The parameters of this rule
	Parameters Parameters `json:"params" valid:"Required"`

	// Selector attached to the rule for filtering tags
	TagSelectors []*Selector `json:"tag_selectors" valid:"Required"`

	// Selector attached to the rule for filtering scope (e.g: repositories or namespaces)
	ScopeSelectors map[string][]*Selector `json:"scope_selectors" valid:"Required"`
}

// Parameters of rule, indexed by the key
type Parameters map[string]Parameter

// Parameter of rule
type Parameter interface{}

// Scope definition
type Scope struct {
	// Scope level declaration
	// 'system', 'project' and 'repository'
	Level string `json:"level" valid:"Required;Match(/^(project)$/)"`

	// The reference identity for the specified level
	// 0 for 'system', project ID for 'project' and repo ID for 'repository'
	Reference int64 `json:"ref" valid:"Required"`
}

// Trigger of the policy
type Trigger struct {
	// Const string to declare the trigger type
	// 'Schedule'
	Kind string `json:"kind" valid:"Required"`

	// Settings for the specified trigger
	// '[cron]="* 22 11 * * *"' for the 'Schedule'
	Settings map[string]interface{} `json:"settings" valid:"Required"`
}

// Metadata of policy
type Metadata struct {
	// ID of the policy
	ID int64 `json:"id"`

	// Algorithm applied to the rules
	// "OR" / "AND"
	Algorithm string `json:"algorithm" valid:"Required;Match(or)"`

	// Rule collection
	Rules []RuleMetadata `json:"rules"`

	// Trigger about how to launch the policy
	Trigger *Trigger `json:"trigger" valid:"Required"`

	// Which scope the policy will be applied to
	Scope *Scope `json:"scope" valid:"Required"`
}

// Selector to narrow down the list
type Selector struct {
	// Kind of the selector
	// "doublestar" or "label"
	Kind string `json:"kind" valid:"Required;Match(doublestar)"`

	// Decorated the selector
	// for "doublestar" : "matching" and "excluding"
	// for "label" : "with" and "without"
	Decoration string `json:"decoration" valid:"Required"`

	// Param for the selector
	Pattern string `json:"pattern" valid:"Required"`

	// Extras for other settings
	Extras string `json:"extras"`
}
