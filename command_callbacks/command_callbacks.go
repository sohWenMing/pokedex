package commandcallbacks

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var callBackMap map[string]func() = map[string]func(){
	"exit": exitCallBackfunc,
}

func exitCallBackfunc() {
	return
}
