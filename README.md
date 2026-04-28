# 🎋 墨心纸城 (Moxin-Zhicheng)

## ✨ 项目特色

### 🎨 **文化设计理念**
- 宣纸、水墨风格的视觉设计，融合中国传统美学
- 四季主题自适应系统，随着季节变化动态切换配色
- 精心定制的"墨心色板"，呈现专业的视觉层级

### 📚 **功能特性**
- ⭐ **每日诗词推荐**：精选每日佳句，支持按日期浏览历史诗词
- 🌟 **星空搜索界面**：交互式Canvas星空，点击标签进行诗词搜索
- 🏷️ **智能标签系统**：基于TextRank算法提取关键词，精准匹配诗词
- 🤖 **AI翻译补全**：集成Qwen2.5模型，基于RAG技术完善诗词译文
- 📱 **响应式设计**：完美适配桌面、平板、移动设备


### 项目截图 

<img width="2369" height="1243" alt="image" src="https://github.com/user-attachments/assets/f3c50a97-000c-43f1-ac73-3c4352d52284" />
<img width="2369" height="1243" alt="image" src="https://github.com/user-attachments/assets/357b1057-e268-4e4a-b92c-e11046f05035" />


## 📁 项目结构

```
mozhi-xincheng/
├── moxin-zhicheng_fronted/        # 前端项目（Vue3 + Vite）
│   ├── src/
│   │   ├── views/                 # 页面组件
│   │   │   ├── HomeView.vue       # 首页：每日诗词展示
│   │   │   └── chenxiang/         # 沉香模块
│   │   │       ├── ChenXiangView.vue       # 星空搜索界面
│   │   │       ├── ChenXiangList.vue       # 诗词列表页
│   │   │       └── ChenXiangPoet.vue       # 诗词详情页
│   │   ├── components/            # 公共组件
│   │   │   ├── SentenceCard.vue   # 诗词卡片组件
│   │   │   ├── HitokotoCard.vue   # 一言卡片组件
│   │   │   ├── Button.vue         # 按钮组件
│   │   │   └── MoxinVFX.vue       # 粒子特效组件
│   │   ├── layouts/               # 布局组件
│   │   │   └── Layout.vue         # 主布局（导航 + 路由）
│   │   ├── router/                # 路由配置
│   │   │   └── index.js           # 路由定义
│   │   ├── utils/                 # 工具函数
│   │   │   └── api/api.js         # API接口调用
│   │   ├── styles/                # 全局样式
│   │   │   └── index.css          # Tailwind配置入口
│   │   ├── App.vue                # 根组件
│   │   └── main.js                # 应用入口
│   ├── index.html                 # HTML模板
│   ├── vite.config.js             # Vite配置
│   ├── tailwind.config.js         # Tailwind配置（墨心色板）
│   ├── package.json               # 前端依赖
│   └── README.md                  # 前端文档
│
├── moxin-zhicheng_backend/        # 后端项目（Go + Gin）
│   ├── main.go                    # 应用入口
│   ├── internal/
│   │   ├── config/                # 配置管理（Viper）
│   │   ├── database/              # 数据库初始化（GORM）
│   │   ├── redis/                 # Redis客户端
│   │   ├── logger/                # 日志模块
│   │   ├── models/                # 数据模型
│   │   │   ├── poetry.go          # 诗词模型
│   │   │   ├── poetry_tag.go      # 标签模型
│   │   │   └── poetry_tag_relation.go  # 诗词-标签关系
│   │   └── handler/               # HTTP处理器
│   ├── routes/                    # 路由定义
│   ├── script/                    # 数据处理脚本
│   │   ├── chinese_poetry.go      # 🌟 诗词数据导入脚本
│   │   └── textrank.go            # 🌟 文本关键词提取脚本
│   ├── go.mod                     # Go模块定义
│   ├── go.sum                     # Go依赖锁定
│   ├── config.yaml                # 配置文件
│   └── .env.example               # 环境变量示例
│
├── claude.md                       # 前端开发规范
└── README.md                       # 本文件
```

---

## 🚀 快速开始

### 前置要求
- **Node.js** >= 16
- **Go** >= 1.21
- **PostgreSQL** 或 **MySQL**（推荐 PostgreSQL）
- **Redis** >= 6.0
- **Docker**（可选，用于Qdrant向量数据库）

### 1️⃣ 后端启动

