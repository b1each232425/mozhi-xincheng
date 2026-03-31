// src/utils/api.js
const BASE_URL = 'http://localhost:8080'; // 后端地址，未来更换为 Go 后端只需修改此处

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

export const getPoetryList = async (keyword, page = 1, pageSize = 10) => {
  try {
    // 必须手动拼接查询字符串
    const url = `${BASE_URL}/chenxiang/list?keyword=${encodeURIComponent(keyword)}&page=${page}&page_size=${pageSize}`;
    
    const res = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    });

    if (!res.ok) throw new Error('网络响应错误');
    
    const data = await res.json(); // 解析 JSON
    return data; // 返回给组件处理
  } catch (err) {
    console.error("检索失败:", err);
    return { code: 500, data: [] };
  }
};

export const getStars = async () => {
  try {
    const res = await fetch(`${BASE_URL}/chenxiang/starTags`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    });
    if (!res.ok) throw new Error('网络响应错误');
    const data = await res.json();
    return data;
  } catch (err) {
    console.error("获取功能星失败:", err);
    return { code: 500, data: [] };
  }
};

export const getPoemDetail = async (id) => {
  try {
    const res = await fetch(`${BASE_URL}/chenxiang/singlePoem/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    });
    if (!res.ok) throw new Error('网络响应错误');
    const data = await res.json();
    return data;
  } catch (err) {
    console.error("获取诗词详情失败:", err);
    return { code: 500, data: {} };
  }
};

export const getLuoBiList = async (page = 1, pageSize = 10) => {
  try{
    const res = await fetch(`${BASE_URL}/luobi/articles?page=${page}&page_size=${pageSize}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'}});
    if (!res.ok) throw new Error('网络响应错误');
    const data = await res.json();
    return data;
  } catch (err) {
    console.error("获取落笔列表失败:", err);
    return { code: 500, data: [] };
}}

export const createLuoBiArticle = async (articleData) => {
  try {
    const res = await fetch(`${BASE_URL}/luobi/article`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(articleData)
    });
    if (!res.ok) throw new Error('网络响应错误');
    const data = await res.json();
    return data;
  } catch (err) {
    console.error("创建落笔文章失败:", err);
    return { code: 500, data: {} };
  }
};

export const getLuoBiArticle = async (id) => {
  try {
    const res = await fetch(`${BASE_URL}/luobi/article/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    });
    if (!res.ok) throw new Error('网络响应错误');
    const data = await res.json();
    return data;
  } catch (err) {
    console.error("获取落笔文章失败:", err);
    return { code: 500, data: {} };
  }
};
