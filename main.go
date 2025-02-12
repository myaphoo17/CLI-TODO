package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"CLI-TODO/db"

	"gopkg.in/ini.v1"
)

func main() {
	config, err := loadConfig("config.ini")
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		fmt.Println("Please create a config.ini file with the following format:")
		fmt.Println("\n[database]\ndsn = root:password@tcp(127.0.0.1:3306)/task_manager")
		os.Exit(1)
	}
	db.InitDB(config["dsn"])
	defer db.Db.Close()

	db.CreateAndInitializeDatabase()

	for {
		fmt.Println("\n1. Add Task\n2. List Tasks\n3. Mark Task as Complete\n4. Delete Task\n5. Exit")
		fmt.Print("Choose an option: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			fmt.Println("Invalid input! Please enter a number.")
			continue
		}

		switch choice {
		case 1:
			addTask()
		case 2:
			listTasks()
		case 3:
			markComplete()
		case 4:
			deleteTask()
		case 5:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice! Please select a valid option.")
		}
	}
}

func loadConfig(filename string) (map[string]string, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found")
	}
	cfg, err := ini.Load(filename)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"dsn": cfg.Section("database").Key("dsn").String(),
	}, nil
}

func addTask() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter task title: ")
	title, _ := reader.ReadString('\n')
	fmt.Print("Enter task description: ")
	description, _ := reader.ReadString('\n')

	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	_, err := db.Db.Exec("INSERT INTO tasks (title, description, completed) VALUES (?, ?, ?)", title, description, false)
	if err != nil {
		fmt.Println("Error adding task:", err)
		return
	}
	fmt.Println("Task added successfully!")
}

func listTasks() {
	rows, err := db.Db.Query("SELECT id, title, description, completed FROM tasks")
	if err != nil {
		fmt.Println("Error fetching tasks:", err)
		return
	}
	defer rows.Close()

	fmt.Println("\nTasks:")
	for rows.Next() {
		var id int
		var title, description string
		var completed bool

		err := rows.Scan(&id, &title, &description, &completed)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}

		status := "Incomplete"
		if completed {
			status = "Complete"
		}
		fmt.Printf("ID: %d\nTitle: %s\nDescription: %s\nStatus: %s\n\n", id, title, description, status)
	}
}

func markComplete() {
	fmt.Print("Enter task ID to mark as complete: ")
	var id int
	fmt.Scan(&id)

	res, err := db.Db.Exec("UPDATE tasks SET completed = TRUE WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error updating task:", err)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		fmt.Println("Task not found!")
	} else {
		fmt.Println("Task marked as complete!")
	}
}

func deleteTask() {
	fmt.Print("Enter task ID to delete: ")
	var id int
	fmt.Scan(&id)

	res, err := db.Db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error deleting task:", err)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		fmt.Println("Task not found!")
	} else {
		fmt.Println("Task deleted successfully!")
	}
}
