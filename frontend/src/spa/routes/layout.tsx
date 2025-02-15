import { GearIcon } from "@radix-ui/react-icons";
import { useMutation, useQuery } from "@tanstack/react-query";
import { LogOutIcon, PlusIcon, XIcon } from "lucide-react";
import { useEffect, useRef, useState, type RefObject } from "react";
import { Link, Outlet, useNavigate } from "react-router";
import { mutations, queries } from "../api";
import { Portal } from "../components/portal";
import { useToggle } from "../hooks/use-toggle";

export default function Layout() {
  const navigate = useNavigate();
  const settingsDialogToggle = useToggle();
  const settingsTriggerRef = useRef<HTMLButtonElement>(null);

  const user = useQuery(queries.profile);
  const logout = useMutation(mutations.logout);

  const handleLogout = () => {
    logout.mutate(undefined, {
      onSuccess() {
        navigate("/login");
      },
    });
  };

  return (
    <>
      <header className="fixed top-0 inset-x-0 z-50 bg-zinc-800 border-b border-zinc-900">
        <div className="flex items-center justify-between px-4 py-2">
          <div className="flex items-center gap-1">
            <Link to="/" className="font-semibold tracking-wide text-lg">
              Tmplist
            </Link>
            <a
              href="/"
              className="border rounded-full px-2 py-px text-[9px] uppercase"
              type="button"
            >
              SPA
            </a>
          </div>

          <div className="flex items-center gap-2">
            <button
              ref={settingsTriggerRef}
              onClick={settingsDialogToggle.toggle}
              className="flex items-center justify-center gap-1 bg-zinc-700 text-zinc-400 w-8 h-8 sm:w-auto sm:px-4 text-xs font-medium rounded-full"
            >
              <GearIcon className="size-5" />
              <span className="hidden sm:inline">Settings</span>
            </button>
            {user.data ? (
              <button
                onClick={handleLogout}
                className="flex items-center justify-center gap-1 bg-sky-200 text-zinc-800 w-8 h-8 sm:w-auto sm:px-4 text-xs font-medium rounded-full"
              >
                <LogOutIcon className="size-5 stroke-[1.5]" />
                <span className="hidden sm:inline">Logout</span>
              </button>
            ) : (
              <Link to="/login">
                <button className="flex items-center justify-center gap-1 bg-sky-200 text-zinc-800 w-8 h-8 sm:w-auto sm:px-4 text-xs font-medium rounded-full">
                  <PlusIcon className="size-5" />
                  <span className="hidden sm:inline">Login</span>
                </button>
              </Link>
            )}
          </div>
        </div>
      </header>

      <main className="mt-12 py-8 max-w-4xl mx-auto px-4">
        <Outlet />
      </main>

      <div id="portal" />

      <SettingsDialog
        open={settingsDialogToggle.open}
        toggle={settingsDialogToggle.toggle}
        anchor={settingsTriggerRef}
      />
    </>
  );
}

function SettingsDialog(props: {
  open: boolean;
  toggle: () => void;
  anchor: RefObject<HTMLButtonElement | null>;
}) {
  const [rect, setRect] = useState(new DOMRect());

  const getCookieValue = (name: string) => {
    const match = document.cookie.match(
      "(^|;)\\s*" + name + "\\s*=\\s*([^;]+)"
    );
    return match?.pop();
  };

  useEffect(() => {
    const anchor = props.anchor.current;
    if (!anchor) return;

    setRect(anchor.getBoundingClientRect());

    const controller = new AbortController();
    window.addEventListener(
      "resize",
      () => {
        setRect(anchor.getBoundingClientRect());
      },
      { signal: controller.signal }
    );
    return () => controller.abort();
  }, [props.anchor]);

  return (
    <Portal>
      <div
        data-open={props.open || undefined}
        style={{
          top: rect.bottom,
          right: window.innerWidth - rect.right,
        }}
        className="hidden data-open:block absolute w-80 max-w-[calc(100vw-72px)] translate-y-2 shadow-xl bg-zinc-700 rounded-3xl p-4"
      >
        <div className="mb-4 flex items-center justify-between">
          <h3 className="font-semibold text-lg">Settings</h3>
          <button onClick={props.toggle} data-role="close">
            <XIcon className="size-5 stroke-1" />
          </button>
        </div>
        <form method="post" action="/" className="grid gap-4">
          <div className="grid gap-1">
            <label className="font-medium text-zinc-400 text-[10px]">
              Mode
            </label>
            <div className="rounded-xl border border-zinc-600 focus-within:border-zinc-400 px-2">
              <select
                name="mode"
                defaultValue={getCookieValue("mode")}
                className="w-full text-xs bg-transparent h-9 focus:outline-none"
              >
                <option value="none">None</option>
                <option value="classic">Classic</option>
                <option value="enhanced">Enhanced</option>
                <option value="spa">SPA</option>
              </select>
            </div>
          </div>
          <button className="flex items-center justify-center gap-1 bg-sky-200 text-zinc-800 h-9 text-xs font-medium rounded-full">
            Save
          </button>
        </form>
      </div>
    </Portal>
  );
}
