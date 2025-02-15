import { useMutation, useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { format } from "date-fns";
import { EllipsisVertical, PlusIcon, SearchIcon, XIcon } from "lucide-react";
import { useEffect, useRef, useState, type FormEvent } from "react";
import { Form, Link, useNavigate, useSearchParams } from "react-router";
import { mutations, queries } from "../api";
import { Portal } from "../components/portal";
import { TodoItem } from "../components/todo-item";
import { useDeboundedFn } from "../hooks/use-debounced-fn";
import { useToggle } from "../hooks/use-toggle";
import type { TodoData } from "../models";

export default function Index() {
  const user = useSuspenseQuery(queries.profile);

  const [searchParams, setSearchParams] = useSearchParams();
  const search = searchParams.get("search") ?? "";
  const setSearch = useDeboundedFn(500, (search: string) => {
    setSearchParams((prev) => {
      prev.set("search", search);
      return prev;
    });
  });
  const [searching, setSearching] = useState(!!search);
  const topics = useQuery(queries.getAllTopics(search));

  const newTopicDialogToggle = useToggle();

  return (
    <section>
      <div className="mb-4 h-7 flex items-center justify-between leading-none">
        <h2 className="text-xl font-medium">Your Topics</h2>
        <p className="text-zinc-300 text-sm">{user.data.email}</p>
      </div>

      <Form
        data-show={searching || undefined}
        className="group flex items-center justify-end gap-1"
      >
        <input
          type="text"
          placeholder="Search topics"
          className="invisible group-data-show:visible sm:max-w-80 flex-1 px-4 h-8 text-xs rounded-full border border-zinc-700 focus:border-zinc-600 focus:outline-none"
          name="search"
          defaultValue={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <button
          type="button"
          onClick={() => setSearching(!searching)}
          className="rounded-full text-zinc-400 hover:text-zinc-200 transition-colors"
        >
          <span className="hidden group-data-show:inline">
            <XIcon className="size-5 stroke-[1.5]" />
          </span>
          <span className="group-data-show:hidden">
            <SearchIcon className="size-5 stroke-[1.5]" />
          </span>
        </button>
      </Form>

      <div className="py-4 grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2 md:gap-6">
        <button
          onClick={newTopicDialogToggle.toggle}
          type="button"
          className="flex flex-col items-center justify-center gap-2 p-2.5 rounded-2xl aspect-square border border-zinc-700 hover:border-zinc-600 hover:bg-zinc-700/20 text-zinc-400 hover:text-zinc-100 transition-colors"
        >
          <PlusIcon className="size-7" />
          <span className="select-none text-sm font-medium">New Topic</span>
        </button>

        {topics.data?.map((topic) => (
          <Link
            key={topic.id}
            to={`/${topic.id}`}
            className="flex flex-col justify-between gap-2 p-2.5 rounded-2xl aspect-square border border-zinc-700 bg-zinc-700 hover:scale-[1.02] transition-transform"
          >
            <div className="flex items-center justify-between">
              <div
                className={`size-5 rounded-full bg-linear-60 ${gradient(
                  topic.id
                )}`}
              />
              <EllipsisVertical className="size-4 text-zinc-400" />
            </div>
            <div className="grid gap-3">
              <p className="text-lg line-clamp-2">{topic.title}</p>
              <p className="text-[10px] text-zinc-300">
                {formatCreatedAt(topic.created_at)} • {topic.todos_count} item
                {topic.todos_count !== 1 && "s"}
              </p>
            </div>
          </Link>
        ))}

        {emptyCells(topics.data?.length ?? 0).map((_, i) => (
          <div
            key={i}
            className="flex flex-col justify-between gap-2 p-2.5 rounded-2xl aspect-square border border-zinc-700/30 bg-zinc-700/30"
          />
        ))}
      </div>

      <NewTopicDialog
        open={newTopicDialogToggle.open}
        toggle={newTopicDialogToggle.toggle}
        setOpen={newTopicDialogToggle.setOpen}
      />
    </section>
  );
}

function NewTopicDialog(props: {
  open: boolean;
  toggle: () => void;
  setOpen: (open: boolean) => void;
}) {
  const navigate = useNavigate();
  const createTopic = useMutation(mutations.createTopic);
  const topicInputRef = useRef<HTMLInputElement>(null);
  const newTodoInputRef = useRef<HTMLInputElement>(null);
  const [todos, setTodos] = useState<TodoData[]>([]);

  useEffect(() => {
    const controller = new AbortController();
    document.addEventListener(
      "keydown",
      (e) => {
        if (e.metaKey && e.key === "k") {
          props.setOpen(true);
        }
      },
      { signal: controller.signal }
    );
    if (props.open) {
      topicInputRef.current?.focus();
      document.addEventListener(
        "keyup",
        (e) => {
          if (e.key === "Escape") {
            props.setOpen(false);
          }
        },
        { signal: controller.signal }
      );
    }
    return () => controller.abort();
  }, [props.open]);

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const topic = topicInputRef.current?.value.trim();
    if (!topic) return;
    const todosData = [...todos];
    const draft = newTodoInputRef.current?.value.trim();
    if (draft) {
      todosData.unshift({ text: draft, checked: false });
      setTodos(todosData);
      newTodoInputRef.current!.value = "";
    }
    createTopic.mutate(
      {
        data: { topic, todos: todosData },
      },
      {
        onSuccess(res) {
          navigate(`/${res.id}`);
        },
      }
    );
  };

  const handleChange = (todo: TodoData, text: string, checked: boolean) => {
    setTodos((prev) => prev.map((t) => (t === todo ? { text, checked } : t)));
  };

  const handleDelete = (todo: TodoData) => {
    setTodos((prev) => prev.filter((t) => t !== todo));
  };

  return (
    <Portal>
      <div
        data-open={props.open || undefined}
        onClick={props.toggle}
        className="hidden data-open:block fixed inset-0 bg-transparent data-open:bg-black/30 transition-colors duration-300"
      />
      <div
        data-open={props.open || undefined}
        className="hidden data-open:block fixed inset-x-0 top-[20vh] w-full max-w-[90vw] sm:max-w-xl mx-auto shadow-xl bg-zinc-800 border border-zinc-700 rounded-xl p-4"
      >
        <form onSubmit={handleSubmit} className="grid">
          <div className="flex items-start gap-1">
            <input
              ref={topicInputRef}
              required
              type="text"
              placeholder="Enter topic..."
              className="flex-1 px-2 h-12 focus:outline-none text-lg font-medium"
            />
            <button
              type="button"
              onClick={props.toggle}
              className="p-1 mt-2 rounded-full text-zinc-400 hover:text-zinc-100 transition-colors"
            >
              <XIcon className="size-5 stroke-1" />
            </button>
          </div>
          <fieldset className="max-h-100 overflow-y-scroll">
            <div className="flex items-center gap-1 text-sm text-zinc-400 border-b border-transparent focus-within:border-zinc-700 transition-colors">
              <label htmlFor="new-todo-input">
                <PlusIcon className="size-5 stroke-2" />
              </label>
              <input
                ref={newTodoInputRef}
                id="new-todo-input"
                type="text"
                placeholder="New todo item"
                className="flex-1 px-2 h-10 focus:outline-none"
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    e.preventDefault();
                    const text = e.currentTarget.value.trim();
                    if (!text) return;
                    setTodos((prev) => [{ text, checked: false }, ...prev]);
                    e.currentTarget.value = "";
                  }
                }}
              />
            </div>
            <div>
              {todos
                .filter((todo) => !todo.checked)
                .map((todo, i) => (
                  <TodoItem
                    key={i}
                    text={todo.text}
                    checked={todo.checked}
                    onCheck={(on) => handleChange(todo, todo.text, on)}
                    onChange={(val) => handleChange(todo, val, todo.checked)}
                    onDelete={() => handleDelete(todo)}
                  />
                ))}
            </div>
            <div
              data-label="Completed todos"
              className="mt-4 pt-4 border-t border-zinc-700 not-empty:before:content-[attr(data-label)] before:text-zinc-500 before:text-xs before:block before:pb-2"
            >
              {todos
                .filter((todo) => todo.checked)
                .map((todo, i) => (
                  <TodoItem
                    key={i}
                    text={todo.text}
                    checked={todo.checked}
                    onCheck={(on) => handleChange(todo, todo.text, on)}
                    onChange={(val) => handleChange(todo, val, todo.checked)}
                    onDelete={() => handleDelete(todo)}
                  />
                ))}
            </div>
          </fieldset>
          <div className="flex justify-end mt-2">
            <button
              type="submit"
              className="flex items-center justify-center gap-1 bg-sky-200 text-zinc-800 h-9 px-4 text-xs font-medium rounded-xl"
            >
              Add Topic
            </button>
          </div>
        </form>
      </div>
    </Portal>
  );
}

