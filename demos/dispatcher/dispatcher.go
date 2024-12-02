package dispatcher

type Job interface {
	Do() error
}

var JobQueue chan Job
var DispatcherInstance *Dispatcher

type worker struct {
	workerId    int
	workerQueue chan Job
	quit        chan bool
}

func (w *worker) Start() {
	go func() {
		for {
			DispatcherInstance.idleWorker <- w.workerId
			select {
			case job := <-w.workerQueue:
				job.Do()
			case <-w.quit:
				return
			}
		}
	}()
}
func (w worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

type Dispatcher struct {
	Name       string
	maxWorkers int
	idleWorker chan int
	workerList []worker
}

func New(maxWorkerNum, maxQueueNum int) *Dispatcher {
	JobQueue = make(chan Job, maxQueueNum)

	DispatcherInstance = &Dispatcher{
		idleWorker: make(chan int, maxWorkerNum),
		Name:       "Dispatcher001",
		maxWorkers: maxWorkerNum,
	}

	return DispatcherInstance
}

func (d *Dispatcher) Run() {
	d.workerList = make([]worker, d.maxWorkers)
	for i := 0; i < d.maxWorkers; i++ {
		d.workerList[i] = worker{
			workerId:    i,
			workerQueue: make(chan Job),
			quit:        make(chan bool),
		}

		d.workerList[i].Start()
	}
	go func() {
		for {
			select {
			case job := <-JobQueue:
				workerId := <-d.idleWorker
				d.workerList[workerId].workerQueue <- job
			}
		}
	}()
}

func (d *Dispatcher) Dispatch(job Job) {
	JobQueue <- job
}

func (d *Dispatcher) Close() {
	for _, w := range d.workerList {
		w.Stop()
	}
}
