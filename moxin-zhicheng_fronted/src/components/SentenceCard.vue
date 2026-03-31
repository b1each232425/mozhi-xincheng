<template>
  <div class="relative w-full max-w-5xl mx-auto pt-32 pb-24 px-6 md:px-12 rounded-xl shadow-2xl border border-white/10 font-['KaiTi','楷体','serif']"
  :class="[currentTheme.bg, currentTheme.border]">
    
    <div class="absolute top-8 left-8 md:left-12 z-30 flex space-x-8 items-center overflow-x-auto no-scrollbar max-w-[85%] pb-4 mask-fade-edges">
      <div 
        v-for="dateItem in dateRange" 
        :key="dateItem.full"
        @click="$emit('update:selectedDate', dateItem.full)" 
        class="flex-shrink-0 cursor-pointer transition-all duration-500 flex flex-col items-center group"
        :class="modelValue === dateItem.full ? 'scale-110' : 'opacity-40 hover:opacity-80'"
      >
        <span class="text-[10px] font-mono text-white/60 tracking-widest mb-1 transition-colors">
          {{ dateItem.year }}
        </span>
        <span 
          class="text-xl font-bold transition-all duration-300" 
          :class="modelValue === dateItem.full ? 'text-[#f6f4ec] drop-shadow-md' : 'text-white/70'"
        >
          {{ dateItem.display }}
        </span>
        <div 
          v-if="modelValue === dateItem.full" 
          class="w-1.5 h-1.5 bg-[#d4af37] rounded-full mt-2 shadow-[0_0_8px_#d4af37]"
        ></div>
      </div>
    </div>

    <section class="bg-[#f8f6f0] p-10 md:p-20 shadow-[0_20px_50px_rgba(0,0,0,0.2)] rounded-sm relative overflow-hidden transition-all duration-500">
      
      <div class="absolute top-0 right-12 w-[1px] h-24 bg-[#a63838] opacity-30"></div>

      <div v-if="loading" class="animate-pulse space-y-8">
        <slot name="loading">
          <div class="h-6 bg-[#e0dcd0] w-3/4 rounded-sm opacity-50"></div>
          <div class="h-6 bg-[#e0dcd0] w-1/2 rounded-sm opacity-50"></div>
        </slot>
      </div>

      <div v-else class="space-y-16 mt-4">
        <article class="relative px-8 md:px-12">
          <span class="absolute left-0 -top-8 text-7xl text-[#a63838] opacity-15 font-serif leading-none">“</span>
          <p class="text-2xl md:text-[28px] leading-[2.2] tracking-[4px] text-[#2b2b2b] text-justify relative z-10">
            {{ sentence.content }}
          </p>
          <span class="absolute right-0 -bottom-12 text-7xl text-[#a63838] opacity-15 font-serif leading-none">”</span>
        </article>

        <div class="flex justify-end items-center space-x-4 text-lg md:text-xl text-[#5c554b]">
          <span class="w-12 h-[1px] bg-[#5c554b] opacity-40"></span>
          <span class="tracking-[2px]">{{ sentence.author }}</span>
          <span v-if="sentence.source" class="tracking-[2px] before:content-['《'] after:content-['》'] text-[#8c222c]">
            {{ sentence.source }}
          </span>
        </div>
      </div>

      <footer class="mt-24 flex items-center justify-between border-t border-[#d3cbba] pt-8">
        <div class="flex space-x-6">
          <slot name="actions"></slot>
        </div>
        <div class="text-[10px] opacity-40 tracking-[0.3em] text-[#2b2b2b] font-mono uppercase">
          Daily Inspiration
        </div>
      </footer>
    </section>

    <div class="absolute bottom-6 left-12 opacity-20 text-[10px] text-white/50 tracking-[0.4em] uppercase font-mono">
      Moxin Project · Est. 2026
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { getCurrentSeason } from '../../utils/season';
defineProps({
  sentence: {
    type: Object,
    required: true,
    default: () => ({ content: '', author: '', source: '' })
  },
  loading: {
    type: Boolean,
    default: false
  },
  modelValue: String, 
  dateRange: Array
});
defineEmits(['update:selectedDate']);

const seasonTheme = {
  spring: { bg: 'bg-[#4a6d55]', border: 'border-white/20', accent: '#8da48f' }, // 春：森绿
  summer: { bg: 'bg-[#1a4a6e]', border: 'border-white/10', accent: '#5da9ad' }, // 夏：深海蓝
  autumn: { bg: 'bg-[#8c222c]', border: 'border-white/10', accent: '#d4af37' }, // 秋：原本的朱红
  winter: { bg: 'bg-[#2c2e3a]', border: 'border-white/10', accent: '#a1a3a6' }  // 冬：玄青/冷灰
};

// 2. 获取当前季节主题
const currentTheme = computed(() => {
  const season = getCurrentSeason();
  return seasonTheme[season] || seasonTheme.autumn; // 默认回退到红色
});
</script>

<style scoped>
@import "tailwindcss";

/* 隐藏滚动条但保留滚动功能 */
.no-scrollbar::-webkit-scrollbar { display: none; }
.no-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }

/* 为日期横向滚动条添加两端渐变遮罩，使其融入背景 */
.mask-fade-edges {
  -webkit-mask-image: linear-gradient(to right, transparent, black 5%, black 95%, transparent);
  mask-image: linear-gradient(to right, transparent, black 5%, black 95%, transparent);
}
</style>