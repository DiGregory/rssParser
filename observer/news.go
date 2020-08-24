package observer

import (
	"context"
	"github.com/DiGregory/s7testTask/proto"
	"github.com/DiGregory/s7testTask/storage"
)


type NewsService struct {
	storage storage.NewsStorager
}

func NewNewsService(storage storage.NewsStorager) *NewsService {
	return &NewsService{storage: storage}
}

func (nss *NewsService) GetNews(ctx context.Context, request *proto.GetNewsRequest) (*proto.GetNewsResponse, error) {
	var limit, offset *int32
	if request.Limit != 0 {
		limit = &request.Limit
	}
	if request.Offset != 0 {
		offset = &request.Offset
	}
	news, err := nss.storage.GetNews(limit, offset)
	if err != nil {
		return nil, err
	}
	response := new(proto.GetNewsResponse)
	for _, n := range news {
		response.News = append(response.News, &proto.News{
			Id:                   n.ID,
			Title:                n.Title,
			Description:          n.Description,
			Link:                 n.Link,
		})
	}
	return response, nil
}
