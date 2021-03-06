package crawler_test

import (
	"context"
	"io"
	"monzoCrawler/domain/crawler/internal/mocks"
	"monzoCrawler/domain/model"
	"net/url"
	"strings"
	"testing"
	"time"

	"monzoCrawler/domain/crawler"
)

func TestCrawler_Crawl(t *testing.T) {
	t.Run("send a signal when the crawler is done", func(t *testing.T) {
		t.Parallel()

		fetcherExtractorMock := &mocks.FetcherExtractorMock{
			FetchFunc: func(ctx context.Context, urlMoqParam url.URL) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("")), nil
			},
			ExtractFunc: func(reader io.Reader) (model.CrawlResult, error) {
				return model.CrawlResult{}, nil
			},
		}

		queue := &mocks.QueueMock{
			PushFunc: func(val interface{}) error {
				return nil
			},
		}

		crawlJob := model.CrawlJob{SeedURL: &url.URL{
			Scheme: "http",
			Host:   "www.google.com",
		}}

		doneChan := make(chan struct{})

		crwler := crawler.New(fetcherExtractorMock, crawlJob, queue, doneChan)

		go crwler.Crawl(context.Background())

		select {
		case <-time.After(time.Second):
			t.Error("expected to receive a signal when the crawler is done")
		case <-doneChan:
			return
		}

	})
}
