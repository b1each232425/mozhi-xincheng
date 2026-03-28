<template>
  <button
    :class="[
      'inline-flex items-center justify-center font-serif transition-all duration-300 select-none cursor-pointer',
      'focus:outline-none focus:ring-2 focus:ring-offset-2',
      shapeClasses,
      themeClasses,
      sizeClasses,
      customClass,
      { 'opacity-50 cursor-not-allowed': disabled }
    ]"
    :disabled="disabled"
    @click="handleClick"
  >
    <slot name="icon-left"></slot>
    
    <span :class="{ 'mx-2': $slots['icon-left'] || $slots['icon-right'] }">
      <slot></slot>
    </span>

    <slot name="icon-right"></slot>
  </button>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  // 1. 主题色调 ( cinnabar | chocolate | brick | sprout | lotus )
  theme: { type: String, default: 'cinnabar' },
  // 2. 形状 ( rounded | rect | pill | circle )
  shape: { type: String, default: 'rounded' },
  // 3. 尺寸 ( sm | md | lg )
  size: { type: String, default: 'md' },
  disabled: { type: Boolean, default: false },
  // 允许传入自定义类名进行微调
  customClass: { type: String, default: '' }
});

const emit = defineEmits(['click']);

// 1. 计算形状类名 (完全依据你的“方形、圆形”要求)
const shapeClasses = computed(() => {
  switch (props.shape) {
    case 'rect':   return 'rounded-none';       // 彻底直角
    case 'pill':   return 'rounded-full px-6';  // 胶囊形 (圆形边框)
    case 'circle': return 'rounded-full w-12 h-12 p-0'; // 严格正圆形 (需配合 icon 使用)
    default:       return 'rounded-lg';         // 默认圆角方形
  }
});

// 2. 计算主题色调类名 (完全采用你提供的 5 个色值)
const themeClasses = computed(() => {
  const themes = {
    // #cf4813 (朱砂红)
     cinnabar: 'bg-[#cf4813] text-white hover:bg-[#a63a0f] focus:ring-[#cf4813]',
    // #D2691E (巧克力色)
    chocolate: 'bg-[#D2691E] text-white hover:bg-[#a65217] focus:ring-[#D2691E]',
    // #B22222 (砖红色)
    brick: 'bg-[#B22222] text-white hover:bg-[#8e1b1b] focus:ring-[#B22222]',
    // #b7d07a (嫩绿色) - 使用深色文字保持对比度
    sprout: 'bg-[#b7d07a] text-[#4a5f2e] hover:bg-[#9cb85a] focus:ring-[#b7d07a]',
    // #f2cac9 (浅粉色) - 使用深色文字保持对比度
    lotus: 'bg-[#f2cac9] text-[#7a4e4d] hover:bg-[#e6a8a7] focus:ring-[#f2cac9]'
  };
  return themes[props.theme] || themes.cinnabar;
});

// 3. 计算尺寸类名
const sizeClasses = computed(() => {
  if (props.shape === 'circle') return ''; // 圆形按钮尺寸独立控制
  const sizes = {
    sm: 'px-3 py-1.5 text-xs tracking-wider',
    md: 'px-5 py-2.5 text-sm tracking-widest', // 默认尺寸，增加字间距提升古风感
    lg: 'px-8 py-4 text-base tracking-widest'
  };
  return sizes[props.size] || sizes.md;
});

const handleClick = (e) => {
  if (!props.disabled) {
    emit('click', e);
  }
};
</script>

<style scoped>
/* 此处无需任何 CSS，全部通过 Tailwind 类名实现 */
/* 如果你想强制覆盖全局字体，可以写在这里 */
.font-serif {
  font-family: 'Noto Serif SC', serif;
}
</style>