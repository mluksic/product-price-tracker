package storage

type Storer interface {
	Get()
	Insert()
	Modifier
}

type Modifier interface {
	Update()
	Delete()
}
