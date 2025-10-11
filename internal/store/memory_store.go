package store

type MemoryStore interface {
	Add()
	Delete()
	Update()
	Data()
}
