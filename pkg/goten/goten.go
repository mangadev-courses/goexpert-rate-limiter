package goten

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Goten struct{}

func New() *Goten {
	return &Goten{}
}

func (v *Goten) LoadTest(url, apiKeyHeader string, requests int, concurrency int) error {
	fmt.Println("Starting the test")

	startTime := time.Now()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	type result struct {
		statusCode int
		latency    time.Duration
		err        error
	}

	jobs := make(chan int, requests)
	resultsCh := make(chan result, requests)

	worker := func() {
		for range jobs {
			start := time.Now()
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				resultsCh <- result{err: err, latency: time.Since(start)}
				continue
			}

			if apiKeyHeader != "" {
				req.Header.Add("API_KEY", apiKeyHeader)
			}

			resp, err := client.Do(req)
			latency := time.Since(start)
			if err != nil {
				resultsCh <- result{err: err, latency: latency}
				continue
			}

			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			resultsCh <- result{statusCode: resp.StatusCode, latency: latency}
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker()
		}()
	}

	for i := 0; i < requests; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	close(resultsCh)

	var totalLatency time.Duration
	var successCount, errorCount int
	statusCounts := make(map[int]int)

	for res := range resultsCh {
		if res.err != nil {
			errorCount++
		} else {
			successCount++
			totalLatency += res.latency
			statusCounts[res.statusCode]++
		}
	}

	elapsedTime := time.Since(startTime)

	fmt.Printf("\n--- Load Test Report ---\n")
	fmt.Printf("Total time execution: %v\n", elapsedTime)
	fmt.Printf("Total requests: %d\n", requests)
	fmt.Printf("Successful requests: %d\n", successCount)
	fmt.Printf("Errored requests: %d\n", errorCount)
	if successCount > 0 {
		avgLatency := totalLatency / time.Duration(successCount)
		fmt.Printf("Average latency: %v\n", avgLatency)
	}
	fmt.Println("Status Code Distribution:")
	for code, count := range statusCounts {
		fmt.Printf("  %d: %d\n", code, count)
	}

	return nil
}
