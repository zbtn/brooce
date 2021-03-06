package web

import (
	"net/http"

	"brooce/heartbeat"
	"brooce/listing"
	"brooce/task"

	"github.com/go-redis/redis"
)

type mainpageOutputType struct {
	Queues         map[string]*listing.QueueInfoType
	RunningJobs    []*task.Task
	RunningWorkers []*heartbeat.HeartbeatType
	TotalThreads   int
}

func mainpageHandler(req *http.Request, rep *httpReply) (err error) {
	output := &mainpageOutputType{}

	output.RunningJobs, err = listing.RunningJobs(true)
	if err != nil {
		return
	}
	output.RunningWorkers, err = listing.RunningWorkers()
	if err != nil {
		return
	}
	output.Queues, err = listing.Queues(false)
	if err != nil {
		return
	}

	for _, worker := range output.RunningWorkers {
		output.TotalThreads += len(worker.Threads)
	}

	err = templates.ExecuteTemplate(rep, "mainpage", output)
	return
}

type listQueueType struct {
	QueueName     string
	Pending       int64
	Running       int
	Done          int64
	Failed        int64
	Delayed       int64
	pendingResult *redis.IntCmd
	runningResult *redis.StringSliceCmd
	doneResult    *redis.IntCmd
	failedResult  *redis.IntCmd
	delayedResult *redis.IntCmd
}
