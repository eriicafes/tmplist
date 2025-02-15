import { InfoCircledIcon } from "@radix-ui/react-icons";
import clsx from "clsx";
import { XIcon } from "lucide-react";
import { useState } from "react";

export function Toast(props: {
  message: string;
  type?: "warning" | "error" | "success";
  className?: string;
}) {
  const [hide, setHide] = useState(false);
  if (hide) return null;

  return (
    <div
      id="toast"
      className={clsx([
        "flex gap-1 border p-3 rounded-lg mb-4",
        props.type === "warning" &&
          "bg-yellow-100 text-yellow-600 border-yellow-400",
        props.type === "error" && "bg-red-100 text-red-500 border-red-300",
        props.type === "success" && "bg-teal-100 text-teal-600 border-teal-400",
        !props.type && "bg-zinc-100 text-zinc-600 border-zinc-400",
        props.className,
      ])}
    >
      <InfoCircledIcon className="size-4 shrink-0" />
      <p className="flex-1 text-xs">{props.message}</p>
      <button className="shrink-0" type="button" onClick={() => setHide(true)}>
        <XIcon className="size-4" />
      </button>
    </div>
  );
}
