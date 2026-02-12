package pages

// ProposalFormData stores user-submitted proposal information.
type ProposalFormData struct {
	Name        string
	Email       string
	Company     string
	ProjectType string
	Budget      string
	Timeline    string
	Details     string
}

// ProposalPageData contains UI state for the proposal page.
type ProposalPageData struct {
	Form    ProposalFormData
	Error   string
	Success bool
}
