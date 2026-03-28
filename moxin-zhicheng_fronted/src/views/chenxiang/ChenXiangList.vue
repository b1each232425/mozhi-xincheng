<template>
  <div class="flex flex-col min-h-screen">
    <div
      class="h-[35vh] relative bg-[#00050a] overflow-hidden flex items-center justify-center"
    >
      <div
        class="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,#0a192f_0%,#00050a_100%)]"
      >
        <div
          @click="router.push('/chenxiang')"
          class="absolute top-6 left-6 z-20 text-white/50 text-md cursor-pointer hover:text-white transition-all flex items-center group"
        >
          <span
            class="mr-1 transform group-hover:-translate-x-1 transition-transform"
            >←</span
          >
          返回星海
        </div>
      </div>
      <div
        class="absolute -top-[10%] -right-[5%] w-[40vw] h-[40vw] blur-[60px] opacity-15 rounded-full bg-[radial-gradient(circle,#4b0082,transparent)]"
      ></div>
      <div
        class="absolute -bottom-[10%] -left-[5%] w-[40vw] h-[40vw] blur-[60px] opacity-15 rounded-full bg-[radial-gradient(circle,#0000ff,transparent)]"
      ></div>

      <div class="relative z-10 text-center w-full px-4">
        <div class="flex justify-center gap-[10px]">
          <input
            type="text"
            v-model="newQuery"
            @keyup.enter="handleNewSearch"
            placeholder="更深处寻觅..."
            class="w-full max-w-[300px] px-5 py-[10px] bg-white/10 border border-white/30 rounded-[20px] text-white outline-none backdrop-blur-md focus:border-white/60 transition-all"
          />
          <Button
            theme="cinnabar"
            shape="pill"
            customClass="h-[44px] px-6 text-sm tracking-widest"
            @click="handleNewSearch"
          >
            探寻
          </Button>
        </div>
      </div>
    </div>

    <div class="flex-1 relative bg-[#f4f1de] text-[#2b2b2b] py-10">
      <div
        class="absolute inset-0 pointer-events-none opacity-60 bg-[url('https://www.transparenttextures.com/patterns/handmade-paper.png')]"
      ></div>

      <div class="relative max-w-[800px] mx-auto px-5 z-10">
        <div
          class="mb-10 font-serif border-b-2 border-[#d4a373]/50 pb-5 text-[#6b705c] text-center"
        >
          已从星海中为你打捞关于 “<span
            class="text-[#bc6c25] font-bold text-lg"
            >{{ route.query.keyword }}</span
          >” 的星光
        </div>

        <div v-if="loading" class="flex flex-col items-center py-16">
          <div
            class="w-12 h-12 border-4 border-[#2b2b2b] border-t-transparent rounded-full animate-spin mb-5 opacity-70"
          ></div>
          <p class="font-serif italic text-lg text-[#6b705c] opacity-80">
            正在研墨...
          </p>
        </div>

        <div v-else-if="poems && poems.length > 0" class="space-y-12">
          <div
            v-for="(item, index) in poems"
            :key="item.ID || index"
            @click="router.push(`/Chenxiang/Poet/${item.ID}`)"
            class="relative group p-10 rounded-xl border border-[#d4a373]/30 bg-[#fdfcf5]/80 backdrop-blur-sm shadow-[5px_5px_15px_rgba(0,0,0,0.02)] transition-all duration-500 ease-out hover:-translate-y-2 hover:shadow-[10px_10px_30px_rgba(188,108,37,0.07)]"
          >
            <div
              class="absolute inset-0 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-700 pointer-events-none bg-[radial-gradient(circle_at_0%_0%,#d4a37310,transparent_20%),radial-gradient(circle_at_100%_0%,#d4a37310,transparent_20%),radial-gradient(circle_at_0%_100%,#d4a37310,transparent_20%),radial-gradient(circle_at_100%_100%,#d4a37310,transparent_20%)]"
            ></div>

            <h2
              class="text-[28px] mb-[18px] font-serif font-bold text-[#1a1a1a] cursor-pointer group-hover:text-[#bc6c25] transition-colors duration-300 relative inline-block"
            >
              {{ item.Title || item.Rhythmic }}
              <span
                class="absolute bottom-0 left-0 w-0 h-0.5 bg-[#bc6c25]/40 group-hover:w-full transition-all duration-500"
              ></span>
            </h2>

            <div
              class="flex items-center text-[#5a6378] mb-6 text-[15px] font-serif opacity-90"
            >
              <!-- <span class="mr-3 text-[#bc6c25]/80">[{{ item.Dynasty || '未知' }}]</span> -->
              <span class="font-medium pr-4 border-r border-[#d4a373]/30">{{
                item.Author
              }}</span>
              <span class="pl-4 text-xs opacity-60">ID: {{ item.ID }}</span>
            </div>

            <div
              class="text-[20px] leading-[2.1] tracking-wider font-serif whitespace-pre-wrap text-[#333] opacity-95"
            >
              {{ formatContent(item.Paragraphs) }}
            </div>

            <div class="mt-8 flex flex-wrap gap-2.5" v-if="item.Tags">
              <span
                v-for="tag in item.Tags.split(',')"
                :key="tag"
                class="text-xs bg-[#f4f1de]/60 text-[#6b705c] px-3 py-1.5 rounded-full border border-[#d4a373]/30 hover:bg-[#bc6c25]/10 hover:border-[#bc6c25]/40 transition-colors cursor-default"
              >
                # {{ tag }}
              </span>
            </div>
          </div>

          <ZenPagination
            v-model="currentPage"
            :total="total"
            v-model:pageSize="pageSize"
            @change="fetchData"
          />

          <div
            class="pt-10 text-center text-[#6b705c] font-serif border-t border-[#d4a373]/40 mt-16 opacity-80"
          >
            共在星海中寻得 {{ total }} 首相关诗词
          </div>
        </div>

        <div
          v-else
          class="text-center py-24 rounded-2xl bg-[#fdfcf5]/50 border border-[#d4a373]/20"
        >
          <p class="text-xl font-serif text-[#6b705c] opacity-70">
            星海茫茫，此处尚未寻得星光。
          </p>
          <p class="mt-4 text-sm text-[#bc6c25]/60">也许可以换个词试试？</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getPoetryList as searchApi } from "../../../utils/api/api";
