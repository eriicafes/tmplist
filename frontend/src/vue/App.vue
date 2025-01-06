<script setup lang="ts">
import { ref } from 'vue';

const count = ref(0)
</script>

<template>
    <div class="max-w-sm mx-auto mt-40 flex items-center justify-between p-4 rounded-xl bg-green-50 text-green-500 font-thin border border-green-100">
      <p>Vue</p>
      <span>{{ count }}</span>
      <button class="text-xl" @click="count++">+</button>
    </div>
</template>

