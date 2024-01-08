package compute

type Query struct {
	commandID CommandID
	arguments []string
}

func NewQuery(commandID CommandID, arguments []string) Query {
	return Query{
		commandID: commandID,
		arguments: arguments,
	}
}

func (c *Query) CommandID() CommandID {
	return c.commandID
}

func (c *Query) Arguments() []string {
	return c.arguments
}
