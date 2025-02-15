import type { ReactNode } from "react";
import { createPortal } from "react-dom";

export function Portal(props: { children: ReactNode }) {
  const portal = document.getElementById("portal");
  if (!portal) return null;
  return createPortal(props.children, portal);
}
