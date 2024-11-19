package ground

// Report Ground Error
type MsgVerify struct {
	Force bool
}

// tick message
type MsgTick struct{}

// tick message
type MsgTick_max struct{}

type MsgGroundOk struct{}
type MsgGroundErr struct{}
type MsgRequestStatus struct{}
type MsgStatus struct{}
type MsgStarted struct{}
