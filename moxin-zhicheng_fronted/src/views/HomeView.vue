<template>
  <div class="h-screen overflow-y-auto snap-y snap-mandatory scroll-smooth">
   
    <section class="h-screen w-full flex flex-col items-center justify-center snap-start bg-moxin-paper relative">
       <MoxinVFX :density="70" class="absolute inset-0 z-10 vfx-mask"/>
       <!-- <div class="z-20 mb-8 w-full max-w-xl overflow-x-auto no-scrollbar px-4 flex space-x-6 items-center border-b border-moxin-shazhu/10 pb-4">
        <div 
          v-for="dateItem in dateRange" 
          :key="dateItem.full"
          @click="handleDateChange(dateItem.full)"
          class="flex-shrink-0 cursor-pointer transition-all duration-300 flex flex-col items-center"
          :class="selectedDate === dateItem.full ? 'scale-110' : 'opacity-30 hover:opacity-60'"
        >
          <span class="text-[10px] font-mono text-moxin-ink">{{ dateItem.year }}</span>
          <span class="text-lg font-serif text-moxin-ink font-bold" :class="selectedDate === dateItem.full ? 'text-moxin-shazhu' : ''">
            {{ dateItem.display }}
          </span>
        </div>
      </div> -->
      <div class="max-w-3xl w-full px-6">
        <SentenceCard 
    :sentence="dailySentence" 
    :loading="loading"
    :date-range="dateRange"       :model-value="selectedDate"   @update:selected-date="handleDateChange" >
          <template #actions>
            <button @click="handleToggleLike" :class="isLiked ? 'text-moxin-shazhu' : ''">
              {{ isLiked ? '♥ 已藏' : '♡ 收藏' }}
            </button>
          </template>
          <template #date>{{ currentDate }}</template>
        </SentenceCard>
      </div>
      
      <p class="mt-12 text-sm opacity-40 tracking-widest text-moxin-ink">
        素纸一张，写尽心城
      </p>

      <div class="absolute bottom-10 animate-bounce opacity-20">
        <span class="text-4xl">↓</span>
      </div>
    </section>

    <section 
      ref="secondSection"
      class="h-screen w-full flex items-center justify-center snap-start relative shadow-inner overflow-hidden"
    >
      <div 
        class="absolute inset-0 transition-all duration-1000 ease-in-out"
        :class="isSecondSectionActive ? 'blur-sm scale-105' : 'blur-0 scale-100'"
        :style="{ 
          backgroundImage: `url(${tableBg})`, 
          backgroundSize: 'cover', 
          backgroundPosition: 'center' 
        }"
      >
        <div class="absolute inset-0 bg-black/20"></div>
      </div>
      
      <div 
        class="max-w-3xl w-full px-6 z-10 transition-all duration-1000 delay-500 ease-out"
        :class="[
          isSecondSectionActive 
            ? 'opacity-100 translate-y-0 scale-100' 
            : 'opacity-0 translate-y-12 scale-95'
        ]"
      >
        <HitokotoCard :Hitokoto="dailyHitokoto" :loading="loading" />
      </div>
    </section>

  </div>
</template>

<script setup>
import { ref, onMounted, computed,onUnmounted } from 'vue';
import { getDailySentence, getDailyHitokoto } from '../../utils/api/api';
import SentenceCard from '../components/SentenceCard.vue';
import HitokotoCard from '../components/HitokotoCard.vue';

// 💡 导入你放在 assets 里的贴图
import tableBg from '../assets/seelean-iRJH2lSvo_E-unsplash.jpg';
import MoxinVFX from '../components/MoxinVFX.vue'; // 引入特效组件
const loading = ref(true);
const isLiked = ref(false);
const dailySentence = ref({ content: '', author: '', source: '' });
const dailyHitokoto = ref({ hitokoto: '', from: '', from_who: '' });
const secondSection = ref(null); // 💡 必须定义，否则 observer 找不到目标
const isSecondSectionActive = ref(false);
let observer = null;


const currentDate = computed(() => {
  const date = new Date();
  return `${date.getFullYear()} / ${date.getMonth() + 1} / ${date.getDate()}`;
});

const selectedDate = ref(new Date().toISOString().split('T')[0]); // 默认今天

// 💡 修改：让日期范围基于 selectedDate 动态计算
const dateRange = computed(() => {
  const dates = [];
  const baseDate = new Date(selectedDate.value); // 以当前选中日期为基准
  
  // 生成以选中日期为中心的 14 天（前 7 天 到 后 6 天）
  // 这样用户往左点，范围就往左移；往右点，范围就往右移
  for (let i = -7; i < 7; i++) {
    const d = new Date(baseDate);
    d.setDate(d.getDate() + i);
    dates.push({
      full: d.toISOString().split('T')[0],
      year: d.getFullYear(),
      display: `${d.getMonth() + 1}/${d.getDate()}`
    });
  }
  // 如果你希望始终包含今天，或者有特定的排序逻辑，可以对 dates 进行 sort
  return dates;
});

const handleDateChange = (date) => {
  selectedDate.value = date; // 更新选中日期，触发 dateRange 重新计算
  loadData(); 
};

const loadData = async () => {
  loading.value = true;
  try {
    const [sentence, Hitokoto] = await Promise.all([
      getDailySentence(),
      getDailyHitokoto()
    ]);
    dailySentence.value = sentence;
    dailyHitokoto.value = Hitokoto;
  } catch (error) {
    console.error("加载失败", error);
  } finally {
    loading.value = false;
  }
};

const handleToggleLike = () => {
  isLiked.value = !isLiked.value;
};

onMounted(() => {
  loadData();
observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      // 当第二幕进入视口超过 50% 时触发动画
      if (entry.isIntersecting) {
        isSecondSectionActive.value = true;
      } else {
        // 如果想往回滑时重置动画，可以取消下面注释
        // isSecondSectionActive.value = false;
      }
    });
  }, { threshold: 0.5 });

  if (secondSection.value) {
    observer.observe(secondSection.value);
  }
});

onUnmounted(() => {
  if (observer) observer.disconnect();
});
</script>

<style scoped>
/* 隐藏滚动条，保持页面整洁 */
.h-screen::-webkit-scrollbar {
  display: none;
}
.h-screen {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.vfx-mask {
  /* polygon 的逻辑是：先画外圈大矩形，再画内圈小矩形（卡片位置） */
  /* 这里的百分比需要根据你卡片在屏幕上的实际占比进行微调 */
  clip-path: polygon(
    0% 0%, 100% 0%, 100% 100%, 0% 100%, 0% 0%, 
    20% 30%, 80% 30%, 80% 70%, 20% 70%, 20% 30%
  );
}
</style>