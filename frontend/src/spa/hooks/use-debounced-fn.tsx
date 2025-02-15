import { useRef } from "react";

export function useDeboundedFn<T extends any[]>(
  delay: number,
  fn: (...args: T) => void
) {
  const timeoutRef = useRef<number>(undefined);

  return (...args: T) => {
    if (timeoutRef.current) clearTimeout(timeoutRef.current);
    timeoutRef.current = setTimeout(() => fn(...args), delay);
  };
}
