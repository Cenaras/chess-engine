package uci

type UCICommand string

const (
	UCI        UCICommand = "uci"
	IsReady    UCICommand = "isready"
	UCINewGame UCICommand = "ucinewgame"
	Position   UCICommand = "position"
	Go         UCICommand = "go"
	SetOptions UCICommand = "setoptions"
)
