package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type ContentFetcher interface {
	Fetch(ctx context.Context, url string) ([]byte, error)
}

type ContentProcessor interface {
	Process(ctx context.Context, content []byte) (ProcessedData, error)
}

type ProcessedData struct {
	Title       string
	Description string
	Keywords    []string
	Timestamp   time.Time
	Source      string
}

type ContentAggregator struct {
	fetcher     ContentFetcher
	processor   ContentProcessor
	workerCount int
	rateLimiter *rate.Limiter
	wg          sync.WaitGroup
	closed      chan struct{}
	once        sync.Once
}

func (ca *ContentAggregator) FanOut(ctx context.Context, urls []string) ([]ProcessedData, []error) {
	jobs := make(chan string, len(urls))           //pecah job sesuai banyak url
	results := make(chan ProcessedData, len(urls)) //collect all results
	errors := make(chan error, len(urls))          //collect the errors

	for _, url := range urls {
		jobs <- url //channeling per jobs
	}
	close(jobs)

	ca.wg.Add(1)
	go func() {
		defer ca.wg.Done()
		ca.WorkerPool(ctx, jobs, results, errors)
	}()

	var (
		collected []ProcessedData
		errList   []error
	)

	for range len(urls) {
		select {
		case <-ctx.Done():
			errList = append(errList, ctx.Err())
		case res := <-results:
			collected = append(collected, res)
		case err := <-errors:
			errList = append(errList, err)
		}
	}
	return collected, errList
}

func (ca *ContentAggregator) Shutdown() error {
	ca.once.Do(func() {
		close(ca.closed)
		ca.wg.Wait()
	})
	return nil
}

func (ca *ContentAggregator) WorkerPool(
	ctx context.Context,
	jobs <-chan string,
	results chan<- ProcessedData,
	errors chan<- error) {
	var wg sync.WaitGroup

	for i := 0; i < ca.workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				select {
				case <-ctx.Done():
					errors <- ctx.Err()
					return
				case <-ca.closed:
					return
				default:
					if err := ca.rateLimiter.Wait(ctx); err != nil {
						errors <- err
						continue
					}
					content, err := ca.fetcher.Fetch(ctx, url)
					if err != nil {
						errors <- fmt.Errorf("fetch error: %w", err)
						continue
					}
					data, err := ca.processor.Process(ctx, content)
					if err != nil {
						errors <- fmt.Errorf("process error: %w", err)
						continue
					}

					data.Source = url
					results <- data
				}
			}
		}()
	}
	wg.Wait()
}

func NewContentAggregator(fetcher ContentFetcher, processor ContentProcessor, workerCount int, requestPerSecond int) *ContentAggregator {
	return &ContentAggregator{
		fetcher:     fetcher,
		processor:   processor,
		workerCount: workerCount,
		rateLimiter: rate.NewLimiter(rate.Limit(requestPerSecond), workerCount),
		closed:      make(chan struct{}),
	}
}

func (ca *ContentAggregator) FetchAndProcess(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, error) {
	results, errs := ca.FanOut(ctx, urls)

	if len(errs) > 0 {
		return results, fmt.Errorf("encountered %d errors", len(errs))
	}
	return results, nil
}

type HTTPFetcher struct {
	Client *http.Client
}

func (f *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

type HTMLProcessor struct{}

func (p *HTMLProcessor) Process(ctx context.Context, content []byte) (ProcessedData, error) {
	html := string(content)
	data := ProcessedData{
		Timestamp: time.Now(),
	}

	titleMatch := regexp.MustCompile(`(?i)<title>(.*?)</title>`).FindStringSubmatch(html)
	if len(titleMatch) > 1 {
		data.Title = titleMatch[1]
	}

	descMatch := regexp.MustCompile(`(?i)<meta\s+name=["']description["']\s+content=["'](.*?)["']`).FindStringSubmatch(html)
	if len(descMatch) > 1 {
		data.Description = descMatch[1]
	}

	keywordsMatch := regexp.MustCompile(`(?i)<meta\s+name=["']keywords["']\s+content=["'](.*?)["']`).FindStringSubmatch(html)
	if len(keywordsMatch) > 1 {
		data.Keywords = strings.Split(keywordsMatch[1], ",")
	}

	return data, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fetcher := &HTTPFetcher{Client: &http.Client{Timeout: 5 * time.Second}}
	processor := &HTMLProcessor{}

	agg := NewContentAggregator(fetcher, processor, 5, 2)
	defer agg.Shutdown()

	urls := []string{
		"https://www.google.com",
		"https://example.com",
	}

	results, err := agg.FetchAndProcess(ctx, urls)
	if err != nil {
		fmt.Println("Errors occurred:", err)
	}

	for _, data := range results {
		fmt.Printf("URL: %s\nTitle: %s\nDesc: %s\nKeywords: %v\n\n", data.Source, data.Title, data.Description, data.Keywords)
	}
}
