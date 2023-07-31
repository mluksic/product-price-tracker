package storage

type Storer interface {
	Connect()
	Get()
	Insert()
	Modifier
}

type Modifier interface {
	Update()
	Delete()
}
