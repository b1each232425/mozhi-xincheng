<template>
  <div class="flex flex-col min-h-screen">
    <div class="h-[30vh] relative bg-[#00050a] overflow-hidden flex items-center justify-center">
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,#0a192f_0%,#00050a_100%)]"></div>
      
      <div
        @click="router.back()"
        class="absolute top-6 left-6 z-30 text-white/50 text-sm cursor-pointer hover:text-white transition-all flex items-center group"
      >
        <span class="mr-1 transform group-hover:-translate-x-1 transition-transform">←</span>
        返回星海
      </div>

      <div class="relative z-10 text-center px-4" v-if="!loading">
        <h1 class="text-white text-3xl md:text-4xl font-serif font-bold tracking-[0.2em] mb-4">
          {{ poetry.Title || poetry.Rhythmic }}
        </h1>
        <div class="flex items-center justify-center gap-4 text-white/70 font-serif italic">
          <span v-if="poetry.Dynasty">〔{{ poetry.Dynasty }}〕</span>
          <span>{{ poetry.Author }}</span>
        </div>
      </div>
    </div>

    <div class="flex-1 relative bg-[#f4f1de] text-[#2b2b2b] py-16">
      <div class="absolute inset-0 pointer-events-none opacity-60 bg-[url('https://www.transparenttextures.com/patterns/handmade-paper.png')]"></div>

      <div class="relative max-w-[800px] mx-auto px-8 z-10">
        <div v-if="loading" class="flex flex-col items-center py-20">
          <div class="w-10 h-10 border-4 border-[#2b2b2b] border-t-transparent rounded-full animate-spin"></div>
          <p class="mt-4 font-serif italic text-[#6b705c]">正在为你铺展画卷...</p>
        </div>

        <div v-else-if="poetry.ID" class="flex flex-col">
          <section class="flex flex-col items-center mb-16">
            <div class="text-[22px] md:text-[25px] leading-[2.6] tracking-[0.18em] font-serif text-[#1a1a1a] text-center whitespace-pre-wrap">
              {{ poetry.Paragraphs }}
            </div>
            
            <div class="flex flex-wrap justify-center gap-3 mt-10" v-if="poetry.Tags">
              <span 
                v-for="tag in poetry.Tags.split(',')" 
                :key="tag"
                class="text-xs border border-[#d4a373]/40 px-3 py-1 rounded-full text-[#6b705c] bg-[#fdfcf5]/50"
              >
                # {{ tag }}
              </span>
            </div>
          </section>

          <div class="flex items-center justify-center gap-4 my-8">
            <div class="w-12 h-px bg-[#d4a373]/30"></div>
            <div class="w-12 h-px bg-[#d4a373]/30"></div>
          </div>

          <section v-if="poetry.Translation" class="mb-12">
            <h3 class="font-serif font-bold text-[#bc6c25] text-lg mb-4 flex items-center">
              <span class="w-1 h-5 bg-[#bc6c25] mr-3 rounded-full"></span>
              今译
            </h3>
            <div class="bg-[#fdfcf5]/80 p-7 rounded-2xl border border-[#d4a373]/20 shadow-sm leading-8 text-[17px] text-[#444] font-serif italic">
              {{ formatAIText(poetry.Translation) }}
            </div>
          </section>

          <section v-if="poetry.Annotation" class="mb-12">
            <h3 class="font-serif font-bold text-[#6b705c] text-lg mb-4 flex items-center">
              <span class="w-1 h-5 bg-[#6b705c] mr-3 rounded-full"></span>
              注释
            </h3>
            <div class="bg-[#fdfcf5]/60 p-7 rounded-2xl border border-[#d4a373]/20 shadow-sm">
              <div class="text-[16px] leading-8 text-[#555] font-serif whitespace-pre-wrap">
                {{ poetry.Annotation }}
              </div>
            </div>
          </section>

        </div>

        <div v-else class="text-center py-20 font-serif text-[#6b705c]">
          未能寻得此中星光
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getPoemDetail } from '../../../utils/api/api'; 
import Button from '../../components/Button.vue';

const route = useRoute();
const router = useRouter();
const poetry = ref({});
const loading = ref(true);

/**
 * 核心逻辑：从路径中获取 ID 并请求后端 SinglePoem 接口
 */
const fetchDetail = async () => {
  const id = route.params.id; // 匹配路由配置中的 /poet/:id
  if (!id) return;
  
  loading.value = true;
  try {
    // 对应后端 SinglePoem (c.Param("id"))
    const res = await getPoemDetail(id);
    if (res && res.code === 200) {
      poetry.value = res.data; // 直接接收后端返回的 model.Poetry 对象
    }
  } catch (err) {
    console.error("获取诗词详情失败:", err);
  } finally {
    loading.value = false;
  }
};

/**
 * 处理 AI 字段标记
 */
const formatAIText = (text) => {
  if (!text) return "";
  return text.replace("[AI补全] ", "");
};

const handleCopy = () => {
  const text = `《${poetry.value.Title || poetry.value.Rhythmic}》\n${poetry.value.Author}\n\n${poetry.value.Paragraphs}`;
  navigator.clipboard.writeText(text);
  alert("已将诗文临摹至剪贴板");
};

onMounted(fetchDetail);
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Noto+Serif+SC:wght@400;700&display=swap");
.font-serif {
  font-family: "Noto Serif SC", serif;
}
</style>