package internal

type Worker struct {
	ID int
}

func SpawnWorker(id int) *Worker {
	return &Worker{ID: id}
}
