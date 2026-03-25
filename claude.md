# moxiang-zhicheng 前端开发规范
## Vue3 开发规范
- 使用 `<script setup>` 语法糖。
- 组件通信: 优先使用 Props 和 Emits，全局状态使用 Pinia。
- UI 库: 使用 。
- 响应式数据: 优先使用 `ref`，仅在对象层级较深时使用 `reactive`。
- 函数使用正则式写法
- 命名约束: 文件夹名小写连字符，组件名大写驼峰。
- 尽量使用Tailwind CSS

## 前端目录规范
- src/views 存放页面代码
- src/components 存放公共组件
- src/utils 存放逻辑复用
- src/styles 存放全局样式
- src/layouts 页面布局

<!-- ## git提交规范
- 格式：<type>(<scope>): <subject> 例子：feat(home): 增加每日诗词卡片, fix(api): 修复收藏接口调用失败。 -->