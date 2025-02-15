import { $onMount, $ref, $state } from "../utils/reactivity";

$onMount(() => {
  // state
  const showSettingsDialog = $state(false);

  // refs
  const settingsTrigger = $ref<HTMLButtonElement>("#settings-trigger")!;
  const settingsDialog = $ref<HTMLDivElement>("#settings-dialog")!;

  // toggle settings dialog
  settingsTrigger.addEventListener("click", () => {
    showSettingsDialog.set((prev) => !prev);
  });

  // close settings dialog from within
  settingsDialog.addEventListener("click", (e) => {
    const target = e.target as HTMLElement;
    if (target.dataset.role === "close") {
      showSettingsDialog.set(false);
    }
  });

  // show/hide settings dialog
  showSettingsDialog.listen((show) => {
    const syncPosition = () => {
      const rect = settingsTrigger.getBoundingClientRect();
      settingsDialog.style.top = `${rect.bottom}px`;
      settingsDialog.style.right = `${window.innerWidth - rect.right}px`;
    };
    syncPosition();
    settingsDialog.toggleAttribute("data-open", show);
    if (!show) return;

    const controller = new AbortController();
    // update anchor position on resize
    window.addEventListener("resize", syncPosition, {
      signal: controller.signal,
    });
    // close on click outside
    document.addEventListener(
      "click",
      (event) => {
        const inDialog = settingsDialog.contains(event.target as Node);
        const inTrigger = settingsTrigger.contains(event.target as Node);
        if (!inDialog && !inTrigger) showSettingsDialog.set(false);
      },
      { signal: controller.signal }
    );
    // close on escape
    document.addEventListener(
      "keyup",
      (e) => {
        if (e.key === "Escape") {
          showSettingsDialog.set(false);
        }
      },
      { signal: controller.signal }
    );
    return () => controller.abort();
  });
});
