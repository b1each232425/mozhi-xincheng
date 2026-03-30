package main

import (
	"context"
	"fmt"
	"log"
	"moxin-zhicheng/internal/config"
	"moxin-zhicheng/internal/database"
	"moxin-zhicheng/internal/logger"
	model "moxin-zhicheng/internal/models"
	"moxin-zhicheng/internal/redis" // 导入你的 redis 包
	"moxin-zhicheng/utils"
	"strconv" // 用于转换 ID

	"github.com/qdrant/go-client/qdrant"
)

func main() {
	logger.InitLogger("config_dev")
	config.InitConfig()
	database.InitDB()
	redis.InitRedis() // 初始化你自定义的 Redis 连接

	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6424,
	})
	if err != nil {
		log.Fatalf("连接 Qdrant 失败: %v", err)
	}
	defer client.Close()

	// 1. 从 redis 包中的 RedisClient 读取进度
	// 注意这里必须使用 redis.RedisClient
	lastIDStr, _ := redis.RedisClient.Get(context.Background(), "sync:last_poetry_id").Result()
	lastID, _ := strconv.ParseUint(lastIDStr, 10, 64)

	var poetries []model.Poetry
	// 2. 核心逻辑：只查 ID 大于上次进度且不是 AI 翻译的数据
	database.DB.Where("id > ? AND translation IS NOT NULL AND translation != '' AND translation NOT LIKE ?",
		lastID, "[AI翻译]%").
		Order("id asc").
		Find(&poetries)

	waitUpsert := true

	for _, p := range poetries {
		corpusText := fmt.Sprintf("原文:%s\n译文:%s", p.Paragraphs, p.Translation)

		// 获取向量
		vector, err := utils.GetLocalEmbedding(corpusText)
		if err != nil {
			log.Printf("向量获取失败 ID %d: %v", p.ID, err)
			continue
		}

		// 3. 严格按照 Qdrant Oneof 结构构造
		_, err = client.Upsert(context.Background(), &qdrant.UpsertPoints{
			CollectionName: "poetry_corpus",
			Wait:           &waitUpsert,
			Points: []*qdrant.PointStruct{
				{
					Id: &qdrant.PointId{
						PointIdOptions: &qdrant.PointId_Num{
							Num: uint64(p.ID),
						},
					},
					Vectors: &qdrant.Vectors{
						VectorsOptions: &qdrant.Vectors_Vector{
							Vector: &qdrant.Vector{
								Data: vector,
							},
						},
					},
					Payload: map[string]*qdrant.Value{
						"title":  {Kind: &qdrant.Value_StringValue{StringValue: p.Title}},
						"author": {Kind: &qdrant.Value_StringValue{StringValue: p.Author}},
						"trans":  {Kind: &qdrant.Value_StringValue{StringValue: p.Translation}},
					},
				},
			},
		})

		if err != nil {
			log.Printf("写入 Qdrant 失败 《%s》: %v", p.Title, err)
			continue
		}

		// 4. 同步成功后，在 Redis 中更新当前 ID 进度
		redis.RedisClient.Set(context.Background(), "sync:last_poetry_id", p.ID, 0)
		fmt.Printf("✅ 已同步 ID %d: 《%s》\n", p.ID, p.Title)
	}

	fmt.Println("[同步流程运行结束]")
}
