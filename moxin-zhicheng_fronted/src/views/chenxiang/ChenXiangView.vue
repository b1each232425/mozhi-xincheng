<template>
  <div 
    class="fixed inset-0 bg-[#00050a] overflow-hidden cursor-crosshair"
    @mousemove="handleMouseMove" 
    @click="handleCanvasClick"
  >
    <div class="absolute inset-0 z-[1] bg-[radial-gradient(circle_at_50%_50%,#0a192f_0%,#00050a_100%)]"></div>
    
    <div class="absolute -top-[10%] -right-[10%] w-[60vw] h-[60vw] rounded-full blur-[80px] opacity-15 pointer-events-none bg-[radial-gradient(circle,#4b0082,transparent)]"></div>
    <div class="absolute -bottom-[10%] -left-[10%] w-[60vw] h-[60vw] rounded-full blur-[80px] opacity-15 pointer-events-none bg-[radial-gradient(circle,#0000ff,transparent)]"></div>

    <canvas ref="canvasRef" class="absolute inset-0 z-[2] block"></canvas>

    <div 
      v-if="hoveredStar" 
      class="absolute pointer-events-none bg-white/10 backdrop-blur-md px-4 py-2 rounded border border-white/20 text-white z-[100] animate-in fade-in slide-in-from-top-1 duration-300"
      :style="{ left: hoveredStar.x + 15 + 'px', top: hoveredStar.y - 20 + 'px' }"
    >
      <span class="text-base font-bold tracking-[2px] block">{{ hoveredStar.label }}</span>
      <div class="text-[10px] opacity-60 text-center">点击探寻</div>
    </div>

    <div class="absolute inset-0 flex justify-center items-center pointer-events-none z-[3]">
      <div class="text-center pointer-events-auto">
        <div class="relative group">
          <input 
            type="text" 
            placeholder="众里寻他千百度..." 
            class="w-[400px] px-6 py-[15px] bg-white/5 border border-white/20 rounded-[30px] color-white text-lg outline-none backdrop-blur-[15px] transition-all duration-300 focus:w-[450px] focus:bg-white/10 focus:border-white/50 text-white"
            v-model="query" 
            @keyup.enter="handleSearch"
          />
          <div class="absolute inset-0 -z-10 bg-blue-500/10 blur-xl opacity-0 group-focus-within:opacity-100 transition-opacity"></div>
        </div>
        <div class="mt-5 text-white/40 text-xs">已在星海中收录 345,618 片星光</div>
      </div>
    </div>
  </div>
</template>

<script setup>
/* JS 逻辑保持不变 */
import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { getStars } from '../../../utils/api/api'; 

const router = useRouter();
const query = ref('');
const canvasRef = ref(null);
const hoveredStar = ref(null);

let ctx = null;
let stars = [];
let functionalStars = [];
let animationId = null;
let mouseX = -1000, mouseY = -1000;

class Star {
  constructor(w, h, isFunctional = false, label = '') {
    this.x = Math.random() * w;
    this.y = Math.random() * h;
    this.isFunctional = isFunctional;
    this.label = label;
    this.size = isFunctional ? Math.random() * 2 + 3 : Math.random() * 1.5;
    this.baseOpacity = Math.random() * 0.5 + 0.2;
    this.opacity = this.baseOpacity;
    this.blinkSpeed = Math.random() * 0.03 + 0.01;
    this.angle = Math.random() * Math.PI * 2;
    this.scale = 1;
    this.targetScale = 1;
    this.color = isFunctional ? (Math.random() > 0.5 ? '#FFD700' : '#87CEFA') : '#FFFFFF';
  }

  update(mX, mY) {
    this.angle += this.blinkSpeed;
    this.opacity = this.baseOpacity + Math.sin(this.angle) * 0.2;
    if (this.isFunctional) {
      const dist = Math.hypot(this.x - mX, this.y - mY);
      if (dist < 50) {
        this.targetScale = 1.8;
        this.x += (mX - this.x) * 0.02;
        this.y += (mY - this.y) * 0.02;
      } else {
        this.targetScale = 1;
      }
      this.scale += (this.targetScale - this.scale) * 0.1;
    }
  }

  draw() {
    ctx.save();
    ctx.fillStyle = this.color;
    ctx.globalAlpha = this.opacity;
    if (this.isFunctional) {
      ctx.shadowBlur = 15 * this.scale;
      ctx.shadowColor = this.color;
    }
    ctx.beginPath();
    ctx.arc(this.x, this.y, this.size * this.scale, 0, Math.PI * 2);
    ctx.fill();
    ctx.restore();
  }
}

const initCanvas = (tags) => {
  const canvas = canvasRef.value;
  if (!canvas) return;
  ctx = canvas.getContext('2d');
  const w = canvas.width = window.innerWidth;
  const h = canvas.height = window.innerHeight;
  stars = [];
  for (let i = 0; i < 600; i++) stars.push(new Star(w, h, false));
  functionalStars = tags.map(tag => {
    const s = new Star(w, h, true, tag);
    s.x = w * 0.2 + Math.random() * w * 0.6;
    s.y = h * 0.2 + Math.random() * h * 0.6;
    return s;
  });
  stars.push(...functionalStars);
};

const fetchAndStart = async () => {
  try {
    const res = await getStars(); 
    const finalTags = (res && res.code === 200) ? res.data : ["唐诗", "宋词", "苏轼", "李白"];
    initCanvas(finalTags);
    render();
  } catch (err) {
    initCanvas(["系统繁忙"]);
    render();
  }
};

const render = () => {
  if (!ctx) return;
  ctx.clearRect(0, 0, canvasRef.value.width, canvasRef.value.height);
  stars.forEach(s => { s.update(mouseX, mouseY); s.draw(); });
  animationId = requestAnimationFrame(render);
};

const handleMouseMove = (e) => {
  mouseX = e.clientX;
  mouseY = e.clientY;
  const active = functionalStars.find(s => Math.hypot(s.x - mouseX, s.y - mouseY) < 30);
  hoveredStar.value = active || null;
};

const handleSearch = () => {
  if (!query.value.trim()) return;
  router.push({ name: 'ChenXiangList', query: { keyword: query.value } });
};

const handleCanvasClick = () => {
  if (hoveredStar.value) {
    router.push({ name: 'ChenXiangList', query: { keyword: hoveredStar.value.label } });
  }
};

onMounted(() => {
  fetchAndStart();
  window.addEventListener('resize', () => {
    const currentTags = functionalStars.map(s => s.label);
    initCanvas(currentTags);
  });
});

onUnmounted(() => {
  cancelAnimationFrame(animationId);
});
</script>

<style scoped>
@import "tailwindcss";
</style>