-- name: CheckIfClusterExists :one
SELECT
    EXISTS(
        SELECT
            1
        FROM
            clusters
        WHERE
            name = ?
        LIMIT
            1
    );

-- name: CheckIfNodeExists :one
SELECT
    EXISTS(
        SELECT
            1
        FROM
            nodes
        WHERE
            cluster_id = ?
            AND host = ?
            AND port = ?
        LIMIT
            1
    );

-- name: CheckIfClassExists :one
SELECT
    EXISTS(
        SELECT
            1
        FROM
            classes
        WHERE
            cluster_id = ?
            AND name = ?
        LIMIT
            1
    );

-- name: GetNodes :many
SELECT
    node_id,
    cluster_id,
    protocol,
    auth_user,
    host,
    port,
    weight,
    max_groups
FROM
    nodes
ORDER BY
    cluster_id;

-- name: GetNodesFromClusterID :many
SELECT
    node_id,
    cluster_id,
    protocol,
    auth_user,
    host,
    port,
    weight,
    max_groups
FROM
    nodes
WHERE
    cluster_id = ?;

-- name: GetClusters :many
SELECT
    cluster_id,
    name,
    description
FROM
    clusters
ORDER BY
    cluster_id;

-- name: GetClasses :many
SELECT
    class_id,
    cluster_id,
    name,
    description
FROM
    classes
WHERE
    cluster_id = ?
ORDER BY
    name;

-- name: GetNodeGroupNamesForClass :many
SELECT
    n.protocol || '://' || n.host || ':' || CAST(n.port AS TEXT) AS node_url,
    g.name AS group_name,
    u.username,
    u.full_name,
    u.default_password
FROM
    classes c
    JOIN groups g ON g.class_id = c.class_id
    JOIN users u ON u.group_id = g.group_id
    JOIN group_assignments ga ON ga.group_id = g.group_id
    JOIN nodes n ON n.node_id = ga.node_id
WHERE
    c.cluster_id = ?
    AND c.name = ?
ORDER BY
    n.node_id,
    g.group_id,
    u.user_id;

-- name: GetClusterByID :one
SELECT
    cluster_id,
    name,
    description
FROM
    clusters
WHERE
    name = ?
LIMIT
    1;

-- name: GetClassDisribution :many
SELECT
    c.name AS class_name,
    c.description AS class_desc,
    n.protocol || '://' || n.host || ':' || CAST(n.port AS TEXT) AS node_url,
    g.name AS group_name,
    u.username AS username
FROM
    classes c
    JOIN groups g ON g.class_id = c.class_id
    JOIN users u ON u.group_id = g.group_id
    JOIN group_assignments ga ON ga.group_id = g.group_id
    JOIN nodes n ON n.node_id = ga.node_id
WHERE
    c.cluster_id = ?
ORDER BY
    c.name,
    n.node_id,
    g.group_id,
    u.user_id;

-- name: GetAllExerciseNameFromCluster :many
SELECT
    DISTINCT(e.name)
FROM
    exercises e
    JOIN group_assignments ga ON e.group_id = ga.group_id
    JOIN groups g ON e.group_id = g.group_id
    JOIN classes c ON g.class_id = c.class_id
    JOIN clusters ca ON c.cluster_id = ca.cluster_id
WHERE
    ca.name = ?;

-- name: GetNodeExercisesForCluster :many
SELECT
    n.protocol || '://' || n.host || ':' || CAST(n.port AS TEXT) AS node_url,
    e.name AS exercise_name,
    e.project_uuid,
    g.name AS group_name,
    e.state,
    c.class_id,
    c.name
FROM
    classes c
    JOIN groups g ON g.class_id = c.class_id
    JOIN exercises e ON e.group_id = g.group_id
    JOIN group_assignments ga ON ga.group_id = g.group_id
    JOIN nodes n ON n.node_id = ga.node_id
WHERE
    c.cluster_id = ?
    AND e.state <> 'deleted'
ORDER BY
    n.node_id,
    g.group_id,
    e.exercise_id;

-- name: CheckClusterExistsWithData :one
SELECT
    cluster_id,
    name,
    description
FROM
    clusters
WHERE
    name = ?
LIMIT
    1;

-- name: GetNodeGroupAssignments :many
SELECT
    ga.node_id,
    COUNT(*)
FROM
    group_assignments ga
    JOIN nodes n ON n.node_id = ga.node_id
WHERE
    n.cluster_id = ?
GROUP BY
    ga.node_id;

-- name: GetExercisesForDeletion :many
SELECT
    e.project_uuid,
    e.name,
    c.name AS class_name,
    g.name AS group_name
FROM
    exercises e
    JOIN groups g ON e.group_id = g.group_id
    JOIN classes c ON g.class_id = c.class_id
WHERE
    e.name = ?
    AND (
        ? = ''
        OR c.name = ?
    )
    AND (
        ? = ''
        OR g.name = ?
    );
