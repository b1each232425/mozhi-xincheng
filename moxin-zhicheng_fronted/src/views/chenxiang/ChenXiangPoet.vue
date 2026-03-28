<template>
  <div class="flex flex-col min-h-screen">
    <div class="h-[25vh] relative bg-[#00050a] overflow-hidden flex items-center justify-center">
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,#0a192f_0%,#00050a_100%)]"></div>
      
      <div
        @click="router.back()"
        class="absolute top-6 left-6 z-30 text-white/50 text-sm cursor-pointer hover:text-white transition-all flex items-center group"
      >
        <span class="mr-1 transform group-hover:-translate-x-1 transition-transform">←</span>
        返回
      </div>

      <div class="relative z-10 text-center">
        <h1 class="text-white text-3xl font-serif font-bold tracking-[0.2em]">
          {{ poetry.Title || poetry.Rhythmic || '品读中...' }}
        </h1>
        <p class="text-white/60 mt-4 font-serif italic">—— {{ poetry.Author }}</p>
      </div>
    </div>

    <div class="flex-1 relative bg-[#f4f1de] text-[#2b2b2b] py-16">
      <div class="absolute inset-0 pointer-events-none opacity-60 bg-[url('https://www.transparenttextures.com/patterns/handmade-paper.png')]"></div>

      <div class="relative max-w-[700px] mx-auto px-8 z-10">
        <div v-if="loading" class="flex flex-col items-center py-20">
          <div class="w-10 h-10 border-4 border-[#2b2b2b] border-t-transparent rounded-full animate-spin"></div>
        </div>

        <div v-else class="flex flex-col items-center">
          <div class="text-[22px] leading-[2.5] tracking-[0.15em] font-serif text-[#1a1a1a] text-center whitespace-pre-wrap">
            {{ poetry.Paragraphs }}
          </div>

          <div class="w-24 h-px bg-[#d4a373]/40 my-12"></div>

          <div class="flex flex-wrap justify-center gap-3 mb-10">
            <span 
              v-for="tag in (poetry.Tags ? poetry.Tags.split(',') : [])" 
              :key="tag"
              class="text-xs border border-[#d4a373]/30 px-3 py-1 rounded-full text-[#6b705c]"
            >
              # {{ tag }}
            </span>
          </div>

          <Button theme="lotus" shape="pill" size="sm" @click="handleCopy">
            复制全文
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
// 假设你的 api 文件中有获取单首诗的接口，如果没有，可以用 getPoetryList 过滤 ID
import { getPoetryList } from '../../../utils/api/api'; 
import Button from '../../components/Button.vue';

const route = useRoute();
const router = useRouter();
const poetry = ref({});
const loading = ref(true);

const fetchDetail = async () => {
  const id = route.params.id;
  loading.value = true;
  try {
    // 这里演示逻辑：如果后端没有 getById，就利用 keyword=id 搜索
    const res = await getPoetryList(id, 1, 1);
    if (res && res.data.list && res.data.list.length > 0) {
      poetry.value = res.data.list[0];
    }
  } catch (err) {
    console.error("加载详情失败:", err);
  } finally {
    loading.value = false;
  }
};

const handleCopy = () => {
  const text = `${poetry.value.Title}\n${poetry.value.Author}\n\n${poetry.value.Paragraphs}`;
  navigator.clipboard.writeText(text);
  alert("已临摹至剪贴板");
};

onMounted(fetchDetail);
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Noto+Serif+SC:wght@400;700&display=swap");
.font-serif {
  font-family: "Noto Serif SC", serif;
}
</style>