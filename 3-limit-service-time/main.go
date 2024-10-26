//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import "time"

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds

}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User, pid int) bool {
	println("starting the process:", pid, " for user:", u.ID, " with time used:", u.TimeUsed)
	if !u.IsPremium && u.TimeUsed >= 10 {
		return false
	}
	var ticker *time.Ticker = time.NewTicker(10 * time.Second)
	if !u.IsPremium {
		tickerDuration := time.Duration(0)
		if 10 >= u.TimeUsed {
			tickerDuration = time.Duration(10 - u.TimeUsed)
		}
		ticker = time.NewTicker(tickerDuration * time.Second)
	}
	startedAt := time.Now()
	defer func() {
		timeDuration := time.Since(startedAt).Seconds()
		u.TimeUsed += int64(timeDuration)
		println("ending the process:", pid, " for user:", u.ID, " after running for duration: ", u.TimeUsed)
	}()

	processChan := make(chan bool)
	go func() {
		process()
		processChan <- true
	}()

	if !u.IsPremium {

		select {
		case <-ticker.C:
			// println("stopping the process:", pid, " due to time exceed for user:", u.ID, " timeUsed:", u.TimeUsed)
			return true
		case <-processChan:
			return true
		}
	}

	<-processChan
	return true

	// return true
}

func main() {
	RunMockServer()
}
