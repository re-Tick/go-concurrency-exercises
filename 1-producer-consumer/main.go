//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func producer(cancel context.CancelFunc, stream Stream, tweetChannel chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			cancel()
			return
		}
		tweetChannel <- tweet
		// tweets = append(tweets, tweet)
	}
}

func consumer(ctx context.Context, tweetChannel chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-tweetChannel:
			if t.IsTalkingAboutGo() {
				fmt.Println(t.Username, "\ttweets about golang")
			} else {
				fmt.Println(t.Username, "\tdoes not tweet about golang")
			}
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	tweetChannel := make(chan *Tweet)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	ctx, cancel := context.WithCancel(context.Background())
	// Producer
	go producer(cancel, stream, tweetChannel, wg)

	// Consumer
	go consumer(ctx, tweetChannel, wg)
	wg.Wait()
	fmt.Printf("Process took %s\n", time.Since(start))
}
