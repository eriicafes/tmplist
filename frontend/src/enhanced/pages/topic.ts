import { $onMount } from "../../utils/reactivity";

$onMount(() => {
  document.addEventListener("focusInput", (e) => {
    const el = document.querySelector((e as CustomEvent).detail.value);
    const length = el.value.length;
    el.focus();
    el.setSelectionRange(length, length);
  });
});
