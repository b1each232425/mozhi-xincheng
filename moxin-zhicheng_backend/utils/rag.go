package utils

import (
	"context" // 解决 Background 报错
	"fmt"     // 解决 Sprintf 报错
	"moxin-zhicheng/internal/logger"
	model "moxin-zhicheng/internal/models"
	"strings" // 解决 HasPrefix, Builder, WriteString 报错

	"github.com/qdrant/go-client/qdrant" // 解决 Search 等 Qdrant 相关报错
)

// QdrantClient 必须由外部初始化并赋值
var QdrantClient *qdrant.Client

func InitQdrant() {
	var err error
	// (gRPC)
	QdrantClient, err = qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6424,
	})
	if err != nil {
		panic("无法连接到 Qdrant 服务: " + err.Error())
	}
}

// SearchSimilarCorpus 修正了所有未解析引用
func SearchSimilarCorpus(paragraphs string) []model.Poetry {
	queryVector, err := GetLocalEmbedding(paragraphs)
	if err != nil || queryVector == nil {
		logger.Error("RAG检索失败", err)
		return nil
	}

	searchResult, err := QdrantClient.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: "poetry_corpus",
		Query:          qdrant.NewQuery(queryVector...), // 传入你的 []float32 向量
		Limit:          qdrant.PtrOf(uint64(20)),
		WithPayload:    qdrant.NewWithPayload(true),
	})

	if err != nil {
		logger.Error("Qdrant查询失败", err)
		return nil
	}

	var references []model.Poetry
	for _, hit := range searchResult {
		// 这里的 Payload 是 map[string]*qdrant.Value，需要类型断言获取字符串值
		paragraphsVal, ok := hit.Payload["paragraphs"]
		translationVal, ok2 := hit.Payload["translation"]

		if !ok || !ok2 {
			continue
		}

		p := model.Poetry{
			Paragraphs:  paragraphsVal.GetStringValue(),
			Translation: translationVal.GetStringValue(),
		}
		// 必须用 strings.HasPrefix
		if !strings.HasPrefix(p.Translation, "[AI补全]") {
			references = append(references, p)
		}
	}
	return references
}

// GenerateWithRAG 修正了 strings.Builder 相关报错
func GenerateWithRAG(target model.Poetry, refs []model.Poetry) model.Poetry {
	// 必须用 strings.Builder
	var contextStr strings.Builder
	if len(refs) > 0 {
		contextStr.WriteString("### 参考示例：\n") // 必须用 contextStr.WriteString
		for i, r := range refs {
			// 必须用 fmt.Sprintf
			contextStr.WriteString(fmt.Sprintf("%d. 原文：%s\n   译文：%s\n", i+1, r.Paragraphs, r.Translation))
		}
	}

	prompt := fmt.Sprintf(`你是一位精通中国古典文学与现代汉语修辞的翻译专家。
当前任务是参考已有的高质量译文，将给定的古诗词（主要是《楚辞》或同类文体）翻译为优雅、准确的现代汉语。

%s

### 翻译规范：
1. **信达雅**：译文需忠实原文，语义准确，且文字组织要优美，保留古诗词的意境。
2. **拒绝冗余**：**禁止**输出“这段话的意思是”、“翻译如下”等废话。直接输出译文内容。
3. **术语统一**：参考示例中的人名、地名、祭祀用语等，如在目标原文中出现，请保持翻译风格一致。
4. **处理生僻字**：遇到生僻字或特定典故，请根据上下文语境转化为现代读者易懂但又不失庄重的表达。

### 待处理目标：
原文：%s

请直接提供对应的现代汉语译文：`, contextStr.String(), target.Paragraphs)

	// 调用你实现的 Ollama 接口
	aiTranslation, err := CallOllama("qwen2.5:3b", prompt)
	fmt.Printf("ID: %d, AI返回内容: [%s]\n", target.ID, aiTranslation) // 加这行
	if err != nil {
		return model.Poetry{Translation: "AI 生成失败"}
	}

	return model.Poetry{
		Translation: aiTranslation,
	}
}
