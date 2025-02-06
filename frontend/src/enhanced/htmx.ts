import htmx from "htmx.org";

htmx.config.responseHandling = [
  { code: "204", swap: false },
  { code: "[23]..", swap: true },
  { code: "[45]..", swap: true, error: true },
  { code: "...", swap: true },
];

export { htmx };

export function swappedTarget(e: Event, selector: string) {
  return (e.target as Element).querySelector(selector) !== null;
}
