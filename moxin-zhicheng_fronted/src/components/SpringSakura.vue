<template>
  <canvas 
    ref="sakuraCanvas" 
    class="absolute inset-0 w-full h-full pointer-events-none z-0"
  ></canvas>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue';

// 🌸 1. 引入樱花 SVG 素材 (请确保路径正确)
// 假设你放在 src/assets/sakura.svg
import sakuraSvg from '@/assets/sakura.svg';

const props = defineProps({
  density: {
    type: Number,
    default: 50 // 🌸 樱花建议密度适中，太密会显得杂乱
  }
});

const sakuraCanvas = ref(null);
let ctx = null;
let animationFrameId = null;
let petals = [];
let width = 0;
let height = 0;
let isImgLoaded = ref(false);

// 🌸 2. 预加载图片对象
const sakuraImg = new Image();
sakuraImg.src = sakuraSvg;

// 🌸 3. 樱花瓣粒子类
class SakuraPetal {
  constructor() {
    this.init(true); // true 表示初始生成，分布在全屏
  }

  // 初始化或重置花瓣状态
  init(isFirstLoad = false) {
    this.x = Math.random() * width;
    
    // 如果是第一次加载，随机分布在整个屏幕高度
    // 如果是重置，则从屏幕顶部上方落下
    this.y = isFirstLoad ? Math.random() * height : -20; 
    
    // 🌸 尺寸：樱花瓣通常较小，保持精致感
    this.size = Math.random() * 7 + 5; // 5px ~ 12px
    
    // 🌸 物理特性：樱花极轻，垂直下落速度很慢
    this.velY = Math.random() * 0.4 + 0.3; // 0.3px ~ 0.7px
    
    // 🌸 水平基础微风 (向右微风)
    this.velX = Math.random() * 0.5 + 0.1;   
    
    // 🌸 摆动 (受 Y 轴影响的 sin 波)
    // 幅度越大，左右飘得越远
    this.swayAmptitude = Math.random() * 2 + 1;
    // 频率越低，摆动越平缓
    this.swayFrequency = Math.random() * 0.01 + 0.005;

    // 🌸 旋转与 3D 翻滚
    // 平面旋转角度
    this.angle = Math.random() * Math.PI * 2;
    // 旋转速度 (很慢)
    this.spinSpeed = Math.random() * 0.01 - 0.005; 
    
    // 🌸 核心技巧：模拟 3D 翻转的视觉角度
    this.flipAngle = Math.random() * Math.PI;
    // 翻转速度 (决定花瓣从全貌到侧影的切换快慢)
    this.flipSpeed = Math.random() * 0.02 + 0.01; 
  }

  // 更新位置与状态
  update() {
    this.y += this.velY;
    
    // 🌸 模拟轻盈感：水平位移受正弦波影响明显
    this.x += this.velX + Math.sin(this.y * this.swayFrequency) * this.swayAmptitude;
    
    // 更新角度
    this.angle += this.spinSpeed;
    this.flipAngle += this.flipSpeed;

    // 🌸 边界处理
    // 飘出底部
    if (this.y > height + this.size) {
      this.init(false); // 重置到顶部
    }
    // 飘出左右边界
    if (this.x > width + this.size) {
      this.x = -this.size;
    } else if (this.x < -this.size) {
      this.x = width;
    }
  }

  // 🌸 4. 核心绘制逻辑
  draw() {
    if (!isImgLoaded.value) return; // 确保图片加载完成

    ctx.save();
    
    // 1. 移动画布原点到花瓣中心
    ctx.translate(this.x, this.y);
    
    // 2. 应用平面旋转
    ctx.rotate(this.angle);
    
    // 3. 🌸 应用 3D 翻转视觉效果
    // 使用 Math.abs(Math.cos()) 让 X 轴缩放呈现 0 -> 1 -> 0 的循环
    // 模拟花瓣翻转时变窄的视觉过程
    ctx.scale(Math.abs(Math.cos(this.flipAngle)), 1); 

    // 4. 绘制图片 (注意：原点已移至中心，所以要偏移 -size/2)
    ctx.drawImage(sakuraImg, -this.size/2, -this.size/2, this.size, this.size);
    
    ctx.restore();
  }
}

// 初始化 Canvas 尺寸
const setCanvasSize = () => {
  if (!sakuraCanvas.value) return;
  // 获取父容器的宽高，确保铺满
  const parent = sakuraCanvas.value.parentElement;
  if (parent) {
    width = parent.clientWidth;
    height = parent.clientHeight;
  } else {
    width = window.innerWidth;
    height = window.innerHeight;
  }
  sakuraCanvas.value.width = width;
  sakuraCanvas.value.height = height;
};

// 创建粒子系统
const createPetals = () => {
  petals = [];
  for (let i = 0; i < props.density; i++) {
    petals.push(new SakuraPetal());
  }
};

// 动画主循环
const animate = () => {
  if (!ctx || !isImgLoaded.value) return;
  
  // 1. 清空画布
  ctx.clearRect(0, 0, width, height);

  // 2. 更新并绘制所有花瓣
  petals.forEach(petal => {
    petal.update();
    petal.draw();
  });

  // 3. 请求下一帧
  animationFrameId = requestAnimationFrame(animate);
};

// 窗口大小改变时，重新调整
const handleResize = () => {
  setCanvasSize();
  // 重新生成粒子以适应新尺寸，防止边界溢出
  createPetals(); 
};

onMounted(() => {
  if (!sakuraCanvas.value) return;
  ctx = sakuraCanvas.value.getContext('2d');
  
  setCanvasSize();
  createPetals();
  
  // 🌸 5. 确保图片加载完成才启动动画
  if (sakuraImg.complete) {
    isImgLoaded.value = true;
    animate();
  } else {
    sakuraImg.onload = () => {
      isImgLoaded.value = true;
      animate();
    };
  }

  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  // 组件卸载时清空动画，防止内存泄漏
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId);
  }
  window.removeEventListener('resize', handleResize);
});

// 监听 density 变化，动态调整花瓣数量
watch(() => props.density, () => {
  createPetals();
});
</script>

<style scoped>
/* 确保 Canvas 在父容器中绝对定位 */
canvas {
  display: block;
  backface-visibility: hidden; /* 优化 Canvas 渲染性能 */
}
</style>