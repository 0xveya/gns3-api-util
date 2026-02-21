PRAGMA foreign_keys = ON;

CREATE TABLE clusters(
    cluster_id integer PRIMARY KEY autoincrement,
    name text NOT NULL UNIQUE,
    description text
);

CREATE TABLE nodes (
    node_id integer PRIMARY KEY autoincrement,
    cluster_id integer NOT NULL,
    protocol text CHECK (protocol IN ('http', 'https')) NOT NULL,
    auth_user text NOT NULL,
    host text NOT NULL,
    port integer NOT NULL,
    weight integer NOT NULL DEFAULT 5 CHECK (weight BETWEEN 0 AND 10),
    max_groups integer DEFAULT 3
    -- unique(protocol, host, port)
);

CREATE TABLE classes(
    class_id integer PRIMARY KEY autoincrement,
    cluster_id integer NOT NULL,
    name text NOT NULL,
    description text,
    FOREIGN KEY (cluster_id) REFERENCES clusters(cluster_id) ON DELETE CASCADE
);

CREATE TABLE groups(
    group_id integer PRIMARY KEY autoincrement,
    class_id integer NOT NULL,
    name text NOT NULL,
    FOREIGN KEY (class_id) REFERENCES classes(class_id) ON DELETE CASCADE
);

CREATE TABLE users (
    user_id integer PRIMARY KEY autoincrement,
    username text NOT NULL UNIQUE,
    full_name text,
    group_id integer NOT NULL,
    default_password text NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(group_id) ON DELETE CASCADE
);

CREATE TABLE group_assignments (
    group_id integer NOT NULL,
    node_id integer NOT NULL,
    assigned_at timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id),
    FOREIGN KEY (group_id) REFERENCES groups(group_id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES nodes(node_id) ON DELETE CASCADE
);

CREATE TABLE exercises (
    exercise_id integer PRIMARY KEY autoincrement,
    project_uuid text NOT NULL CHECK (length(project_uuid) = 8),
    group_id integer NOT NULL,
    name text NOT NULL,
    state text CHECK (
        state IN ('created', 'running', 'completed', 'deleted')
    ) DEFAULT 'created',
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(group_id) ON DELETE CASCADE
);
