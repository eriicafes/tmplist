import { XIcon } from "lucide-react";

export function TodoItem(props: {
  text: string;
  checked: boolean;
  onCheck(checked: boolean): void;
  onChange?(text: string): void;
  onEnter?(text: string): void;
  onDelete(): void;
}) {
  return (
    <div className="group flex items-center gap-1 text-sm text-zinc-400 border-y border-transparent focus-within:border-zinc-700 hover:border-zinc-700 transition-colors">
      <div className="size-5 flex items-center justify-center">
        <input
          type="checkbox"
          checked={props.checked}
          onChange={(e) => props.onCheck(e.target.checked)}
          className="size-4 appearance-none checked:appearance-auto border-2 border-zinc-600 accent-sky-200 rounded-xs focus:outline-none"
        />
      </div>
      <input
        required
        type="text"
        value={props.onChange ? props.text : undefined}
        onChange={
          props.onChange ? (e) => props.onChange!(e.target.value) : undefined
        }
        defaultValue={props.onEnter ? props.text : undefined}
        onKeyUp={
          props.onEnter
            ? (e) => {
                if (e.key === "Enter") {
                  e.preventDefault();
                  props.onEnter!(e.currentTarget.value);
                }
              }
            : undefined
        }
        className="flex-1 px-2 h-10 focus:outline-none"
      />
      <button
        onClick={props.onDelete}
        type="button"
        className="rounded-full p-1 hidden group-focus-within:block hover:bg-zinc-700 hover:text-zinc-300 transition-colors"
      >
        <XIcon className="size-4 stroke-1" />
      </button>
    </div>
  );
}
