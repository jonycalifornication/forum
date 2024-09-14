package database

import "log"

func createTables() {
	userTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	_, err := DB.Exec(userTable)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	postTable := `CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		title TEXT NOT NULL,
		text TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		image_path TEXT
	);`

	_, err = DB.Exec(postTable)
	if err != nil {
		log.Fatalf("Failed to create posts table: %v", err)
	}

	commentTable := `CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		text TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err = DB.Exec(commentTable)
	if err != nil {
		log.Fatalf("Failed to create comments table: %v", err)
	}

	categoryTable := `CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);`

	_, err = DB.Exec(categoryTable)
	if err != nil {
		log.Fatalf("Failed to create categories table: %v", err)
	}

	addInitialCategories()

	postCategoryTable := `CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		PRIMARY KEY(post_id, category_id),
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(category_id) REFERENCES categories(id)
	);`

	_, err = DB.Exec(postCategoryTable)
	if err != nil {
		log.Fatalf("Failed to create post_categories table: %v", err)
	}

	postReactionTable := `CREATE TABLE IF NOT EXISTS post_reactions (
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		reaction_type TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY(post_id, user_id),
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err = DB.Exec(postReactionTable)
	if err != nil {
		log.Fatalf("Failed to create post_reactions table: %v", err)
	}

	// Создание таблицы реакций на комментарии
	commentReactionTable := `CREATE TABLE IF NOT EXISTS comment_reactions (
		comment_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		reaction_type TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY(comment_id, user_id),
		FOREIGN KEY(comment_id) REFERENCES comments(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err = DB.Exec(commentReactionTable)
	if err != nil {
		log.Fatalf("Failed to create comment_reactions table: %v", err)
	}
}
