package postgre

const MigrationQuery = `
	CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		login VARCHAR(255) UNIQUE NOT NULL,
		pwd_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		last_update TIMESTAMP DEFAULT NOW()
	);
	-- Таблица логинов и паролей для сайтов
	CREATE TABLE logins_data (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		service_name VARCHAR(255) NOT NULL,
		login VARCHAR(255) NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		last_update TIMESTAMP DEFAULT NOW(),
		metadata TEXT
	);
	
	CREATE TABLE text_data (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		data TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		last_update TIMESTAMP DEFAULT NOW(),
		metadata TEXT
	);
	
	CREATE TABLE credit_cards_data (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		card_id VARCHAR(255) NOT NULL,
		owner_name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		last_update TIMESTAMP DEFAULT NOW(),
		metadata TEXT
	);
	
	CREATE TABLE files_data (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		path VARCHAR(1024) NOT NULL,
		file_type VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		last_update TIMESTAMP DEFAULT NOW(),
		metadata TEXT
	);
`
