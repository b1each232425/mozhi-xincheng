// src/utils/api.js


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

// 获取随机短文
export const getDailyHitokoto = async () => {
  try {
    // 这里以一言的“文学”和“诗词”分类为例
    const res = await fetch('https://v1.hitokoto.cn/?c=d&c=i');
    const data = await res.json();
    return data; 
  } catch (e) {
    return "此时无声，唯笔尖游走。"; // 降级兜底
  }
};