import { mount } from "svelte";
import App from "./App.svelte";

const target = document.getElementById("app")!;
target.innerText = "";
const app = mount(App, { target });

export default app;
