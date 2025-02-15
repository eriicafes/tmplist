import { $onMount, $ref, $state } from "../../utils/reactivity";
import { htmx, swappedTarget } from "../htmx";

htmx.on("htmx:afterSwap", (e) => {
  if (swappedTarget(e, "#new-topic-trigger")) {
    newTopicTrigger.reset();
  }
  if (swappedTarget(e, "#topics")) {
    newTopicDialog.reset();
    topicForm.reset();
  }
});

htmx.on("htmx:afterRequest", () => {
  showNewTopicDialog.set(false);
});

// state
const showNewTopicDialog = $state(false);

const newTopicTrigger = $onMount(() => {
  // refs
  const newTopicTrigger = $ref<HTMLButtonElement>("#new-topic-trigger")!;

  // toggle new topic dialog
  newTopicTrigger.addEventListener("click", () =>
    showNewTopicDialog.set((prev) => !prev)
  );
});

const newTopicDialog = $onMount((signal) => {
  // refs
  const newTopicDialog = $ref<HTMLDivElement>("#new-topic-dialog")!;
  const newTopicOverlay = $ref<HTMLDivElement>("#new-topic-overlay")!;

  // close new topic dialog from within
  newTopicDialog.addEventListener("click", (e) => {
    const target = e.target as HTMLElement | null;
    if (target?.dataset.role === "close") {
      showNewTopicDialog.set(false);
    }
  });

  // close new topic dialog on overlay click
  newTopicOverlay.addEventListener("click", () => {
    showNewTopicDialog.set(false);
  });

  // show dialog on cmd+k
  document.addEventListener(
    "keydown",
    (e) => {
      if (e.metaKey && e.key === "k") {
        showNewTopicDialog.set(true);
      }
    },
    { signal }
  );

  // show/hide new topic dialog
  showNewTopicDialog.listen(
    (show) => {
      newTopicDialog.toggleAttribute("data-open", show);
      newTopicOverlay.toggleAttribute("data-open", show);
      if (!show) return;

      newTopicDialog.querySelector("input")?.focus();

      // listen to events only when dialog is open
      const controller = new AbortController();
      // close on escape
      document.addEventListener(
        "keyup",
        (e) => {
          if (e.key === "Escape") {
            showNewTopicDialog.set(false);
          }
        },
        { signal: controller.signal }
      );
      return () => controller.abort();
    },
    { signal }
  );
});

const topicForm = $onMount(() => {
  // refs
  const topicForm = $ref<HTMLFormElement>("#topic-form")!;
  const newTodoInput = $ref<HTMLInputElement>("#new-todo-input")!;
  const newTodoInputTemplate = $ref<HTMLTemplateElement>(
    "#new-todo-input-template"
  )!;
  const newTodos = $ref<HTMLDivElement>("#new-todos")!;
  const completedTodos = $ref<HTMLDivElement>("#completed-todos")!;

  // add new todo
  newTodoInput.addEventListener("keypress", (e) => {
    if (e.key === "Enter") {
      e.preventDefault();
      const value = newTodoInput.value.trim();
      if (!value) return;
      newTodoInput.value = "";
      const clone = newTodoInputTemplate.content.firstElementChild?.cloneNode(
        true
      ) as Element;
      const input = clone.querySelector(
        "input[name='todo']"
      ) as HTMLInputElement;
      input.value = value;
      newTodos.insertAdjacentElement("afterbegin", clone);
    }
  });

  // add draft todo input before submit
  topicForm
    .querySelector("button[type=submit]")
    ?.addEventListener("click", () => {
      const value = newTodoInput.value.trim();
      if (!value) return;
      newTodoInput.value = "";
      const clone = newTodoInputTemplate.content.firstElementChild?.cloneNode(
        true
      ) as Element;
      const input = clone.querySelector(
        "input[name='todo']"
      ) as HTMLInputElement;
      input.value = value;
      newTodos.insertAdjacentElement("afterbegin", clone);
    });

  // toggle todo completion
  topicForm.addEventListener("change", (e) => {
    const target = e.target as HTMLInputElement | null;
    if (target?.name === "todo-checked") {
      const todoItem = target.closest("div.group")!;
      if (target.checked) {
        completedTodos.insertAdjacentElement("afterbegin", todoItem);
      } else {
        newTodos.insertAdjacentElement("beforeend", todoItem);
      }
    }
  });
});
