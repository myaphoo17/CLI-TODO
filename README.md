# CLI TODO Application

This is a basic command-line tool built with Go to manage tasks, utilizing MySQL as the database backend. The app allows users to add, view, complete, and remove tasks. It also includes automatic database creation if it’s not already present.
## Features

- **Task Management**: Add, list, mark as complete, and delete tasks.
- **Database Handling**: Automatically creates the database and initializes the `tasks` table if it doesn't exist.
- **Configuration**: The database connection details and other settings are managed using a `config.ini` file.

## Prerequisites

- **Go**: The application is written in Go and requires a Go environment set up on your machine.
- **MySQL**: You need MySQL (or MariaDB) installed locally or remotely for the application to work.
- **Go Modules**: Ensure your Go project uses Go modules for package management.

## Setup Instructions

1. **Install Go (if not already installed)**:
   - Follow the instructions from the [Go official website](https://golang.org/dl/) to install Go.

2. **Install MySQL**:
   - You need a MySQL database set up either locally or remotely. 
   - If you're using a local MySQL installation, ensure it’s running on port `3306` by default.

3. **Clone the Repository**:
   ```bash
   git clone https://github.com/myaphoo17/CLI-TODO.git
   cd CLI-TODO
