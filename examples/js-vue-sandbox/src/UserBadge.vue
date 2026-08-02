<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';

interface Props {
  displayName: string;
  isActive?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  isActive: false
});

const emit = defineEmits<{
  (e: 'click:badge', value: string): void;
}>();

const renderTime = ref<string>('');
let timerId: number | null = null;

onMounted(() => {
  renderTime.value = new Date().toLocaleTimeString();

  // Lifecycle memory protection - setup timer
  timerId = window.setInterval(() => {
    renderTime.value = new Date().toLocaleTimeString();
  }, 1000);
});

onUnmounted(() => {
  // Clear hooks dynamically on unmount to prevent leaks
  if (timerId !== null) {
    window.clearInterval(timerId);
  }
});

function handleBadgeClick() {
  emit('click:badge', props.displayName);
}
</script>

<template>
  <div
    class="user-badge"
    :class="{ 'user-badge--active': isActive }"
    @click="handleBadgeClick"
  >
    <span class="user-badge__name">{{ displayName }}</span>
    <span class="user-badge__time">Loaded: {{ renderTime }}</span>
  </div>
</template>

<style scoped>
.user-badge {
  display: inline-flex;
  gap: 8px;
  padding: 6px 12px;
  background-color: #f3f4f6;
  border-radius: 4px;
  cursor: pointer;
}

.user-badge--active {
  background-color: #d1fae5;
}
</style>
