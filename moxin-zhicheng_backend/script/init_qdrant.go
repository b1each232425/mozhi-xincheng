package main

import (
	"context"
	"fmt"
	"log"

	"github.com/qdrant/go-client/qdrant"
)

func main() {
	// 连接本地 Docker 运行的 Qdrant
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6424,
	})
	if err != nil {
		log.Fatal("无法连接 Qdrant:", err)
	}

	// 创建名为 poetry_corpus 的集合
	err = client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "poetry_corpus",
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     1024, // bge-m3 模型的维度
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		fmt.Printf("集合可能已存在或创建失败: %v\n", err)
	} else {
		fmt.Println("Qdrant 集合 poetry_corpus 创建成功！")
	}
}
