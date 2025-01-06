import { useState } from "react";

export default function App() {
  const [count, setCount] = useState(0);

  return (
    <div className="max-w-sm mx-auto mt-40 flex items-center justify-between p-4 rounded-xl bg-blue-50 text-blue-500 font-thin border border-blue-100">
      <p>React</p>
      <span>{count}</span>
      <button className="text-xl" onClick={() => setCount(count + 1)}>
        +
      </button>
    </div>
  );
}
