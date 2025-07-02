const API_URL = '/api/tasks';

const tasksList = document.getElementById('tasks-list');
const form = document.getElementById('task-form');
const titleInput = document.getElementById('task-title');
const descInput = document.getElementById('task-desc');

function createTaskElement(task) {
    const div = document.createElement('div');
    div.className = 'task-item' + (task.done ? ' done' : '');
    div.dataset.id = task.id;

    const checkbox = document.createElement('div');
    checkbox.className = 'task-checkbox' + (task.done ? ' checked' : '');
    checkbox.title = task.done ? 'Выполнено' : 'Отметить как выполнено';
    checkbox.onclick = () => updateTask(task.id, title.textContent, desc.textContent, !task.done);
    const checkboxInner = document.createElement('div');
    checkboxInner.className = 'task-checkbox-inner';
    checkbox.appendChild(checkboxInner);

    const content = document.createElement('div');
    content.className = 'task-content';

    const title = document.createElement('div');
    title.className = 'task-title';
    title.textContent = task.title;
    title.contentEditable = true;
    title.spellcheck = false;
    title.addEventListener('blur', () => updateTask(task.id, title.textContent, desc.textContent, task.done));

    const desc = document.createElement('div');
    desc.className = 'task-desc';
    desc.textContent = task.description || '';
    desc.contentEditable = true;
    desc.spellcheck = false;
    desc.addEventListener('blur', () => updateTask(task.id, title.textContent, desc.textContent, task.done));

    const time = document.createElement('div');
    time.className = 'task-time';
    time.textContent = formatTime(task.created_at);

    content.appendChild(title);
    content.appendChild(desc);
    content.appendChild(time);

    const actions = document.createElement('div');
    actions.className = 'task-actions';
    const delBtn = document.createElement('button');
    delBtn.className = 'delete-btn';
    delBtn.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-x" viewBox="0 0 16 16"><path d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708"/></svg>';
    delBtn.title = 'Удалить';
    delBtn.onclick = () => deleteTask(task.id);
    actions.appendChild(delBtn);

    div.appendChild(checkbox);
    div.appendChild(content);
    div.appendChild(actions);
    return div;
}

function formatTime(iso) {
    if (!iso) return '';
    const d = new Date(iso);
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) +
        ' · ' + d.toLocaleDateString([], { day: '2-digit', month: 'short' });
}

function renderTasks(tasks) {
    tasksList.innerHTML = '';
    tasks.forEach(task => {
        tasksList.appendChild(createTaskElement(task));
    });
}

async function fetchTasks() {
    const res = await fetch(API_URL);
    if (!res.ok) return;
    const tasks = await res.json();
    renderTasks(tasks.sort((a, b) => new Date(b.created_at) - new Date(a.created_at)));
}

async function createTask(title, description) {
    const res = await fetch(API_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, description })
    });
    if (res.ok) fetchTasks();
}

async function updateTask(id, title, description, done) {
    await fetch(`${API_URL}/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, description, done })
    });
    fetchTasks();
}

async function deleteTask(id) {
    await fetch(`${API_URL}/${id}`, { method: 'DELETE' });
    const el = document.querySelector(`.task-item[data-id='${id}']`);
    if (el) el.remove();
}

form.addEventListener('submit', e => {
    e.preventDefault();
    const title = titleInput.value.trim();
    const desc = descInput.value.trim();
    if (!title) return;
    createTask(title, desc);
    titleInput.value = '';
    descInput.value = '';
});

document.addEventListener('DOMContentLoaded', fetchTasks);