#### 步骤 1：配置环境变量
```bash
cd moxin-zhicheng_backend
cp .env.example .env
# 编辑 .env 配置数据库、Redis等参数
```

#### 步骤 2：初始化数据库
```bash
# 创建表结构
go run main.go
```

#### 步骤 3：导入诗词数据
```bash
# 将 Chinese Poetry 数据集放在 script/chinese_poetry_Data/ 目录
# 运行导入脚本
cd script
go run chinese_poetry.go

# 输出示例：
# ✅ 导入成功: tang.json (1000条记录)
# ✅ 导入成功: song.json (1500条记录)
# 🎉 数据库灌顶完成！
```

#### 步骤 4：提取文本关键词
```bash
# 基于TextRank算法提取高质量标签
cd script
go run textrank.go

# 输出示例：
# 🔄 任务启动：从 ID > 0 开始处理
# 🚀 已投递至 ID: 150
# ✅ 任务全部完成！
```

#### 步骤 5：初始化Qdrant，爬取古诗文网译文和注解数据
cd script
go run init_qdrant.go(需要docker运行)
go run crawler.go

#### 步骤 6：构建知识库，AI补全其余诗歌的译文
cd script
go run sync_vector.go
go run ai_annotation.go


#### 启动后端服务
```bash
# 回到项目根目录
cd ..
go run main.go

# 服务将在 http://localhost:8080 启动
# 可访问 http://localhost:8080/debug/pprof/ 进行性能分析
```

### 2️⃣ 前端启动

#### 步骤 1：安装依赖
```bash
cd moxin-zhicheng_fronted
npm install
# 或使用 pnpm
pnpm install
```

#### 步骤 2：开发模式运行
```bash
npm run dev
# 前端将在 http://localhost:5220 启动
```


## 🛠️ 脚本说明

### `script/chinese_poetry.go` - 诗词数据导入
**用途**：将 Chinese Poetry 数据集（JSON格式）导入到数据库

**功能**：
- 🔄 自动识别朝代类型（唐、宋、楚辞、曹操等）
- 📊 统计导入成功/失败数量
- ✅ 提供详细的导入日志
- ⚠️ 标记未被识别的文件

**运行方式**：
```bash
# 确保数据集在以下路径
# script/chinese_poetry_Data/poet.tang.*.json
# script/chinese_poetry_Data/poet.song.*.json
# 等...

go run chinese_poetry.go
```

**输出示例**：
```
正在导入: poet.tang.0.json (type: tang)
✅ 导入成功: 1000条记录
正在导入: poet.song.0.json (type: song)
✅ 导入成功: 1500条记录
🎉 数据库灌顶完成！
```

---

### `script/textrank.go` - 文本关键词提取
**用途**：使用 TextRank 算法对诗词进行关键词提取和标签生成

**功能**：
- 🏷️ 基于 TextRank 算法自动提取关键词
- 🔄 支持断点续传（通过 `poetry_checkpoint.txt`）
- ⚡ 8个协程并行处理，高效提取
- 📈 自动更新标签权重
- 🔗 建立诗词-标签关系表

**运行方式**：
```bash
cd script
go run textrank.go

# 如果需要从特定ID继续（如ID=5000）
# 编辑 poetry_checkpoint.txt，改为 5000，再次运行
```

**断点续传机制**：
- 每处理150条诗词，保存检查点到 `poetry_checkpoint.txt`
- 如果任务中断，再次运行时会从上次断点继续
- 这样可以避免重复处理和网络中断导致的数据丢失

**输出示例**：
```
🔄 任务启动：从 ID > 0 开始处理
🚀 已投递至 ID: 150
🚀 已投递至 ID: 300
...
✅ 任务全部完成！
```

**提取结果**：
- `poetry_tag` 表：存储所有唯一标签及其权重
- `poetry_tag_relation` 表：诗词与标签的多对多关系

---


## 📝 开发规范

### Git 提交规范
```bash
# 格式: <type>(<scope>): <subject>
# 例子:
git commit -m "feat(home): 增加每日诗词卡片"
git commit -m "fix(api): 修复收藏接口调用失败"
git commit -m "refactor(search): 优化搜索算法性能"
```

### 前端开发规范
详见 [claude.md](./claude.md)

```
- src/views       → 页面代码
- src/components  → 公共组件
- utils       → 逻辑复用，api接口定义
- src/styles      → 全局样式
- src/layouts     → 页面布局
```

---

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request





Made with ❤️ by [b1each232425](https://github.com/b1each232425)

</div>
