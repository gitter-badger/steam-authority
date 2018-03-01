package web

import (
	"net/http"

	"github.com/steam-authority/steam-authority/queue"
)

func QueuesHandler(w http.ResponseWriter, r *http.Request) {

	queues, err := queue.GetQeueus()
	if err != nil {
		returnErrorTemplate(w, r, 500, err.Error())
		return
	}

	// Template
	template := queueTemplate{}
	template.Fill(r)
	template.Queues = queues

	returnTemplate(w, r, "queues", template)
	return
}

type queueTemplate struct {
	GlobalTemplate
	Queues []queue.Queue
}
