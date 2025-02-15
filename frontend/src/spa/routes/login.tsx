import { useMutation } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router";
import { mutations, ServerError } from "../api";
import { Toast } from "../components/toast";

export default function Login() {
  const navigate = useNavigate();
  const login = useMutation(mutations.login);
  const error = ServerError.check<{ email: string; password: string }>(
    login.error
  );
  const form = useForm<{ email: string; password: string }>();

  const handleSubmit = form.handleSubmit((data) => {
    login.mutate(data, {
      onSuccess() {
        navigate("/");
      },
    });
  });

  return (
    <section className="max-w-xl mx-auto text-sm">
      <div className="mb-4">
        <p className="text-2xl font-medium mb-2">Login your account</p>
        <p className="font-light">
          Enter email address and password to login your account.
        </p>
      </div>

      <form onSubmit={handleSubmit} className="grid gap-4">
        <div className="grid gap-1">
          <label className="text-zinc-500 px-1">Email</label>
          <input
            type="text"
            placeholder="Your email address"
            className="border h-12 w-full rounded-md px-3 bg-transparent border-zinc-700 focus:border-zinc-400 focus:outline-none"
            {...form.register("email", { required: true })}
          />
          {error?.errors?.email && (
            <p className="text-xs text-red-400 px-1">{error.errors.email}</p>
          )}
        </div>

        <div className="grid gap-1">
          <label className="text-zinc-500 px-1">Password</label>
          <input
            type="password"
            placeholder="Your password"
            className="border h-12 w-full rounded-md px-3 bg-transparent border-zinc-700 focus:border-zinc-400 focus:outline-none"
            {...form.register("password", { required: true })}
          />
          {error?.errors?.password && (
            <p className="text-xs text-red-400 px-1">{error.errors.password}</p>
          )}
        </div>

        {error?.message && <Toast message={error.message} type="error" />}

        <div className="grid pt-4 border-t border-zinc-700">
          <button className="bg-sky-200 text-zinc-800 font-semibold h-12 rounded-md">
            Log In
          </button>
        </div>

        <p className="text-center">
          Don't have an account yet?{" "}
          <Link to="/register" className="font-medium underline">
            Register
          </Link>
        </p>
      </form>
    </section>
  );
}
