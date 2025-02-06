import { $state } from "../../utils/reactivity";

// refs
const topicForm = document.querySelector<HTMLFormElement>("#topic-form")!;
const newTopicTrigger =
  document.querySelector<HTMLButtonElement>("#new-topic-trigger")!;
const newTopicDialog =
  document.querySelector<HTMLDivElement>("#new-topic-dialog")!;
const newTopicOverlay =
  document.querySelector<HTMLDivElement>("#new-topic-overlay")!;
const newTodoInput =
  document.querySelector<HTMLInputElement>("#new-todo-input")!;
const newTodoInputTemplate = document.querySelector<HTMLTemplateElement>(
  "#new-todo-input-template"
)!;
const newTodos = document.querySelector<HTMLDivElement>("#new-todos")!;
const completedTodos =
  document.querySelector<HTMLDivElement>("#completed-todos")!;

// state
const showNewTopicDialog = $state(false);

// toggle new topic dialog
newTopicTrigger.addEventListener("click", () => {
  showNewTopicDialog.set((prev) => !prev);
});

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
document.addEventListener("keydown", (e) => {
  if (e.metaKey && e.key === "k") {
    showNewTopicDialog.set(true);
  }
});

// show/hide new topic dialog
showNewTopicDialog.listen((show) => {
  newTopicDialog.toggleAttribute("data-open", show);
  newTopicOverlay.toggleAttribute("data-open", show);
  if (!show) return;

  newTopicDialog.querySelector("input")?.focus();

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
});

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
    const input = clone.querySelector("input[name='todo']") as HTMLInputElement;
    input.value = value;
    newTodos.insertAdjacentElement("afterbegin", clone);
  }
});

// add draft todo input before submit
topicForm
  .querySelector("button[type=submit]")
  ?.addEventListener("click", (e) => {
    if (newTodoInput.value) {
      e.preventDefault();
      newTodoInput.dispatchEvent(
        new KeyboardEvent("keypress", { key: "Enter" })
      );
      topicForm.submit();
    }
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
