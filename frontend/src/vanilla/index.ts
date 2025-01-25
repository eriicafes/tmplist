import { $state } from "./utils/reactivity";

// refs
const settingsTrigger =
  document.querySelector<HTMLButtonElement>("#settings-trigger")!;
const settingsDialog =
  document.querySelector<HTMLDivElement>("#settings-dialog")!;

// state
const showSettingsDialog = $state(false);
const triggerRect = $state(settingsTrigger.getBoundingClientRect());

// toggle settings dialog
settingsTrigger.addEventListener("click", () => {
  showSettingsDialog.set((prev) => !prev);
});

// close settings dialog from within
settingsDialog.addEventListener("click", (e) => {
  const target = e.target as HTMLElement | null;
  if (target?.dataset.role === "close") {
    showSettingsDialog.set(false);
  }
});

// sync settings dialog anchor position
triggerRect.listen((rect) => {
  settingsDialog.style.top = `${rect.bottom}px`;
  settingsDialog.style.right = `${window.innerWidth - rect.right}px`;
});

// show/hide settings dialog
showSettingsDialog.listen((show) => {
  settingsDialog.toggleAttribute("data-open", show);
  triggerRect.set(settingsTrigger.getBoundingClientRect());
  if (!show) return;

  const controller = new AbortController();
  // update anchor position on resize
  window.addEventListener(
    "resize",
    () => {
      triggerRect.set(settingsTrigger.getBoundingClientRect());
    },
    { signal: controller.signal }
  );
  // close on click outside
  document.addEventListener(
    "click",
    (event) => {
      const inDialog = settingsDialog.contains(event.target as Node);
      const inTrigger = settingsTrigger.contains(event.target as Node);
      if (!inDialog && !inTrigger) showSettingsDialog.set(false);
    },
    { signal: controller.signal, capture: true }
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
