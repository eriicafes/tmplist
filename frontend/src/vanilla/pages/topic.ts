// refs
const formContainer =
  document.querySelector<HTMLDivElement>("#form-container")!;
const topicForm = document.querySelector<HTMLFormElement>("#topic-form")!;
const deleteTopicTrigger = document.querySelector<HTMLButtonElement>(
  "#delete-topic-trigger"
)!;

// enable submit button when topic input is dirty
topicForm.addEventListener("input", () => {
  const submitButton = formContainer.querySelector(
    'button[form="topic-form"]'
  ) as HTMLButtonElement;
  submitButton.disabled = false;
});

// move cursor to end for autofocused todo inputs
formContainer
  .querySelectorAll<HTMLInputElement>("input[autofocus]")
  .forEach((el) => {
    const length = el.value.length;
    if (length) el.setSelectionRange(length, length);
  });

// submit form when checkbox is clicked
formContainer.addEventListener("change", (e) => {
  const target = e.target as HTMLInputElement | null;
  if (target?.name === "todo-checked") {
    return target?.closest("form")?.submit();
  }
});

// submit forms on enter
formContainer.addEventListener("keypress", (e) => {
  if (e.key === "Enter") {
    e.preventDefault();
    const target = e.target as HTMLInputElement | null;
    target?.closest("form")?.submit();
  }
});

// delete topic with confirmation
deleteTopicTrigger.addEventListener("click", (e) => {
  if (!confirm("Are you sure you want to delete this topic?")) {
    e.preventDefault();
  }
});
