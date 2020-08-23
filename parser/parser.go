package parser

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
	"strings"
	"github.com/DiGregory/s7testTask/storage"
)

type ParserPooler interface {
	Start(links []string)
}

type News struct {
	Link      string
	UpdatedAt time.Time
}

type ParserPool struct {
	workers  int
	timeout  time.Duration
	keyWords []string
	storage  *storage.NewsStorage
}

func NewPool(threads, timeout int, keyWords []string, newsStorage *storage.NewsStorage) ParserPooler {
	return &ParserPool{
		workers:  threads,
		timeout:  time.Minute * time.Duration(timeout),
		keyWords: keyWords,
		storage:  newsStorage,
	}
}
func (p *ParserPool) Start(links []string) {
	jobs := make(chan string, 1)

	for i := 0; i < p.workers; i++ {
		p.newWorker(i, jobs)
	}

	news := make([]*News, 0, len(links))
	for _, link := range links {
		news = append(news, &News{
			Link: link,
		})
	}

	for {
		select {
		case <-time.NewTicker(time.Second).C:
			for _, n := range news {
				if time.Since(n.UpdatedAt).Minutes() > p.timeout.Minutes() {
					n.UpdatedAt = time.Now()
					jobs <- n.Link
				}
			}

		}

	}
}

func (p ParserPool) newWorker(id int, work chan string) {
	go func() {
		for {
			select {
			case link := <-work:
				fp := gofeed.NewParser()
				rss, err := fp.ParseURL(link)
				if err != nil {
					fmt.Println(err)
				}
				news := make([]*storage.News, 0)
				for _, i := range rss.Items {
					ok := false
					for _, kw := range p.keyWords {
						if strings.Contains(i.Title, kw) || strings.Contains(i.Description, kw) {
							ok = true
							news = append(news, &storage.News{
								Title:       i.Title,
								Description: i.Description,
								Link:        i.Link,
							})
							break
						}
					}
					if !ok {
						continue
					}
				}
				err = p.storage.CreateNews(news)
			}
		}
	}()

}
