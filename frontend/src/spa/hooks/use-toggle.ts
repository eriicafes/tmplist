import { useState } from "react";

export function useToggle() {
  const [open, setOpen] = useState(false);

  const toggle = () => setOpen(!open);

  return { open, setOpen, toggle };
}
