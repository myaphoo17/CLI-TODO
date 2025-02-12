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
		printMainMenu()

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			fmt.Println("Invalid input! Please enter a valid number.")
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
			fmt.Println("Exiting... Goodbye!")
			return
		default:
			fmt.Println("Invalid choice! Please select a valid option.")
		}
	}
}

func printMainMenu() {
	fmt.Println("===============================")
	fmt.Println("       CLI TODO Application")
	fmt.Println("===============================")
	fmt.Println("1. Add Task")
	fmt.Println("2. List Tasks")
	fmt.Println("3. Mark Task as Complete")
	fmt.Println("4. Delete Task")
	fmt.Println("5. Exit")
	fmt.Println("-------------------------------")
	fmt.Print("Choose an option: ")
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
	fmt.Println("\n>>> Task added successfully!")
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
		fmt.Println("Task not found! Please make sure the ID is correct.")
	} else {
		fmt.Println(">>> Task marked as complete successfully!")
	}
}

func deleteTask() {

	listTasks()

	fmt.Print("Enter task ID to delete: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	id, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		fmt.Println("Invalid input! Please enter a valid number.")
		return
	}

	res, err := db.Db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error deleting task:", err)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		fmt.Println("Task not found! Please make sure the ID is correct.")
	} else {

		fmt.Println("\n>>> Task deleted successfully!")
		listTasks()
	}
}

func listTasks() {
	rows, err := db.Db.Query("SELECT id, title FROM tasks")
	if err != nil {
		fmt.Println("Error fetching tasks:", err)
		return
	}
	defer rows.Close()

	fmt.Println("\n-------- Tasks List --------")
	if !rows.Next() {
		fmt.Println("No tasks found in the database.")
		return
	}

	for rows.Next() {
		var id int
		var title string

		err := rows.Scan(&id, &title)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}

		fmt.Printf("ID: %d , Title: %s\n", id, title)
	}
	fmt.Println("----------------------------")
}
