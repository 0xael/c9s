package ui

type ToastMsg struct {
	Text string
}

type ConfirmAction int

const (
	ConfirmStop ConfirmAction = iota
	ConfirmRemove
)

type ConfirmResultMsg struct {
	Action    ConfirmAction
	ID        string
	Confirmed bool
}

type DispatchResultMsg struct {
	ID  string
	Err error
}
