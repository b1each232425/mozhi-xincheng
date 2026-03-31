<template>
  <div class="min-h-screen bg-[#edebe1] flex justify-center items-start py-[6vh] px-[4vw] font-['KaiTi','楷体','serif']">
    <div class="gufeng-paper relative w-full max-w-[900px] min-h-[85vh] bg-[#f6f4ec] shadow-[0_12px_40px_rgba(92,85,75,0.08)] rounded-[2px] overflow-hidden">
      
      <header class="flex justify-between items-center pt-12 px-16 pb-4 bg-transparent">
        <div class="flex items-center">
          <button @click="goBack" class="text-btn">【退回】</button>
          <span class="ml-4 text-[0.9rem] text-[#887f70] border-l border-[#d3cbba] pl-4">
            {{ isNewArticle ? "新卷" : "阅卷" }}
          </span>
        </div>

        <div class="right-action flex items-center gap-4">
          <div v-if="isEditMode || isNewArticle" class="flex items-center gap-2 mr-4">
  <span 
    class="text-xs transition-colors duration-300"
    :class="editForm.is_public ? 'text-[#a63838]' : 'text-[#7a7364]'"
  >
    {{ editForm.is_public ? '公开' : '私语' }}
  </span>
  
  <button 
    @click="editForm.is_public = !editForm.is_public"
    class="relative inline-flex h-5 w-9 items-center rounded-full transition-colors duration-300"
    :class="editForm.is_public ? 'bg-[#a63838]' : 'bg-[#7a7364]'"
  >
    <span 
      class="inline-block h-3 w-3 transform rounded-full bg-white transition-transform duration-300"
      :class="editForm.is_public ? 'translate-x-5' : 'translate-x-1'" 
    />
  </button>
</div>

          <button v-if="!isEditMode" @click="enterEditMode" class="text-btn">【修订】</button>
          <button v-if="isEditMode || isNewArticle" @click="saveArticle" class="seal-button">
            <span>{{ isNewArticle ? '落' : '存' }}</span>
            <span>{{ isNewArticle ? '笔' : '念' }}</span>
          </button>
        </div>
      </header>

      <main class="px-16 py-8">
        <div class="title-section mb-12 text-center">
          <div v-if="isEditMode || isNewArticle" class="max-w-md mx-auto">
            <input v-model="editForm.title" type="text" placeholder="题字..." class="gufeng-title-input" />
            <div class="title-underline"></div>
          </div>
          <div v-else>
            <h1 class="text-[2.2rem] font-normal text-[#1a1a1a] mb-4 tracking-[4px]">{{ articleData?.title }}</h1>
            <span class="text-[0.9rem] text-[#7a7364]">{{ formatDate(articleData?.CreatedAt) }}</span>
          </div>
        </div>

        <div class="content-section min-h-[400px]">
          <div v-if="isEditMode || isNewArticle" class="editor-container">
            <Toolbar
              style="border-bottom: 1px solid #d3cbba"
              :editor="editorRef"
              :defaultConfig="toolbarConfig"
              mode="simple"
            />
            <Editor
              style="height: 500px; overflow-y: hidden;"
              v-model="valueHtml"
              :defaultConfig="editorConfig"
              mode="simple"
              @onCreated="handleCreated"
            />
          </div>
          <div v-else class="gufeng-content" v-html="articleData?.content"></div>
        </div>

        <footer v-if="!isEditMode" class="mt-20 flex flex-col items-center opacity-30">
          <div class="w-[1px] h-12 bg-[#5c554b] mb-4"></div>
          <p class="text-xs tracking-[4px]">终</p>
        </footer>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, shallowRef, onBeforeUnmount, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Editor, Toolbar } from '@wangeditor/editor-for-vue';
import '@wangeditor/editor/dist/css/style.css';
import { getLuoBiArticle, createLuoBiArticle } from '../../../utils/api/api';

const route = useRoute();
const router = useRouter();

// 状态控制
const isEditMode = ref(false);
const isNewArticle = computed(() => route.params.id === "0");
const articleData = ref(null);

// wangEditor 逻辑
const editorRef = shallowRef(); // 必须用 shallowRef
const valueHtml = ref('');
const toolbarConfig = { excludeKeys: ['fullScreen', 'insertVideo','codeBlock','insertImage',] };
const editorConfig = { placeholder: '展纸研墨...', MENU_CONF: {} };

const editForm = ref({
  title: "",
  is_public: true
});

const handleCreated = (editor) => { editorRef.value = editor; };

// 销毁编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value;
  if (editor) editor.destroy();
});

const fetchDetail = async () => {
  if (isNewArticle.value) {
    isEditMode.value = true;
  } else {
    const res = await getLuoBiArticle(route.params.id);
    if (res?.code === 200) {
      articleData.value = res.data;
      valueHtml.value = res.data.content;
      editForm.value.title = res.data.title;
      editForm.value.is_public = res.data.is_public === 1;
    }
  }
};

const saveArticle = async () => {
  if (!editForm.value.title.trim()) return alert("请先题写标题");
  
  // 获取纯文本做摘要
  const plainText = editorRef.value.getText();
  const summary = plainText.slice(0, 30);

  const res = await createLuoBiArticle({
    title: editForm.value.title,
    content: valueHtml.value, // 富文本 HTML
    summary: summary,
    is_public: editForm.value.is_public
  });

  if (res?.code === 200) {
    alert("已存念");
    router.push("/luobi");
  }
};

const formatDate = (d) => d ? new Date(d).toLocaleDateString() : '';
const goBack = () => router.push("/luobi");
const enterEditMode = () => { isEditMode.value = true; };

onMounted(fetchDetail);
</script>

<style scoped>
@import "tailwindcss";
.text-btn { @apply bg-none border-none text-[1.1rem] text-[#5c554b] cursor-pointer hover:text-[#a63838] transition-colors; }

.seal-button {
  @apply flex flex-col items-center justify-center w-12 h-12 border-2 border-[#a63838] text-[#a63838] 
  font-bold leading-tight rounded-[2px] hover:bg-[#a63838] hover:text-white transition-all active:scale-95;
}

.gufeng-title-input {
  @apply w-full bg-transparent border-none text-center text-2xl tracking-[4px] outline-none text-[#1a1a1a];
}

.title-underline { @apply w-24 h-[1px] bg-[#d3cbba] mx-auto mt-2; }

.gufeng-content {
  @apply text-[1.15rem] leading-[2.6] text-[#2b2b2b] text-justify;
  text-indent: 2.3rem;
}

/* 编辑器样式覆盖，使其融入背景 */
:deep(.w-e-text-container) { background-color: transparent !important; border: none !important; }
:deep(.w-e-toolbar) { background-color: transparent !important; border: none !important; }
</style>