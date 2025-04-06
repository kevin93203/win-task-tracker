# Windows Task Tracker

Windows Task Tracker is a web application designed for monitoring and managing Windows scheduled tasks. It provides an intuitive user interface that allows users to easily view and track scheduled task status across multiple Windows computers.

## Features

- 🔐 User authentication and authorization
- 🖥️ Support for monitoring multiple Windows computers
- 🔑 Credential management for remote computers
- 📊 Real-time display of scheduled task status
- 🔍 Task name search and filtering capabilities
- 🕒 Display of task execution times and results
- 📝 Detailed task information display
- 🎯 Task trigger and execution command details
- 🔄 Auto-refresh and manual refresh options
- ⚙️ Advanced task management (disable, enable, start, stop)

## Technology Stack

### Frontend
- Vue.js 3.5.13 with Composition API
- Vue Router 4.5.0 for navigation
- Tailwind CSS 3.4.1 for styling
- Vite 6.2.0 as build tool
- Axios 1.8.4 for HTTP requests
- Vue-Toastification for notifications

### Backend
- Go for the server implementation
- SQLite for persistent data storage
- JWT for authentication
- PowerShell Scripts for task management
- Built-in CORS support
- RESTful API architecture

## Installation Guide

### Prerequisites
- Node.js 16+
- Go 1.16+
- Windows System (with PowerShell support)
- SQLite3

### Backend Setup

1. Navigate to backend directory:
```bash
cd backend
```

2. Install Go dependencies:
```bash
go mod download
```

3. Start the backend service:
```bash
go run main.go
```

The backend will create a SQLite database file automatically and listen on port 8080.

### Frontend Setup

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start development server:
```bash
npm run dev
```

4. Build for production:
```bash
npm run build
```

## Usage Instructions

1. Ensure the backend service is running (default port: 8080)
2. Access the frontend page in your browser (default development URL: http://localhost:5173)
3. Register an account or login with existing credentials
4. Add remote computers and their credentials in the Remote Computers section
5. View and manage scheduled tasks:
   - Use the search bar to find specific tasks
   - Filter tasks by status (Ready, Running, Disabled)
   - Filter tasks by computer name
   - View detailed task information including triggers and commands
   - Control tasks with disable, enable, start, and stop actions
   - Click the refresh button to update data

## API Endpoints

### Authentication
- `/api/register` - Register new user
- `/api/login` - User login
- `/api/verify` - Verify authentication token
- `/api/logout` - User logout

### Tasks Management
- `/api/tasks` - Get all tasks
- `/api/tasks/disable` - Disable a task
- `/api/tasks/enable` - Enable a task
- `/api/tasks/start` - Start a task
- `/api/tasks/stop` - Stop a task

### Remote Computer Management
- `/api/computers/list` - List all user computers
- `/api/computers` - Create a new remote computer
- `/api/computers/delete` - Delete a remote computer
- `/api/credentials/list` - List user credentials
- `/api/credentials` - Create a new credential
- `/api/credentials/update` - Update a credential
- `/api/credentials/delete` - Delete a credential
- `/api/computers/map-credential` - Map credential to computer

## Important Notes

- Ensure appropriate Windows permissions to read scheduled tasks
- Credentials are stored securely in the SQLite database
- Task information is fetched using PowerShell remoting
- HTTPS is recommended for production deployment
- Firewall settings may need to be adjusted for required ports:
  - Backend API: 8080
  - Frontend dev server: 5173
  - WinRM: 5985 (HTTP) or 5986 (HTTPS)

## Security Considerations

- JWT tokens are used for authentication
- Passwords are hashed before storage
- HTTP-only cookies for token storage
- CORS is configured for security
- User authentication required for all operations
- Credential access is restricted to owner

## License

MIT License

## Author

kevin93203
