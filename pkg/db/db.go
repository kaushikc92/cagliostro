package db

const DatabaseUrl string = "mongodb://mongo-0.mongo,mongo-1.mongo:27017"

type ErrRecordDoesntExist struct {
	Err error
}

func (err *ErrRecordDoesntExist) Error() string {
	return err.Err.Error()
}