import Button from "../../components/Button.vue";
import ZenPagination from "../../components/ZenPagination.vue";
const route = useRoute();
const router = useRouter();
const newQuery = ref(route.query.keyword || "");
const poems = ref([]);
const total = ref(0);
const loading = ref(false);
const currentPage = ref(Number(route.query.page) || 1);
const pageSize = ref(10);

/**
 * 格式化诗词内容：截取前两句
 */
const formatContent = (content) => {
  if (!content) return "";
  const lines = content
    .split(/[\n,，。！？]/)
    .filter((line) => line.trim() !== "");
  if (lines.length > 2) {
    return `${lines[0]}，${lines[1]} ...`;
  }
  return content;
};

/**
 * 获取后端数据
 */
const fetchData = async () => {
  const keyword = route.query.keyword;
  if (!keyword) return;

  loading.value = true;
  try {
    const res = await searchApi(keyword, currentPage.value, pageSize.value);
    if (res && res.code === 200) {
      poems.value = res.data.list || [];
      total.value = res.data.total || 0;
    } else {
      poems.value = [];
      total.value = 0;
    }
  } catch (err) {
    console.error("加载失败:", err);
    poems.value = [];
  } finally {
    loading.value = false;
  }
};

/**
 * 触发新搜索
 */
const handleNewSearch = () => {
  if (!newQuery.value.trim()) return;
  router.push({ query: { keyword: newQuery.value } });
};

onMounted(fetchData);

watch(
  () => [route.query.keyword, route.query.page],
  ([nextK, nextP]) => {
    newQuery.value = nextK || "";
    currentPage.value = Number(nextP) || 1;
    fetchData();
  },
);
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Noto+Serif+SC:wght@400;700&display=swap");

.font-serif {
  font-family: "Noto Serif SC", "Source Han Serif SC", "SimSun", serif;
}
</style>
