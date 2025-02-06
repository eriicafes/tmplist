export type State<T> = {
  get(): T;
  set(value: T | ((prev: T) => T)): void;
  listen(
    listener: (value: T) => void | (() => void),
    options?: { eager?: boolean; signal?: AbortSignal }
  ): () => void;
};
export type ReadableState<T> = Omit<State<T>, "set">;

export function $state<T>(initialValue: T): State<T> {
  let _value = initialValue;
  let _callbacks: (() => void)[] = [];
  const _listeners = new Set<(value: T) => void | (() => void)>();

  return {
    get() {
      return _value;
    },
    set(value) {
      let newValue;
      // compute new value
      if (typeof value === "function") {
        newValue = (value as (prev: T) => T)(_value);
      } else {
        newValue = value;
      }
      // if value has not changed return early
      if (_value === newValue) return;

      // execute listener callbacks before updating value
      for (const cb of _callbacks) cb();
      _callbacks = [];
      _value = newValue;

      // execute listeners with new value
      for (const fn of _listeners) {
        const cb = fn(_value);
        if (cb) _callbacks.push(cb);
      }
    },
    listen(listener, options) {
      // add listener and return an unsubscribe function
      _listeners.add(listener);
      if (options?.eager) {
        const cb = listener(_value);
        if (cb) _callbacks.push(cb);
      }
      if (options?.signal) {
        options.signal.addEventListener("abort", () => {
          _listeners.delete(listener);
        });
      }
      return () => _listeners.delete(listener);
    },
  };
}

export function $derived<T, U>(
  state: ReadableState<T>,
  mapper: (value: T) => U
): ReadableState<U> {
  return {
    get() {
      return mapper(state.get());
    },
    listen(listener) {
      return state.listen((value) => listener(mapper(value)));
    },
  };
}

export function $computed<
  T extends [ReadableState<any>, ...ReadableState<any>[]],
  U
>(deps: T, mapper: (...values: InferDeps<T>) => U): ReadableState<U> {
  // compute initial value
  const values = deps.map((dep) => dep.get()) as InferDeps<T>;
  const state = $state(mapper(...values));

  // recompute value with effect
  $effect(deps, (...values) => state.set(mapper(...values)));

  return {
    get: state.get,
    listen: state.listen,
  };
}

type InferDeps<T extends ReadableState<any>[]> = {
  [K in keyof T]: T[K] extends ReadableState<infer U> ? U : never;
};

export function $effect<
  T extends [ReadableState<any>, ...ReadableState<any>[]]
>(
  deps: T,
  callback: (...values: InferDeps<T>) => void | (() => void),
  options?: { signal?: AbortSignal }
) {
  let _cleanup: void | (() => void);
  let _values: InferDeps<T> | null = null;
  let _timeout: number | null = null;

  const listener = () => {
    // debounce repeated calls to listener using setTimeout
    if (_timeout) clearTimeout(_timeout);
    _timeout = setTimeout(() => {
      const values = deps.map((dep) => dep.get()) as InferDeps<T>;
      const unchanged = _values && values.every((v, i) => v === _values?.[i]);
      // if deps have not changed return early
      if (unchanged) return;
      // call cleanup function and rerun effect with new values
      if (_cleanup) _cleanup();
      _values = values;
      _cleanup = callback(...values);
    });
  };
  // call listener for any dep change
  for (const dep of deps) dep.listen(listener, options);
}

export type MountHandle = {
  unmount(): void;
  reset(): void;
};

export function $onMount(
  initializer: (signal: AbortSignal) => void | (() => void)
): MountHandle {
  let _controller = new AbortController();
  let _cb = initializer(_controller.signal);
  return {
    unmount() {
      if (_cb) _cb();
      _controller.abort();
    },
    reset() {
      this.unmount();
      _controller = new AbortController();
      _cb = initializer(_controller.signal);
    },
  };
}

export function $ref<E extends Element = Element>(selector: string) {
  return document.querySelector<E>(selector);
}
