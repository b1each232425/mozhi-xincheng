/** @type {import('tailwindcss').Config} */
export default {
  // 按照你的规范，扫描 src 目录下所有 Vue 组件
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      // 1. 定义“墨心纸城”专属色板
      colors: {
        'moxin': {
          'paper': '#FBFBF2',   // 主背景：淡淡的宣纸黄，护眼且具质感
          'ink': '#2C2C2C',     // 主文字：深灰黑，模拟水墨在纸上的晕染感
          'shazhu': '#B22222',  // 点缀色：朱砂红，用于收藏、印章等关键交互
          'border': '#E0DDD5',  // 边框色：极浅的灰色，模拟纸张折痕或界格
        }
      },
      // 2. 字体配置：强制开启衬线体模式
      fontFamily: {
        'serif': [
          'Noto Serif SC', 
          'Source Han Serif SC', 
          'Source Serif Pro', 
          'STSong', 
          'serif'
        ],
      },
      // 3. 间距与阴影增强
      boxShadow: {
        'moxin-soft': '0 4px 20px -2px rgba(44, 44, 44, 0.05)', // 极其轻微的墨色阴影
      },
      // 4. 自定义交互动画
      transitionDuration: {
        '600': '600ms',
      }
    },
  },
  plugins: [],
}