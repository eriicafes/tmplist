import { LoaderCircleIcon } from "lucide-react";
import { lazy } from "react";
import {
  createBrowserRouter,
  redirect,
  type LoaderFunction,
  type RouteObject,
} from "react-router";
import { queries, queryClient } from "./api";

const route = (routeObject: RouteObject) => routeObject;

const Layout = lazy(() => import("./routes/layout"));
const Index = lazy(() => import("./routes/index"));
const Topic = lazy(() => import("./routes/topic"));
const Login = lazy(() => import("./routes/login"));
const Register = lazy(() => import("./routes/register"));

const auth = authGuard("auth");
const guest = authGuard("guest");

export const router = createBrowserRouter(
  [
    route({
      element: <Layout />,
      hydrateFallbackElement: (
        <div className="h-screen w-screen grid place-items-center">
          <LoaderCircleIcon className="size-16 stroke-[0.5] animate-spin text-sky-200" />
        </div>
      ),
      children: [
        route({ index: true, loader: auth, element: <Index /> }),
        route({ path: ":id", loader: auth, element: <Topic /> }),
        route({ path: "login", loader: guest, element: <Login /> }),
        route({ path: "register", loader: guest, element: <Register /> }),
      ],
    }),
  ],
  { basename: "/spa" }
);

function authGuard(mode: "auth" | "guest"): LoaderFunction {
  return async () => {
    const user = await queryClient
      .ensureQueryData(queries.profile)
      .catch(() => null);
    if (mode === "auth" && !user) throw redirect("/login");
    if (mode === "guest" && user) throw redirect("/");
  };
}
