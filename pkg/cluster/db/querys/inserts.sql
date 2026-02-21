-- name: CreateCluster :one
INSERT INTO
    clusters (name, description)
VALUES
    (?, ?)
RETURNING
    cluster_id,
    name,
    description;

-- name: InsertNodeIntoCluster :one
INSERT INTO
    nodes (
        cluster_id,
        protocol,
        host,
        port,
        weight,
        max_groups,
        auth_user
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?)
RETURNING
    node_id,
    cluster_id,
    protocol,
    auth_user,
    host,
    port,
    weight,
    max_groups;

-- name: InsertNode :exec
INSERT INTO
    nodes (
        cluster_id,
        protocol,
        host,
        port,
        weight,
        max_groups,
        auth_user
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?);

-- name: AssignGroupToNode :exec
INSERT INTO
    group_assignments (node_id, group_id)
VALUES
    (?, ?);

-- name: CreateClass :exec
INSERT INTO
    classes (cluster_id, name, description)
VALUES
    (?, ?, ?);

-- name: CreateGroup :exec
INSERT INTO
    groups (class_id, name)
VALUES
    (?, ?);

-- name: CreateUser :exec
INSERT INTO
    users (group_id, username, full_name, default_password)
VALUES
    (?, ?, ?, ?);

-- name: InsertExerciseRecord :exec
INSERT INTO
    exercises (project_uuid, group_id, name, state)
VALUES
    (?, ?, ?, ?);

-- name: CreateClassReturning :one
INSERT INTO
    classes (cluster_id, name, description)
VALUES
    (?, ?, ?)
RETURNING
    class_id;

-- name: CreateGroupReturning :one
INSERT INTO
    groups (class_id, name)
VALUES
    (?, ?)
RETURNING
    group_id;

-- name: CreateUserReturning :one
INSERT INTO
    users (group_id, username, full_name, default_password)
VALUES
    (?, ?, ?, ?)
RETURNING
    user_id;
