package goose

import (
	resty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type HtmlRequester interface {
	fetchHTML(string) (string, error)
}

// Crawler can fetch the target HTML page
type htmlrequester struct {
	config Configuration
}

// NewCrawler returns a crawler object initialised with the URL and the [optional] raw HTML body
func NewHtmlRequester(config Configuration) HtmlRequester {
	return htmlrequester{
		config: config,
	}
}

func (hr htmlrequester) fetchHTML(url string) (string, error) {
	client := resty.New()
	client.SetTimeout(hr.config.timeout)
	resp, err := client.R().
		SetHeader("Content-Type", "text/html").
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36").
		SetHeader("Accept", "text/javascript, application/javascript, application/ecmascript, application/x-ecmascript, */*; q=0.01").
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetHeader("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8").
		SetHeader("Cookie", "_gid=GA1.2.1990237071.1698737112; _gat_gtag_UA_46560124_1=1; _ga_VS4T955HM4=GS1.1.1698737111.1.1.1698738565.0.0.0; _ga=GA1.1.2108271516.1698737112").
		SetHeader("Referer", "https://doaj.org/search/articles?ref=homepage-box&source=%7B%22query%22%3A%7B%22query_string%22%3A%7B%22query%22%3A%22climate%20change%22%2C%22default_operator%22%3A%22AND%22%7D%7D%2C%22size%22%3A%22200%22%2C%22track_total_hits%22%3Atrue%7D").
		SetHeader("Sec-Ch-Ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"").
		SetHeader("Sec-Ch-Ua-Mobile", "?0").
		SetHeader("Sec-Ch-Ua-Platform", "\"macOS\"").
		SetHeader("Sec-Fetch-Dest", "empty").
		SetHeader("Sec-Fetch-Mode", "cors").
		SetHeader("Sec-Fetch-Site", "same-origin").
		SetHeader("X-Requested-With", "XMLHttpRequest").

		Get(url)

	if err != nil {
		return "", errors.Wrap(err, "could not perform request on "+url)
	}
	if resp.IsError() {
		return "", &badRequest{Message: "could not perform request with " + url + " status code " + string(resp.StatusCode())}
	}
	return resp.String(), nil
}

type badRequest struct {
	Message string `json:"message,omitempty"`
}

func (BadRequest *badRequest) Error() string {
	return "Required request fields are not filled"
}
