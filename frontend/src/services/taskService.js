import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // 自動帶上 cookie
});

// 不再需要手動設定 token

// Trigger APIs
export async function addTrigger({ computer_id, task_name, cron_expression }) {
  return api.post('http://localhost:8080/api/tasks/triggers', {
    computer_id,
    task_name,
    cron_expression,
  });
}

export async function updateTrigger({ computer_id, task_name, cron_expression, index }) {
  return api.patch('http://localhost:8080/api/tasks/triggers', {
    computer_id,
    task_name,
    cron_expression,
    index,
  });
}

export async function deleteTrigger({ computer_id, task_name, index }) {
  return api.delete('http://localhost:8080/api/tasks/triggers', {
    data: {
      computer_id,
      task_name,
      index,
    },
  });
}

// Action APIs
export async function addAction({ computer_id, task_name, execute, args = '', working_directory = '' }) {
  return api.post('http://localhost:8080/api/tasks/actions', {
    computer_id,
    task_name,
    execute,
    args,
    working_directory,
  });
}

export async function updateAction({ computer_id, task_name, execute, args = '', working_directory = '', index }) {
  return api.patch('http://localhost:8080/api/tasks/actions', {
    computer_id,
    task_name,
    execute,
    args,
    working_directory,
    index,
  });
}

export async function deleteAction({ computer_id, task_name, index }) {
  return api.delete('http://localhost:8080/api/tasks/actions', {
    data: {
      computer_id,
      task_name,
      index,
    },
  });
}
