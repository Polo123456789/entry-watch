CREATE TABLE goose_db_version (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		version_id INTEGER NOT NULL,
		is_applied INTEGER NOT NULL,
		tstamp TIMESTAMP DEFAULT (datetime('now'))
	);
CREATE TABLE sqlite_sequence(name,seq);
CREATE TABLE condominiums (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    address TEXT NOT NULL,

    created_at INTEGER NOT NULL, -- Unix timestamp
    updated_at INTEGER NOT NULL,  -- Unix timestamp
    created_by INTEGER,
    updated_by INTEGER,

    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    condominium_id INTEGER,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone TEXT,
    role TEXT NOT NULL DEFAULT 'user', -- superadmin, admin, guard, user
    password TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT 0,
    hidden BOOLEAN NOT NULL DEFAULT 0, -- registerd, not approved, I dont want
                                       -- to see it in the dashboard

    created_at INTEGER NOT NULL, -- Unix timestamp
    updated_at INTEGER NOT NULL, -- Unix timestamp
    created_by INTEGER,
    updated_by INTEGER,

    FOREIGN KEY (condominium_id) REFERENCES condominiums(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE TABLE visits (
    id TEXT NOT NULL UNIQUE PRIMARY KEY,
    condominium_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    visitor_name TEXT NOT NULL,
    max_uses INTEGER NOT NULL DEFAULT 0, -- 0 means unlimited uses
    uses INTEGER NOT NULL DEFAULT 0,
    valid_from INTEGER NOT NULL, -- Unix timestamp, 0 means no restriction
    valid_to INTEGER NOT NULL, -- Unix timestamp

    created_at INTEGER NOT NULL, -- Unix timestamp
    updated_at INTEGER NOT NULL, -- Unix timestamp

    FOREIGN KEY (condominium_id) REFERENCES condominiums(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE audit_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    level INTEGER NOT NULL, -- 1=info, 2=important, 3=critical
    message TEXT NOT NULL,
    created_at INTEGER NOT NULL, -- Unix timestamp

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);
