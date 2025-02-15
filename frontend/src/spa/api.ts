import {
  QueryClient,
  queryOptions,
  type MutationOptions,
} from "@tanstack/react-query";
import type {
  LoginData,
  RegisterData,
  Todo,
  TodoData,
  Topic,
  TopicData,
  User,
} from "./models";

const mutationOptions = <
  TData = unknown,
  TError = Error,
  TVariables = void,
  TContext = unknown
>(
  opts: MutationOptions<TData, TError, TVariables, TContext>
) => opts;

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      refetchOnWindowFocus: false,
      staleTime: 1000, // 1 second
    },
  },
});

export const queries = {
  profile: queryOptions({
    queryKey: ["profile"],
    async queryFn({ signal }) {
      return await request<User>("/api/profile", { signal });
    },
  }),

  getAllTopics: (search: string) =>
    queryOptions({
      queryKey: ["topics", { search }],
      async queryFn({ signal }) {
        const query = new URLSearchParams();
        if (search) query.append("search", search);
        return request<Topic[]>(`/api?${query}`, { signal });
      },
    }),

  getTopic: (data: { topicId: number }) =>
    queryOptions({
      queryKey: ["topics", data.topicId],
      async queryFn({ signal }) {
        return request<{ topic: Topic; todos: Todo[] }>(
          `/api/${data.topicId}`,
          { signal }
        );
      },
    }),
};

export const mutations = {
  createTopic: mutationOptions({
    async mutationFn(data: { data: TopicData }) {
      return request<Topic>(`/api`, {
        method: "POST",
        body: data.data,
      });
    },
    onSuccess() {
      queryClient.invalidateQueries({ ...queries.getAllTopics, exact: true });
    },
  }),

  updateTopic: mutationOptions({
    async mutationFn(data: {
      topicId: number;
      data: Pick<TopicData, "topic">;
    }) {
      return request<Topic>(`/api/${data.topicId}`, {
        method: "PUT",
        body: data.data,
      });
    },
    onSuccess(data) {
      return queryClient.invalidateQueries(
        queries.getTopic({ topicId: data.id })
      );
    },
  }),

  deleteTopic: mutationOptions({
    async mutationFn(data: { topicId: number }) {
      return request<Topic>(`/api/${data.topicId}`, {
        method: "DELETE",
      });
    },
    onSuccess() {
      queryClient.invalidateQueries({ ...queries.getAllTopics, exact: true });
    },
  }),

  createTodo: mutationOptions({
    async mutationFn(data: { topicId: number; data: TodoData }) {
      return request<Todo[]>(`/api/${data.topicId}/todos`, {
        method: "POST",
        body: data.data,
      });
    },
    onSuccess(data) {
      if (!data[0]) return;
      return queryClient.invalidateQueries(
        queries.getTopic({ topicId: data[0].topic_id })
      );
    },
  }),

  updateTodo: mutationOptions({
    async mutationFn(data: {
      topicId: number;
      todoId: number;
      data: TodoData;
    }) {
      return request<Todo>(`/api/${data.topicId}/todos/${data.todoId}`, {
        method: "PUT",
        body: data.data,
      });
    },
    onSuccess(data) {
      return queryClient.invalidateQueries(
        queries.getTopic({ topicId: data.topic_id })
      );
    },
  }),

  deleteTodo: mutationOptions({
    async mutationFn(data: { topicId: number; todoId: number }) {
      return request<Todo>(`/api/${data.topicId}/todos/${data.todoId}`, {
        method: "DELETE",
      });
    },
    onSuccess(data) {
      return queryClient.invalidateQueries(
        queries.getTopic({ topicId: data.topic_id })
      );
    },
  }),

  login: mutationOptions({
    async mutationFn(data: LoginData) {
      return request<User>("/api/login", {
        method: "POST",
        body: data,
      });
    },
    onSuccess() {
      return queryClient.invalidateQueries(queries.profile);
    },
  }),

  register: mutationOptions({
    async mutationFn(data: RegisterData) {
      return request<{ message: string; profile: User }>("/api/register", {
        method: "POST",
        body: data,
      });
    },
    onSuccess() {
      return queryClient.invalidateQueries(queries.profile);
    },
  }),

  logout: mutationOptions({
    async mutationFn() {
      return request<void>("/api/logout", {
        method: "POST",
      });
    },
    onSuccess() {
      return queryClient.removeQueries(queries.profile);
    },
  }),
};

interface JSONRequestInit extends Omit<RequestInit, "body"> {
  body?: Record<string, any>;
}

async function request<T>(
  input: RequestInfo | URL,
  init: JSONRequestInit = {}
) {
  const { body, ...rest } = init;
  let requestInit: RequestInit = rest;
  requestInit.headers = new Headers(requestInit.headers);
  requestInit.headers.set("Accept", "application/json");
  if (body) {
    requestInit.body = JSON.stringify(body);
    requestInit.headers.set("Content-Type", "application/json");
  }
  const res = await fetch(input, requestInit);
  const json = await res.json();
  if (!res.ok) throw ServerError.fromResponse(json);
  return json as Promise<T>;
}

export class ServerError<
  T extends Record<string, string> = Record<string, string>
> extends Error {
  constructor(message: string, public errors: T | null) {
    super(message);
  }

  public static check<
    T extends Record<string, string> = Record<string, string>
  >(error: unknown): ServerError<T> | null {
    return error instanceof ServerError ? error : null;
  }

  public static fromResponse(json: any) {
    if ("message" in json && "errors" in json) {
      return new ServerError(json.message, json.errors);
    }
    return new ServerError("Invalid response", null);
  }
}
