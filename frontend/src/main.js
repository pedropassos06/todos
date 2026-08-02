import "./style.css";

const runtimeConfig = window.__APP_CONFIG__ || {};
const buildApiBaseUrl = import.meta.env.VITE_API_BASE_URL;
const browserHost = window.location.hostname || "localhost";
const defaultApiBaseUrl = `http://${browserHost}:8081`;
const apiBaseUrl = String(runtimeConfig.API_BASE_URL || buildApiBaseUrl || defaultApiBaseUrl).replace(/\/$/, "");

const form = document.querySelector("#todo-form");
const input = document.querySelector("#todo-input");
const list = document.querySelector("#todo-list");
const status = document.querySelector("#status");
const filterButtons = [...document.querySelectorAll(".filter")];

let todos = [];
let activeFilter = "all";

form.addEventListener("submit", async (event) => {
  event.preventDefault();

  const title = input.value.trim();
  if (!title) {
    setStatus("Digite um titulo para continuar.", true);
    return;
  }

  try {
    setStatus("Criando tarefa...");
    await request("/todos", {
      method: "POST",
      body: JSON.stringify({ title })
    });

    input.value = "";
    setStatus("Tarefa criada com sucesso.");
    await loadTodos();
  } catch (error) {
    setStatus(error.message, true);
  }
});

filterButtons.forEach((button) => {
  button.addEventListener("click", () => {
    activeFilter = button.dataset.filter;

    filterButtons.forEach((candidate) => {
      const selected = candidate === button;
      candidate.classList.toggle("is-active", selected);
      candidate.setAttribute("aria-selected", String(selected));
    });

    renderTodos();
  });
});

async function loadTodos() {
  try {
    setStatus("Carregando tarefas...");
    const data = await request("/todos", { method: "GET" });
    todos = Array.isArray(data.todos) ? data.todos : [];
    renderTodos();
    setStatus(`Total: ${todos.length} tarefa(s).`);
  } catch (error) {
    setStatus(error.message, true);
  }
}

async function toggleTodo(id, completed) {
  try {
    await request(`/todos/${id}`, {
      method: "PATCH",
      body: JSON.stringify({ completed })
    });
    await loadTodos();
  } catch (error) {
    setStatus(error.message, true);
  }
}

async function removeTodo(id) {
  try {
    await request(`/todos/${id}`, { method: "DELETE" });
    setStatus("Tarefa removida.");
    await loadTodos();
  } catch (error) {
    setStatus(error.message, true);
  }
}

function renderTodos() {
  const visible = todos.filter((todo) => {
    if (activeFilter === "open") return !todo.completed;
    if (activeFilter === "done") return todo.completed;
    return true;
  });

  if (visible.length === 0) {
    list.innerHTML = '<li class="empty">Nenhuma tarefa nesse filtro.</li>';
    return;
  }

  list.innerHTML = "";

  visible.forEach((todo) => {
    const item = document.createElement("li");
    item.className = "todo-item";
    if (todo.completed) {
      item.classList.add("done");
    }

    const title = document.createElement("p");
    title.className = "todo-title";
    title.textContent = todo.title;

    const meta = document.createElement("p");
    meta.className = "todo-meta";
    meta.textContent = `Atualizada em ${formatDate(todo.updatedAt)}`;

    const controls = document.createElement("div");
    controls.className = "todo-controls";

    const toggleButton = document.createElement("button");
    toggleButton.className = "toggle";
    toggleButton.textContent = todo.completed ? "Reabrir" : "Concluir";
    toggleButton.addEventListener("click", () => toggleTodo(todo.id, !todo.completed));

    const deleteButton = document.createElement("button");
    deleteButton.className = "danger";
    deleteButton.textContent = "Excluir";
    deleteButton.addEventListener("click", () => removeTodo(todo.id));

    controls.append(toggleButton, deleteButton);
    item.append(title, meta, controls);
    list.append(item);
  });
}

function setStatus(message, isError = false) {
  status.textContent = message;
  status.classList.toggle("error", isError);
}

function formatDate(value) {
  if (!value) return "agora";

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "agora";

  return new Intl.DateTimeFormat("pt-BR", {
    dateStyle: "short",
    timeStyle: "short"
  }).format(date);
}

async function request(path, options = {}) {
  const response = await fetch(`${apiBaseUrl}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {})
    },
    ...options
  });

  if (response.status === 204) {
    return null;
  }

  const payload = await response.json().catch(() => ({}));

  if (!response.ok) {
    const message = payload.message || "Falha na comunicacao com a API.";
    throw new Error(message);
  }

  return payload;
}

loadTodos();
