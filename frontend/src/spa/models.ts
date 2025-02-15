// models
export interface User {
  id: number;
  email: string;
  created_at: string;
}

export interface Session {
  id: string;
  user_id: number;
  expires_at: string;
}

export interface Topic {
  id: number;
  user_id: number;
  title: string;
  todos_count: number;
  created_at: string;
}

export interface Todo {
  id: number;
  topic_id: number;
  body: string;
  done: boolean;
  created_at: string;
}

// schemas
export interface LoginData {
  email: string;
  password: string;
}

export interface RegisterData {
  email: string;
  password: string;
}

export interface TopicData {
  topic: string;
  todos: TodoData[];
}

export interface TodoData {
  text: string;
  checked: boolean;
}
