<template>
  <main class="max-w-2xl mx-auto py-8 md:py-16">
    <section 
      class="bg-white p-8 md:p-12 shadow-moxin-soft border border-moxin-border rounded-sm relative overflow-hidden transition-all duration-600 hover:shadow-md"
    >
      <div class="absolute top-0 right-8 w-px h-12 bg-moxin-shazhu opacity-20"></div>

      <div v-if="loading" class="animate-pulse space-y-4">
        <div class="h-6 bg-gray-100 w-3/4"></div>
        <div class="h-6 bg-gray-100 w-1/2"></div>
      </div>

      <div v-else class="space-y-12">
        <article class="relative">
          <span class="absolute -left-4 -top-2 text-4xl opacity-10 font-serif">“</span>
          <p class="text-2xl md:text-3xl leading-relaxed tracking-wide py-4">
            {{ dailySentence.content }}
          </p>
          <span class="absolute -right-2 -bottom-4 text-4xl opacity-10 font-serif">”</span>
        </article>

        <div class="flex justify-end items-center space-x-2 text-sm md:text-base opacity-70">
          <span class="w-8 h-px bg-moxin-ink opacity-30"></span>
          <span>{{ dailySentence.author }}</span>
          <span v-if="dailySentence.source" class="before:content-['《'] after:content-['》']">
            {{ dailySentence.source }}
          </span>
        </div>
      </div>

      <footer class="mt-16 flex items-center justify-between border-t border-moxin-border pt-6">
        <div class="flex space-x-6">
          <button 
            @click="handleToggleLike" 
            class="group flex items-center space-x-1 transition-colors"
            :class="isLiked ? 'text-moxin-shazhu' : 'hover:text-moxin-shazhu'"
          >
            <span class="text-xl">{{ isLiked ? '♥' : '♡' }}</span>
            <span class="text-xs tracking-tighter uppercase">{{ isLiked ? '已藏' : '收藏' }}</span>
          </button>
          
          <button class="flex items-center space-x-1 hover:text-moxin-ink opacity-60 hover:opacity-100 transition-all">
            <span class="text-lg">♞</span>
            <span class="text-xs tracking-tighter uppercase">生成海报</span>
          </button>
        </div>

        <div class="text-[10px] opacity-30 tracking-widest text-moxin-ink">
          {{ currentDate }}
        </div>
      </footer>
    </section>

    <p class="text-center mt-12 text-xs opacity-30 tracking-widest">
      素纸一张，写尽心城
    </p>
  </main>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { getDailySentence } from '../../utils/api/api'; // 引入封装好的逻辑
const loading = ref(true);
const isLiked = ref(false);
const dailySentence = ref({ content: '', author: '', source: '' });

// 计算当前日期（逻辑保持不变）
const currentDate = computed(() => {
  const date = new Date();
  return `${date.getFullYear()} / ${date.getMonth() + 1} / ${date.getDate()}`;
});

// 核心逻辑：现在只需调用封装好的 API
const loadData = async () => {
  loading.value = true;
  dailySentence.value = await getDailySentence();
  loading.value = false;
};

const handleToggleLike = () => {
  isLiked.value = !isLiked.value;
};

onMounted(() => {
  loadData();
});
</script>