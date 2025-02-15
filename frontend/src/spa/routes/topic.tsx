import { useMutation, useSuspenseQuery } from "@tanstack/react-query";
import { PlusIcon, TrashIcon } from "lucide-react";
import { useRef } from "react";
import { Link, useNavigate, useParams } from "react-router";
import { mutations, queries } from "../api";
import { TodoItem } from "../components/todo-item";
import type { Todo } from "../models";

export default function Topic() {
  const user = useSuspenseQuery(queries.profile);
  const { id } = useParams();
  const topic = useSuspenseQuery(queries.getTopic({ topicId: +id! }));
  const navigate = useNavigate();
  const updateTopic = useMutation(mutations.updateTopic);
  const deleteTopic = useMutation(mutations.deleteTopic);
  const createTodo = useMutation(mutations.createTodo);
  const updateTodo = useMutation(mutations.updateTodo);
  const deleteTodo = useMutation(mutations.deleteTodo);
  const topicInputRef = useRef<HTMLInputElement>(null);
  const newTodoInputRef = useRef<HTMLInputElement>(null);

  const handleDelete = () => {
    deleteTopic.mutate(
      {
        topicId: topic.data.topic.id,
      },
      {
        onSuccess() {
          navigate("/");
        },
      }
    );
  };

  const handleChangeTodo = (todo: Todo, text: string, checked: boolean) => {
    updateTodo.mutate({
      topicId: todo.topic_id,
      todoId: todo.id,
      data: { text, checked },
    });
  };

  const handleDeleteTodo = (todo: Todo) => {
    deleteTodo.mutate({
      topicId: todo.topic_id,
      todoId: todo.id,
    });
  };

  return (
    <section>
      <div className="mb-4 h-7 flex items-center justify-between flex-wrap leading-none">
        <Link to="/" className="text-sm text-zinc-300 hover:text-zinc-100">
          &larr; Back to topics
        </Link>
        <p className="text-sm text-zinc-300">{user.data.email}</p>
      </div>

      <div className="max-w-3xl mx-auto py-4 bg-zinc-800 borders border-zinc-700 rounded-xl">
        <div id="form-container" className="grid">
          <div className="flex items-start gap-1">
            <input
              ref={topicInputRef}
              required
              type="text"
              placeholder="Enter topic..."
              defaultValue={topic.data.topic.title}
              className="flex-1 px-2 h-12 focus:outline-none text-lg font-medium"
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  e.preventDefault();
                  const text = e.currentTarget.value.trim();
                  if (!text) return;
                  updateTopic.mutate({
                    topicId: topic.data.topic.id,
                    data: { topic: text },
                  });
                }
              }}
            />
            <button
              onClick={handleDelete}
              type="button"
              className="p-1 mt-2 rounded-full text-zinc-400 hover:text-zinc-100 transition-colors"
            >
              <TrashIcon className="size-5 stroke-1" />
            </button>
          </div>
          <div>
            <form
              method="post"
              action="/classic/{{.Topic.Id}}/todos"
              className="flex items-center gap-1 text-sm text-zinc-400 border-b border-transparent focus-within:border-zinc-700 transition-colors"
            >
              <label htmlFor="new-todo-input">
                <PlusIcon className="size-5 stroke-2" />
              </label>
              <input
                ref={newTodoInputRef}
                id="new-todo-input"
                required
                type="text"
                autoFocus
                placeholder="New todo item"
                className="flex-1 px-2 h-10 focus:outline-none"
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    e.preventDefault();
                    const text = e.currentTarget.value.trim();
                    if (!text) return;
                    const target = e.currentTarget;
                    createTodo.mutate(
                      {
                        topicId: topic.data.topic.id,
                        data: { text, checked: false },
                      },
                      {
                        onSuccess() {
                          target.value = "";
                        },
                      }
                    );
                  }
                }}
              />
            </form>
            <div>
              {topic.data.todos
                .filter((todo) => !todo.done)
                .map((todo) => (
                  <TodoItem
                    key={todo.id}
                    text={todo.body}
                    checked={todo.done}
                    onCheck={(on) => handleChangeTodo(todo, todo.body, on)}
                    onEnter={(val) => handleChangeTodo(todo, val, todo.done)}
                    onDelete={() => handleDeleteTodo(todo)}
                  />
                ))}
            </div>
            <div
              data-label="Completed todos"
              className="mt-4 pt-4 border-t border-zinc-700 not-empty:before:content-[attr(data-label)] before:text-zinc-500 before:text-xs before:block before:pb-2"
            >
              {topic.data.todos
                .filter((todo) => todo.done)
                .map((todo) => (
                  <TodoItem
                    key={todo.id}
                    text={todo.body}
                    checked={todo.done}
                    onCheck={(on) => handleChangeTodo(todo, todo.body, on)}
                    onEnter={(val) => handleChangeTodo(todo, val, todo.done)}
                    onDelete={() => handleDeleteTodo(todo)}
                  />
                ))}
            </div>
          </div>
          <div className="mt-2">
            <p className="text-xs text-zinc-400">
              {topic.data.todos.filter((todos) => todos.done).length} /{" "}
              {topic.data.todos.length} completed
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
