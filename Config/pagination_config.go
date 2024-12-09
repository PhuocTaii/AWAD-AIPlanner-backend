package config

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PagingConfig struct {
	Limit int
	Page  int
}

func NewPagingConfig(c *gin.Context, inputLimit int, inputPage int) *PagingConfig {
	page := inputPage
	limit := inputLimit

	return &PagingConfig{
		Limit: limit,
		Page:  page,
	}
}

// GetOffset tính toán số lượng bản ghi bỏ qua dựa trên trang và giới hạn
func (p *PagingConfig) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

func PaginatedFind[T any](ctx *gin.Context, collection *mongo.Collection, config *PagingConfig, filter bson.M, sort bson.M) ([]T, int, int, error) {
	// Tính toán offset và limit từ cấu hình phân trang
	offset := config.GetOffset()

	// Thiết lập tùy chọn truy vấn
	findOptions := options.Find()
	findOptions.SetSkip(int64(offset))        // Sử dụng skip để bỏ qua bản ghi
	findOptions.SetLimit(int64(config.Limit)) // Giới hạn số lượng kết quả trả về
	if sort != nil {
		findOptions.SetSort(sort) // Sắp xếp kết quả nếu cần
		//handle mongoDB sort

	}

	// fmt.Println("Sort: ", sort)

	// Lấy dữ liệu từ collection
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, 0, err
	}
	defer cursor.Close(ctx)

	totalItems, err := collection.CountDocuments(ctx, filter)

	//Get total pages:
	totalPages := (totalItems + int64(config.Limit) - 1) / int64(config.Limit)

	// Kết quả sẽ là một slice chứa các phần tử kiểu T
	var results []T
	for cursor.Next(ctx) {
		var item T
		if err := cursor.Decode(&item); err != nil {
			return nil, 0, 0, err
		}
		results = append(results, item)
	}

	// Kiểm tra lỗi khi lặp qua cursor
	if err := cursor.Err(); err != nil {
		return nil, 0, 0, err
	}

	return results, int(totalPages), int(totalItems), nil
}
