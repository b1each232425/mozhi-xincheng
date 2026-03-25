// src/utils/api.js

/**
 * 自动将诗词存入本地存储，用于未来数据迁移和离线显示
 */
const saveToLocal = (sentence) => {
  try {
    const history = JSON.parse(localStorage.getItem('moxin_history') || '[]');
    // 简单去重：如果内容不重复则存入
    if (!history.find(item => item.content === sentence.content)) {
      history.push({
        ...sentence,
        savedAt: new Date().toISOString()
      });
      localStorage.setItem('moxin_history', JSON.stringify(history));
    }
  } catch (e) {
    console.error("本地存储失败", e);
  }
};

/**
 * 获取每日诗词
 * 封装后，未来更换为 Go 后端只需修改此处的 fetch 地址
 */
export const getDailySentence = async () => {
  try {
    const response = await fetch('https://v1.jinrishici.com/all.json');
    if (!response.ok) throw new Error('网络请求失败');
    
    const data = await response.json();
    
    // 统一数据格式，方便未来数据库建模
    const sentence = {
      content: data.content,
      author: data.author,
      source: data.origin // 将外部字段统一为你的 source
    };

    // 执行你的“积累”计划：存入本地
    saveToLocal(sentence);
    
    return sentence;
  } catch (error) {
    console.warn("使用备用诗词:", error);
    // 降级处理：返回兜底数据
    return {
      content: "人生如逆旅，我亦是行人。",
      author: "苏轼",
      source: "临江仙·送钱穆父"
    };
  }
};