function formatCreatedAt(date: string) {
  return format(new Date(date), "MMM d, yyyy");
}

function emptyCells(topicsCount: number) {
  const maxEmptyCells = 7;
  if (topicsCount >= maxEmptyCells) {
    return [];
  }
  return Array.from({ length: maxEmptyCells - topicsCount });
}

function gradient(topicId: number) {
  const gradients = [
    "from-purple-500 to-pink-500",
    "from-green-400 to-blue-500",
    "from-yellow-400 to-red-500",
    "from-blue-400 to-indigo-500",
    "from-red-400 to-yellow-500",
    "from-pink-400 to-purple-500",
    "from-indigo-400 to-green-500",
    "from-teal-400 to-cyan-500",
    "from-orange-400 to-red-500",
    "from-lime-400 to-green-500",
    "from-amber-400 to-yellow-500",
    "from-emerald-400 to-teal-500",
    "from-sky-400 to-blue-500",
    "from-rose-400 to-pink-500",
    "from-fuchsia-400 to-purple-500",
    "from-violet-400 to-indigo-500",
    "from-cyan-400 to-teal-500",
    "from-lime-500 to-green-600",
    "from-amber-500 to-orange-600",
    "from-emerald-500 to-teal-600",
    "from-sky-500 to-blue-600",
    "from-rose-500 to-pink-600",
    "from-fuchsia-500 to-purple-600",
    "from-violet-500 to-indigo-600",
    "from-cyan-500 to-teal-600",
    "from-lime-600 to-green-700",
    "from-amber-600 to-orange-700",
    "from-emerald-600 to-teal-700",
    "from-sky-600 to-blue-700",
    "from-rose-600 to-pink-700",
    "from-fuchsia-600 to-purple-700",
    "from-violet-600 to-indigo-700",
    "from-cyan-600 to-teal-700",
    "from-lime-700 to-green-800",
    "from-amber-700 to-orange-800",
    "from-emerald-700 to-teal-800",
    "from-sky-700 to-blue-800",
    "from-rose-700 to-pink-800",
    "from-fuchsia-700 to-purple-800",
    "from-violet-700 to-indigo-800",
    "from-cyan-700 to-teal-800",
    "from-lime-800 to-green-900",
    "from-amber-800 to-orange-900",
    "from-emerald-800 to-teal-900",
    "from-sky-800 to-blue-900",
    "from-rose-800 to-pink-900",
    "from-fuchsia-800 to-purple-900",
    "from-violet-800 to-indigo-900",
    "from-cyan-800 to-teal-900",
  ];
  return gradients[topicId % gradients.length];
}
