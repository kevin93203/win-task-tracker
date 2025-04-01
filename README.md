# Windows Task Tracker

Windows Task Tracker is a web application designed for monitoring and managing Windows scheduled tasks. It provides an intuitive user interface that allows users to easily view and track scheduled task status across multiple Windows computers.

## Features

- 🖥️ Support for monitoring multiple Windows computers
- 📊 Real-time display of scheduled task status
- 🔍 Task name search and filtering capabilities
- 🕒 Display of task execution times and results
- 📝 Detailed task information display

## Technology Stack

### Frontend
- Vue.js 3
- Tailwind CSS
- Vite
- Axios

### Backend
- Go
- PowerShell Scripts

## Installation Guide

### Prerequisites
- Node.js 16+
- Go 1.16+
- Windows System (with PowerShell support)

### Backend Setup

1. Navigate to backend directory:
```bash
cd backend
```

2. Configure connection information:
Edit `config.json`, add the computers you want to monitor:
```json
[
    {
        "userName": "your-username",
        "password": "your-password",
        "ComputerName": "your-computer-name"
    }
]
```

3. Start the backend service:
```bash
go run main.go
```

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
3. Use the search bar to find specific tasks
4. Filter tasks by status and computer
5. Click the refresh button to update data

## Important Notes

- Ensure appropriate Windows permissions to read scheduled tasks
- Password management should use environment variables or secure configuration methods
- Firewall settings may need to be adjusted for required ports

## License

MIT License

## Author

kevin93203
