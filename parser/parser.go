package parser

import (
	"strings"
	"time"

	"github.com/DiGregory/rssParser/storage"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

type ParserPooler interface {
	Start(links []string)
	newWorker(id int, work chan string)
}

type News struct {
	Link      string
	UpdatedAt time.Time
}

type ParserPool struct {
	workers  int
	timeout  time.Duration
	keyWords []string
	storage  storage.NewsStorager
}

func NewPool(threads, timeout int, keyWords []string, newsStorage storage.NewsStorager) ParserPooler {
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
	for range time.NewTicker(time.Second).C {
		for _, n := range news {
			if time.Since(n.UpdatedAt) >= p.timeout {
				n.UpdatedAt = time.Now()
				jobs <- n.Link
			}
		}
	}
}

func (p *ParserPool) newWorker(id int, work chan string) {
	fp := gofeed.NewParser()
	go func() {
		for {
			link := <-work
			rss, err := fp.ParseURL(link)
			if err != nil {
				logrus.WithError(err).WithField("worker", id).Error("Parse url error")
				continue
			}
			news := make([]*storage.News, 0)
			for _, i := range rss.Items {
				for _, kw := range p.keyWords {
					if strings.Contains(i.Title, kw) || strings.Contains(i.Description, kw) {
						news = append(news, &storage.News{
							Title:       i.Title,
							Description: i.Description,
							Link:        i.Link,
						})
						break
					}
				}
			}
			err = p.storage.CreateNews(news)
			if err != nil {
				logrus.WithError(err).WithField("worker", id).Error("Create news error")
				continue
			}

		}
	}()

}
