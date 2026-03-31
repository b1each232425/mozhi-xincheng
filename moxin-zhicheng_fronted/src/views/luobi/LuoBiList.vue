<template>
  <div class="max-w-4xl mx-auto p-6 md:p-10 font-serif text-moxin-ink selection:bg-moxin-shazhu/20">
    
    <header class="mb-12 border-b border-moxin-border/50 pb-10 flex flex-col md:flex-row md:items-end justify-between gap-6">
      <div class="space-y-4">
        <h1 class="text-3xl font-bold tracking-[0.4em] text-moxin-ink">落笔</h1>
        <div class="space-y-1">
          <p class="text-base text-moxin-ink/80 tracking-[0.2em]">闲写抛书林下坐，有时随事记流年。</p>
          <p class="text-[10px] text-moxin-ink/40 tracking-[0.3em] uppercase italic">Mud on the snow, traces of the soul</p>
        </div>
      </div>
      
      <button @click="navigateToCreate" class="px-6 py-2.5 bg-moxin-shazhu text-white text-sm rounded-full shadow-sm hover:bg-moxin-shazhu/90 transition-all active:scale-95 tracking-widest flex items-center space-x-2 self-start md:self-auto">
        <span>落笔</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-pen-line"><path d="M12 20h9"/><path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z"/></svg>
      </button>
    </header>

    <div class="space-y-8">
      <article 
        v-for="article in articles" 
        :key="article.ID"
        class="group bg-white/40 backdrop-blur-sm p-8 rounded-2xl border border-moxin-border/30 shadow-sm hover:shadow-md hover:border-moxin-shazhu/30 transition-all duration-300 cursor-pointer relative overflow-hidden"
        @click="viewDetail(article.ID)"
      >
        <div class="absolute inset-0 bg-gradient-to-br from-moxin-shazhu/0 via-transparent to-moxin-shazhu/5 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>

        <div class="relative z-10 flex flex-col md:flex-row md:items-start md:justify-between gap-6">
          <div class="flex-1">
            <h2 class="text-xl font-bold tracking-wider text-moxin-ink group-hover:text-moxin-shazhu transition-colors">
              {{ article.title }}
            </h2>
            
            <p class="mt-4 text-moxin-ink/70 leading-relaxed text-sm tracking-wide line-clamp-2 md:line-clamp-3">
              {{ article.summary }}
            </p>

            <div class="mt-6 flex items-center gap-6 text-xs text-moxin-ink/50 tracking-widest border-t border-moxin-border/30 pt-4">
              <span class="flex items-center gap-1.5">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M8 2v4"/><path d="M16 2v4"/><rect width="18" height="18" x="3" y="4" rx="2"/><path d="M3 10h18"/></svg>
                {{ formatDate(article.CreatedAt) }}
              </span>
              <span class="flex items-center gap-1.5">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"/><circle cx="12" cy="12" r="3"/></svg>
                {{ article.view_count }} 查看
              </span>
            </div>
          </div>
          
          <div class="shrink-0 flex items-center gap-3">
  <span 
    class="text-xs tracking-widest transition-colors duration-300"
    :class="article.is_public === 1 ? 'text-moxin-shazhu' : 'text-moxin-ink/40'"
  >
    {{ article.is_public === 1 ? '公开' : '私语' }}
  </span>

  <button 
    @click.stop="togglePublicStatus(article)"
    class="relative inline-flex h-5 w-9 items-center rounded-full transition-colors duration-300 focus:outline-none"
    :class="article.is_public === 1 ? 'bg-moxin-shazhu' : 'bg-moxin-ink/10'"
  >
    <span
      class="inline-block h-3 w-3 transform rounded-full bg-white transition-transform duration-300"
      :class="article.is_public === 1 ? 'translate-x-5' : 'translate-x-1'"
    />
  </button>
</div>
        </div>
      </article>
      
      <div v-if="articles.length === 0" class="py-20 text-center text-moxin-ink/30 tracking-widest">
        暂无墨迹，待君落笔
      </div>
    </div>

    <footer class="mt-16 text-center text-sm text-moxin-ink/40 tracking-[0.2em]">
      <p>仅显示最近的墨迹</p>
    </footer>

  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { getLuoBiList } from '../../../utils/api/api';

const router = useRouter();
const articles = ref([]); // 存放真实数据

const viewDetail = (id) => {
  router.push(`/luobi/article/${id}`);
};

const navigateToCreate = () => {
  router.push('/luobi/article/0');
};

// 格式化日期函数
const formatDate = (dateStr) => {
  if (!dateStr) return '';
  return new Date(dateStr).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  });
};

const fetchArticles = async () => {
  try {
    const res = await getLuoBiList(); // 调用后端接口
    if (res && res.code === 200) {
      articles.value = res.data.list; // 绑定后端返回的列表数据
    }
  } catch (error) {
    console.error('获取墨迹列表失败:', error);
  }
};

const togglePublicStatus = async (article) => {
  // 切换本地状态 (0 变 1, 1 变 0)
  const newStatus = article.is_public === 1 ? 0 : 1;
  
  // 如果你有更新接口，可以在这里调用
  // try {
  //   await updateArticleStatus(article.ID, { is_public: newStatus });
  //   article.is_public = newStatus;
  // } catch (err) {
  //   alert("同步状态失败");
  // }

  // 演示直接修改
  article.is_public = newStatus;
};

onMounted(fetchArticles);
</script>