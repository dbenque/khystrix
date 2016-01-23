package main

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

func startHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	baseStr := r.URL.Query().Get("base")
	floorStr := r.URL.Query().Get("floor")

	if (len(name) == 0) || (len(baseStr) == 0) || (len(floorStr) == 0) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte)("Missing or empty parameter\n"))
		return
	}

	base, err := strconv.Atoi(baseStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte)("base parameter not ' int' "))
		return
	}

	floor, err := strconv.Atoi(floorStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte)("floor parameter not ' int' "))
		return
	}

	go LoopOverCmd(name, base, floor)

	w.WriteHeader(http.StatusOK)
	w.Write(([]byte)("New command loop started\n"))
}

// hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
// 	Timeout:               100,
// 	MaxConcurrentRequests: 100,
// 	ErrorPercentThreshold: 50,
// })

func configureHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	timeoutStr := r.URL.Query().Get("timeout")
	maxConcurrentRequestsStr := r.URL.Query().Get("maxConcurrentRequests")
	errorPercentThresholdStr := r.URL.Query().Get("errorPercentThreshold")

	if (len(name) == 0) || (len(timeoutStr) == 0) || (len(maxConcurrentRequestsStr) == 0) || (len(errorPercentThresholdStr) == 0) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte)("Missing or empty parameter\n"))
		return
	}

	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte)("timeout parameter not ' int' "))
		return
	}

	maxConcurrentRequests, err := strconv.Atoi(maxConcurrentRequestsStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte)("maxConcurrentRequests parameter not ' int' "))
		return
	}

	errorPercentThreshold, err := strconv.Atoi(errorPercentThresholdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte)("errorPercentThreshold parameter not ' int' "))
		return
	}

	hystrix.ConfigureCommand(name, hystrix.CommandConfig{
		Timeout:               timeout,
		MaxConcurrentRequests: maxConcurrentRequests,
		ErrorPercentThreshold: errorPercentThreshold,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(([]byte)("Configuration done.\n"))

}

func closeHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte)("Missing or empty parameter\n"))
		return
	}

	if cb, _, err := hystrix.GetCircuit(name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(([]byte)("Can get Circuit Breaker\n"))
	} else {
		cb.ReportEvent([]string{"success"}, time.Now(), 0)

		w.WriteHeader(http.StatusOK)
		w.Write(([]byte)("Closed.\n"))
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte)("Missing or empty parameter\n"))
		return
	}

	if cb, _, err := hystrix.GetCircuit(name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(([]byte)("Can get Circuit Breaker\n"))
	} else {
		w.WriteHeader(http.StatusOK)
		allow := "Allow:No\n"
		if cb.AllowRequest() {
			allow = "Allow:Yes\n"
		}
		isOpen := "Opened:No\n"
		if cb.IsOpen() {
			isOpen = "Opened:Yes\n"
		}

		w.Write(([]byte)(allow + isOpen))
	}
}

func toggleOpenHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	tg := r.URL.Query().Get("value")
	if len(name) == 0 || (tg != "true" && tg != "false") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte)("Missing or empty parameter\n"))
		return
	}

	if cb, _, err := hystrix.GetCircuit(name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(([]byte)("Can get Circuit Breaker\n"))
	} else {
		cb.ToggleForceOpen(tg == "true")
		cb.AllowRequest()
		w.WriteHeader(http.StatusOK)
		w.Write(([]byte)("Done.\n"))

		log.Printf("%#v", *cb)

	}
}

func main() {

	// hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
	// 	Timeout:               100,
	// 	MaxConcurrentRequests: 100,
	// 	ErrorPercentThreshold: 50,
	// })

	// Launch streamHandler for Hystrix dashboard
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "8222"), hystrixStreamHandler)

	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/configure", configureHandler)
	http.HandleFunc("/toggleOpen", toggleOpenHandler)
	http.HandleFunc("/close", closeHandler)
	http.HandleFunc("/status", statusHandler)
	http.ListenAndServe(":8221", nil)

	output := make(chan bool, 1)
	<-output

}
