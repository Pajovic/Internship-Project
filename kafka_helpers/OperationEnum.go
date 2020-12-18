package kafka_helpers

type OperationEnum int

const (
	Created OperationEnum = iota
	Updated
	Deleted
)

func OperationEnumString(e OperationEnum) string {
	switch e {
	case Created:
		return "CREATED"
	case Updated:
		return "UPDATED"
	case Deleted:
		return "DELETED"
	default:
		return ""
	}
}

func (e OperationEnum) FromString(s string) {

}
