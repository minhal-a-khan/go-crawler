package urlFetcherExtractor

import (
	"context"
	"golang.org/x/net/html"
	"io"
	"monzoCrawler/dao"
	"net/http"
	"time"
)

type HTTPFetcherExtractor struct {
	client *http.Client
}

func NewHTTPFetcherExtractor(timeout time.Duration) HTTPFetcherExtractor {
	return HTTPFetcherExtractor{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (fe HTTPFetcherExtractor) Fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	resp, err := fe.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (fe HTTPFetcherExtractor) Extract(contents io.Reader) (dao.CrawlResult, error) {
	return fe.getLinks(contents), nil
}

//Collect all links from response body and return it as an array of strings
func (fe *HTTPFetcherExtractor) getLinks(body io.Reader) dao.CrawlResult {
	crawlResult := dao.CrawlResult{NewJobs: []dao.CrawlJob{}}

	z := html.NewTokenizer(body)
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			//todo: links list shoudn't contain duplicates
			return crawlResult
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						crawlResult.NewJobs = append(crawlResult.NewJobs, dao.CrawlJob{SeedURL: attr.Val})
					}
				}
			}

		}
	}
}