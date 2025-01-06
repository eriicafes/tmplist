<script lang="ts">
  let count = $state(0);
</script>

<div
  class="max-w-sm mx-auto mt-40 flex items-center justify-between p-4 rounded-xl bg-red-50 text-red-500 font-thin border border-red-100"
>
  <p>Svelte</p>
  <span>{count}</span>
  <button class="text-xl" onclick={() => count++}>+</button>
</div>
