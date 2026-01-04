<script setup>
import { ref, onMounted } from 'vue'

const tasks = ref([])
const newTaskTitle = ref('')

const fetchTasks = async () => {
  try {
    const res = await fetch('/api/v1/tasks')
    if (res.ok) {
      tasks.value = await res.json()
    }
  } catch (e) {
    console.error("Error fetching tasks:", e)
  }
}

const createTask = async () => {
  if (!newTaskTitle.value.trim()) return
  try {
    await fetch('/api/v1/tasks', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: newTaskTitle.value })
    })
    newTaskTitle.value = ''
    await fetchTasks()
  } catch (e) {
    console.error("Error creating task:", e)
  }
}

const toggleTask = async (task) => {
  try {
    await fetch(`/api/v1/tasks/${task.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: task.title, done: !task.done })
    })
    await fetchTasks()
  } catch (e) {
    console.error("Error updating task:", e)
  }
}

const deleteTask = async (id) => {
  try {
    await fetch(`/api/v1/tasks/${id}`, { method: 'DELETE' })
    await fetchTasks()
  } catch (e) {
    console.error("Error deleting task:", e)
  }
}

onMounted(() => {
  fetchTasks()
})
</script>

<template>
  <div class="app-container">
    <h1>Todo App</h1>
    <div class="input-group">
      <input v-model="newTaskTitle" placeholder="New task..." @keyup.enter="createTask" />
      <button @click="createTask">Add</button>
    </div>
    <ul class="task-list">
      <li v-for="task in tasks" :key="task.id" :class="{ completed: task.done }">
        <span @click="toggleTask(task)" class="task-title">
          {{ task.done ? '✅' : '⬜' }} {{ task.title }}
        </span>
        <button class="delete-btn" @click="deleteTask(task.id)">🗑</button>
      </li>
    </ul>
  </div>
</template>

<style>
body { font-family: 'Segoe UI', sans-serif; background: #f4f4f9; display: flex; justify-content: center; }
.app-container { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); width: 400px; margin-top: 50px; }
h1 { text-align: center; color: #333; }
.input-group { display: flex; gap: 10px; margin-bottom: 20px; }
input { flex: 1; padding: 10px; border: 1px solid #ddd; border-radius: 4px; }
button { padding: 10px 20px; background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
.task-list { list-style: none; padding: 0; }
li { display: flex; justify-content: space-between; align-items: center; padding: 10px; border-bottom: 1px solid #eee; }
.task-title { cursor: pointer; flex: 1; text-align: left; }
.completed .task-title { text-decoration: line-through; color: #888; }
.delete-btn { background: #ff4d4d; padding: 5px 10px; border: none; color: white; cursor: pointer; }
</style>