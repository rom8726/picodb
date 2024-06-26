package wal

import (
	"picodb/internal/database/compute"
	"picodb/internal/tools"
)

type LogData struct {
	LSN       int64
	CommandID compute.CommandID
	Arguments []string
}

type Log struct {
	data         LogData
	writePromise tools.Promise[error]
}

func NewLog(lsn int64, commandID compute.CommandID, args []string) Log {
	return Log{
		data: LogData{
			LSN:       lsn,
			CommandID: commandID,
			Arguments: args,
		},
		writePromise: tools.NewPromise[error](),
	}
}

func (l *Log) Data() LogData {
	return l.data
}

func (l *Log) LSN() int64 {
	return l.data.LSN
}

func (l *Log) CommandID() int {
	return l.data.CommandID.Int()
}

func (l *Log) Arguments() []string {
	return l.data.Arguments
}

func (l *Log) SetResult(err error) {
	l.writePromise.Set(err)
}

func (l *Log) Result() tools.Future[error] {
	return l.writePromise.GetFuture()
}
