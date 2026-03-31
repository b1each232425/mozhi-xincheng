import { createRouter, createWebHistory } from "vue-router";
// 导入布局
import Layout from "../layouts/Layout.vue";
// 导入页面
import HomeView from "../views/HomeView.vue";

const routes = [
  {
    path: "/",
    component: Layout, // 使用布局组件作为父级
    children: [
      {
        path: "", // 默认子路由，即访问 / 时显示
        name: "Home",
        component: HomeView,
      },
      // 后续可以在这里继续添加 ExploreView 等
    ],
  },
  {
    path: "/Chenxiang",
    component: Layout,
    // name: 'ChenXiang',
    children: [
      {
        path: "",
        component: () => import("../views/chenxiang/ChenXiangView.vue"),
      },
      {
        path: "List",
        name: "ChenXiangList",
        component: () => import("../views/chenxiang/ChenXiangList.vue"),
      },
      {
        path: "Poet/:id",
        name: "ChenXiangPoet",
        component: () => import("../views/chenxiang/ChenXiangPoet.vue"),
      },
    ],
  },
  {
    path: "/luobi",
    component: Layout,
    children: [
      {
        path: "",
        name: "luobilist",
        component: () => import("../views/luobi/LuoBiList.vue"),
      },
      {
        path: "article/:id",
        name: "luobiarticle",
        component: () => import("../views/luobi/LuoBiArc.vue"),
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
