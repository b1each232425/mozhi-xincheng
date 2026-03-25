import { createRouter, createWebHistory } from 'vue-router';
// 导入布局
import Layout from '../layouts/Layout.vue';
// 导入页面
import HomeView from '../views/HomeView.vue';

const routes = [
  {
    path: '/',
    component: Layout, // 使用布局组件作为父级
    children: [
      {
        path: '', // 默认子路由，即访问 / 时显示
        name: 'Home',
        component: HomeView
      },
      // 后续可以在这里继续添加 ExploreView 等
    ]
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;