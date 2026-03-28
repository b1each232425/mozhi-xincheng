<template>
 <div class="relative w-full max-w-6xl mx-auto pt-32 pb-24 px-8 md:px-16 bg-[#C21F31] rounded-2xl shadow-[10px_10px_40px_rgba(0,0,0,0.05)] border border-moxin-border/20">
    
    <div class="absolute top-10 left-10 md:left-16 z-30 flex space-x-6 items-center overflow-x-auto no-scrollbar max-w-[80%] pb-2">
      <div 
        v-for="dateItem in dateRange" 
        :key="dateItem.full"
        @click="$emit('update:selectedDate', dateItem.full)" 
        class="flex-shrink-0 cursor-pointer transition-all duration-300 flex flex-col items-center"
        :class="modelValue === dateItem.full ? 'scale-110' : 'opacity-30 hover:opacity-60'"
      >
        <span class="text-[10px] font-mono text-moxin-ink">{{ dateItem.year }}</span>
        <span class="text-lg font-serif text-moxin-ink font-bold" :class="modelValue === dateItem.full ? 'text-moxin-shazhu' : ''">
          {{ dateItem.display }}
        </span>
      </div>
    </div>

    <div class="absolute top-0 left-0 w-full flex justify-around px-20 -translate-y-1/2 z-40">
       </div>

    <section 
      class="bg-white p-10 md:p-16 shadow-moxin-soft border border-moxin-border/50 rounded-sm relative overflow-hidden transition-all duration-600 hover:shadow-md"
    >
      <div class="absolute top-0 right-10 w-px h-16 bg-moxin-shazhu opacity-20"></div>

      <div v-if="loading" class="animate-pulse space-y-6">
        <slot name="loading">
          <div class="h-8 bg-gray-100 w-3/4 rounded-sm"></div>
          <div class="h-8 bg-gray-100 w-1/2 rounded-sm"></div>
        </slot>
      </div>

      <div v-else class="space-y-16">
        <article class="relative">
          <span class="absolute -left-6 -top-4 text-5xl opacity-10 font-serif">“</span>
          <p class="text-2xl md:text-4xl leading-relaxed tracking-wide py-6 text-moxin-ink/90 font-serif">
            {{ sentence.content }}
          </p>
          <span class="absolute -right-4 -bottom-6 text-5xl opacity-10 font-serif">”</span>
        </article>

        <div class="flex justify-end items-center space-x-3 text-base md:text-xl opacity-70">
          <span class="w-10 h-px bg-moxin-ink opacity-30"></span>
          <span class="font-serif">{{ sentence.author }}</span>
          <span v-if="sentence.source" class="font-serif before:content-['《'] after:content-['》']">
            {{ sentence.source }}
          </span>
        </div>
      </div>

      <footer class="mt-20 flex items-center justify-between border-t border-moxin-border/40 pt-8">
        <div class="flex space-x-8 scale-110 origin-left">
          <slot name="actions"></slot>
        </div>

        <div class="text-xs opacity-20 tracking-widest text-moxin-ink font-mono uppercase">
          Daily Inspiration
        </div>
      </footer>
    </section>

    <div class="absolute bottom-8 left-16 opacity-10 text-[10px] font-serif tracking-[0.5em] uppercase">
      Moxin Calendar Project / Est. 2026
    </div>
  </div>
</template>

<script setup>
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
  // 💡 新增：接收可选日期列表
  dateRange: Array
});
defineEmits(['update:selectedDate']);
</script>

<style scoped>
/* 保持纯净的平面质感，取消变换 */
div {
  transition: all 0.5s ease;
}

.no-scrollbar::-webkit-scrollbar { display: none; }
.no-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }
</style>