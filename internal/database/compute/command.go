package compute

type CommandID int

const (
	UnknownCommandID CommandID = iota
	SetCommandID
	GetCommandID
	DelCommandID
)

var (
	UnknownCommand = "UNKNOWN"
	SetCommand     = "SET"
	GetCommand     = "GET"
	DelCommand     = "DEL"
)

var commandNamesToID = map[string]CommandID{
	UnknownCommand: UnknownCommandID,
	SetCommand:     SetCommandID,
	GetCommand:     GetCommandID,
	DelCommand:     DelCommandID,
}

func (c CommandID) Int() int {
	return int(c)
}

func CommandNameToCommandID(command string) CommandID {
	id, found := commandNamesToID[command]
	if !found {
		return UnknownCommandID
	}

	return id
}
