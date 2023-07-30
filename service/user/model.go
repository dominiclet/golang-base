package user

type AccountType int

const (
	TrialAccount AccountType = iota
	BasicAccount
	TestAccount // Account with no expiry (remove after beta test)
)

func (a AccountType) String() string {
	switch a {
	case TrialAccount:
		return "Free trial"
	case BasicAccount:
		return "Basic"
	case TestAccount:
		return "Test"
	}
	return "Unknown"
}

const (
	TRIAL_DURATION_DAYS = 14
)

const (
	EmailVerificationTokenLength     = 16
	verificationEmailDisableDuration = 2 // minutes
)
