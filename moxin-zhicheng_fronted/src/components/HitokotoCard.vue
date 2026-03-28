<template>
  <div 
    v-if="loading" 
    class="mt-12 bg-moxin-paper p-8 border border-moxin-border/70 rounded-sm animate-pulse space-y-3 shadow-lg"
  >
    <div class="h-4 bg-gray-200/50 w-full rounded-sm"></div>
    <div class="h-4 bg-gray-200/50 w-5/6 rounded-sm"></div>
  </div>

  <div 
    v-else 
    class="mt-12 bg-moxin-paper p-1 md:p-1 shadow-[0_10px_30px_rgba(0,0,0,0.2)] relative group transition-all duration-500 hover:-translate-y-1"
  >
    <div class="envelope-border p-6 md:p-10 bg-moxin-paper border border-moxin-border/30 h-full w-full">
      
      <div class="absolute top-10 right-10 w-16 h-16 border-2 border-moxin-shazhu/20 rounded-full flex items-center justify-center rotate-12 pointer-events-none select-none">
        <span class="text-[10px] text-moxin-shazhu/30 font-serif text-center leading-tight uppercase tracking-tighter">
          Moxin <br/> Post
        </span>
      </div>

      <div class="flex items-center space-x-2 mb-8 border-b border-moxin-shazhu/10 pb-2">
        <span class="text-[10px] bg-moxin-shazhu text-white px-2 py-0.5 rounded-sm scale-90">寄</span>
        <h3 class="text-xs font-bold tracking-[0.3em] text-moxin-ink/40 uppercase">见字如面</h3>
      </div>

      <div class="text-base md:text-lg text-moxin-ink/80 leading-relaxed font-serif px-2 min-h-[100px] whitespace-pre-line">
        {{ Hitokoto.hitokoto }}
      </div>

      <div class="mt-10 flex flex-col items-end text-sm text-moxin-ink/50 space-y-1 font-serif">
        <div class="flex items-center">
          <span class="opacity-80 mr-2">来自：</span>
          <span class="text-moxin-ink/70">《{{ Hitokoto.from }}》</span>
        </div>
        <div v-if="Hitokoto.from_who" class="flex items-center">
          <span class="opacity-80 mr-2">笔者：</span>
          <span class="text-moxin-ink/70">{{ Hitokoto.from_who }}</span>
        </div>
      </div>
      
      <div class="absolute bottom-4 left-6 opacity-10">
        <span class="text-[10px] font-mono tracking-tighter">#{{ Math.floor(Math.random()*9000)+1000 }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  Hitokoto: {
    type: Object,
    default: () => ({ hitokoto: '', from: '', from_who: '' })
  },
  loading: {
    type: Boolean,
    default: false
  }
});
</script>

<style scoped>
/* 💡 实现经典的红蓝相间邮封边框 */
.envelope-border {
  background: 
    linear-gradient(white, white) padding-box,
    repeating-linear-gradient(
      -45deg, 
      #B22222 0, #B22222 15px, 
      transparent 15px, transparent 30px, 
      #2C3E50 30px, #2C3E50 45px, 
      transparent 45px, transparent 60px
    ) border-box;
  border: 4px solid transparent;
}
</style>