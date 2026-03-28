<template>
  <div class="flex items-center justify-center my-12 font-serif select-none z-10 relative">
    
    <button
      @click="emitPage(modelValue - 1)"
      :disabled="modelValue <= 1"
      class="h-10 w-10 flex items-center justify-center rounded-none transition-colors 
             text-[#a38a70] hover:text-[#bc6c25] hover:bg-white/50 
             disabled:opacity-20 disabled:cursor-not-allowed mx-4"
      title="上一页"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
    </button>

    <div class="flex items-center bg-[#fdfcf5] border border-r-0 border-[#d4a373]/40 shadow-sm">
      <template v-for="(page, index) in visiblePages" :key="index">
        <button
          v-if="page !== '...'"
          @click="emitPage(page)"
          :class="[
            'h-10 w-10 flex items-center justify-center text-sm border-r border-[#d4a373]/40 transition-all rounded-none',
            modelValue === page 
              ? 'bg-[#cf4813] text-white border-[#cf4813] font-bold z-10' 
              : 'text-[#6b705c] hover:bg-[#bc6c25]/10'
          ]"
        >
          {{ page }}
        </button>
        
        <span 
          v-else 
          class="h-10 w-10 flex items-center justify-center text-[#6b705c]/50 border-r border-[#d4a373]/40 tracking-widest"
        >
          ...
        </span>
      </template>
    </div>

    <button
      @click="emitPage(modelValue + 1)"
      :disabled="modelValue >= totalPages"
      class="h-10 w-10 flex items-center justify-center rounded-none transition-colors 
             text-[#a38a70] hover:text-[#bc6c25] hover:bg-white/50 
             disabled:opacity-20 disabled:cursor-not-allowed mx-4"
      title="下一页"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
      </svg>
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  modelValue: { type: Number, default: 1 },
  total: { type: Number, default: 0 },
  pageSize: { type: Number, default: 10 }
});

const emit = defineEmits(['update:modelValue', 'change']);

const totalPages = computed(() => Math.ceil(props.total / props.pageSize) || 1);

const visiblePages = computed(() => {
  const current = props.modelValue;
  const last = totalPages.value;
  const delta = 2;
  const range = [];
  const rangeWithDots = [];
  let l;

  for (let i = 1; i <= last; i++) {
    if (i === 1 || i === last || (i >= current - delta && i <= current + delta)) {
      range.push(i);
    }
  }

  for (let i of range) {
    if (l) {
      if (i - l === 2) { rangeWithDots.push(l + 1); } 
      else if (i - l !== 1) { rangeWithDots.push('...'); }
    }
    rangeWithDots.push(i);
    l = i;
  }
  return rangeWithDots;
});

const emitPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    emit('update:modelValue', page);
    emit('change', page);
  }
};
</script>