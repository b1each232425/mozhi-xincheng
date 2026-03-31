<template>
  <canvas 
    ref="vfxCanvas" 
    class="absolute inset-0 w-full h-full pointer-events-none z-0"
  ></canvas>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue';

// 💡 1. 引入你的 SVG 素材 (请确保路径正确)
// 假设你放在 src/assets/maple-leaf.svg
import mapleSvg from '../assets/maple-leaf.svg';

const props = defineProps({
  density: {
    type: Number,
    default: 70 // 💡 改进：默认枫叶数量多一点，营造落叶满天感
  }
});

const vfxCanvas = ref(null);
let ctx = null;
let animationFrameId = null;
let leaves = [];
let width = 0;
let height = 0;
let isImgLoaded = ref(false);

// 💡 2. 预加载 SVG 图片对象 (不需动态染色)
const leafImg = new Image();
leafImg.src = mapleSvg;

// 💡 粒子类 (枫叶原型)
class MapleLeaf {
  constructor() {
    this.init();
  }

  init() {
    this.x = Math.random() * width; // 随机 X 轴
    // 💡 改进：初始化分布在更广的高度范围内，减少组件刚加载时的堆叠感
    this.y = Math.random() * height * -1.5; 
    
    // 💡 形状与大小 (适当缩小，保持精致感)
    this.size = Math.random() * 10 + 10; // 10px ~ 20px
    
    // 💡 3. 改进：动态属性 (更慢、更灵动)
    // 基础下落速度大大降低 (0.2px ~ 0.5px 之间)
    this.velY = Math.random() * 0.3 + 0.2; 
    // 左右基础飘移也相应降低
    this.velX = Math.random() * 0.6 - 0.3;   
    
    // 水平摆动的幅度系数
    this.swayAmptitude = Math.random() * 1.5 + 0.5;
    // 水平摆动的频率系数
    this.swayFrequency = Math.random() * 0.01 + 0.005;

    // 💡 核心：旋转与翻转 (也调慢阻尼)
    this.angle = Math.random() * Math.PI * 2; // 初始角度
    this.spinSpeed = Math.random() * 0.006 - 0.003;  // 💡 旋转速度调得极慢
    
    // 💡 改进：优化翻转动画原理
    this.flipAngle = Math.random() * Math.PI;
    this.flipSpeed = Math.random() * 0.01 + 0.005; // 💡 翻转速度调极慢
  }

  // 更新位置与状态
  update() {
    this.y += this.velY;
    
    // 💡 4. 改进：增强水平摆动 (随 y 值变化的 sin 波)
    this.x += this.velX + Math.sin(this.y * this.swayFrequency) * this.swayAmptitude; 
    
    // 更新旋转和翻转
    this.angle += this.spinSpeed;
    this.flipAngle += this.flipSpeed;

    // 💡 边界处理 (飘出屏幕后，从顶部重新落下)
    if (this.y > height + this.size) {
      this.y = -this.size * 2;
      this.x = Math.random() * width;
    }
    if (this.x > width + this.size || this.x < -this.size) {
      this.x = Math.random() * width;
    }
  }

  // 💡 5. 改进：核心绘制逻辑 (直接绘制，无需动态染色)
  draw() {
    if (!isImgLoaded.value) return; // 确保图片加载完成

    ctx.save(); // 保存当前画布状态
    
    // 1. 移动画布原点到叶片中心
    ctx.translate(this.x, this.y);
    // 应用极其缓慢的打转旋转
    ctx.rotate(this.angle);
    
    // 2. 应用真实翻转动画 (模拟立体叶片变窄)
    // 使用 Math.abs 控制 scaleX 在 0~1 之间变化
    ctx.scale(Math.abs(Math.cos(this.flipAngle)), 1); 

    // 3. 💡 改进：直接绘制图片素材，性能最高且颜色精准
    ctx.drawImage(leafImg, -this.size/2, -this.size/2, this.size, this.size);
    
    ctx.restore(); // 恢复画布状态
  }
}

// 初始化 Canvas 尺寸
const setCanvasSize = () => {
  if (!vfxCanvas.value) return;
  // 直接使用 window 宽高，确保全屏
  width = window.innerWidth;
  height = window.innerHeight;
  vfxCanvas.value.width = width;
  vfxCanvas.value.height = height;
};

// 创建粒子系统
const createLeaves = () => {
  leaves = [];
  for (let i = 0; i < props.density; i++) {
    leaves.push(new MapleLeaf());
  }
};

// 动画主循环
const animate = () => {
  if (!ctx || !isImgLoaded.value) return;
  // 1. 清空画布
  ctx.clearRect(0, 0, width, height);

  // 2. 更新并绘制所有叶片
  leaves.forEach(leaf => {
    leaf.update();
    leaf.draw();
  });

  // 3. 请求下一帧
  animationFrameId = requestAnimationFrame(animate);
};

// 窗口大小改变时，重新调整
const handleResize = () => {
  setCanvasSize();
  createLeaves(); // 重新生成粒子以适应新尺寸
};

onMounted(() => {
  if (!vfxCanvas.value) return;
  ctx = vfxCanvas.value.getContext('2d');
  
  setCanvasSize();
  createLeaves();
  
  // 💡 6. 改进：只有图片加载完成才启动动画循环，防止 ctx 报错
  if (leafImg.complete) {
    isImgLoaded.value = true;
    animate();
  } else {
    leafImg.onload = () => {
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

// 监听 density 变化
watch(() => props.density, () => {
  createLeaves();
});
</script>