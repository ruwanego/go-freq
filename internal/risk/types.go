package risk

type Decision string

const (
	DecisionApprove Decision = "APPROVE"
	DecisionReject  Decision = "REJECT"
)

type Result struct {
	Decision Decision
	Reason   string
}
