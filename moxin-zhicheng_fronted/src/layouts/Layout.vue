<template>
  <div class="min-h-screen bg-moxin-paper font-serif text-moxin-ink selection:bg-moxin-shazhu selection:text-white">
    
    <nav class="sticky top-0 z-50 bg-moxin-paper/80 backdrop-blur-md border-b border-moxin-border/50">
      <div class="max-w-4xl mx-auto px-6 h-16 flex items-center justify-between">
        <div 
          @click="handleNavigate('/')" 
          class="text-xl font-bold tracking-[0.3em] cursor-pointer hover:opacity-70 transition-opacity"
        >
          墨心
        </div>

        <div class="flex items-center space-x-8 text-sm tracking-widest">
          <button @click="handleNavigate('/explore')" class="hover:text-moxin-shazhu transition-colors">发现</button>
          <button @click="handleNavigate('/editor')" class="hover:text-moxin-shazhu transition-colors">落笔</button>
          <button @click="handleNavigate('/user')" class="group relative py-1">
            <span class="hover:text-moxin-shazhu transition-colors">纸城</span>
            <span class="absolute -top-1 -right-2 w-1 h-1 bg-moxin-shazhu rounded-full opacity-0 group-hover:opacity-100 transition-opacity"></span>
          </button>
        </div>
      </div>
    </nav>

    <main class="relative z-10">
      <router-view v-slot="{ Component }">
        <transition 
          name="fade" 
          mode="out-in"
        >
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <div class="fixed inset-0 pointer-events-none opacity-[0.03] mix-blend-multiply bg-[url('https://www.transparenttextures.com/patterns/p6-dark.png')]"></div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';

const router = useRouter();

// 函数使用正则式（箭头函数）写法
const handleNavigate = (path) => {
  router.push(path);
};
</script>

<style scoped>
/* 页面切换的淡入淡出效果，符合“小而美”的柔和感 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.4s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>