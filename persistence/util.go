package persistence

type PersistenceExtension interface {
	Extend(name string, val interface{})
}
