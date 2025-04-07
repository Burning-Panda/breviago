-- This file is used to initialize the database.
-- It is executed when the app is started.
CREATE TABLE
    IF NOT EXISTS acronyms (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid TEXT NOT NULL UNIQUE,
        short_form TEXT NOT NULL,
        long_form TEXT NOT NULL,
        description TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        UNIQUE (short_form, long_form)
    );

CREATE INDEX IF NOT EXISTS idx_uuid ON acronyms (uuid);

CREATE TABLE
    IF NOT EXISTS acronym_categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        acronym_id INTEGER NOT NULL,
        category_id INTEGER NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (acronym_id) REFERENCES acronyms (id) ON DELETE CASCADE,
        FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid TEXT NOT NULL UNIQUE,
        name TEXT NOT NULL UNIQUE,
        description TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX IF NOT EXISTS idx_uuid ON categories (uuid);

CREATE TABLE
    IF NOT EXISTS note_joiner (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        acronym_id INTEGER,
        note_id INTEGER NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (acronym_id) REFERENCES acronyms (id) ON DELETE SET NULL,
        FOREIGN KEY (note_id) REFERENCES notes (id) ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS notes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid TEXT NOT NULL UNIQUE,
        note TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX IF NOT EXISTS idx_uuid ON notes (uuid);