-- name: DeleteClass :exec
DELETE FROM
    classes
WHERE
    class_id = ?;

-- name: UnassignGroupFromNode :exec
DELETE FROM
    group_assignments
WHERE
    node_id = ?
    AND group_id = ?;

-- name: DeleteExerciseRecord :exec
DELETE FROM
    exercises
WHERE
    project_uuid = ?;

-- name: DeleteExerciseForClass :exec
DELETE FROM
    exercises
WHERE
    group_id IN (
        SELECT
            group_id
        FROM
            groups
        WHERE
            class_id IN (
                SELECT
                    class_id
                FROM
                    classes
                WHERE
                    classes.name = ?
            )
    );

-- name: DeleteExerciseForClassOnNode :exec
DELETE FROM
    exercises
WHERE
    group_id IN (
        SELECT
            g.group_id
        FROM
            groups g
        WHERE
            g.class_id IN (
                SELECT
                    class_id
                FROM
                    classes
                WHERE
                    classes.name = ?
            )
            AND g.group_id IN (
                SELECT
                    group_id
                FROM
                    group_assignments
                WHERE
                    node_id = ?
            )
    );

-- name: DeleteNode :exec
DELETE FROM
    nodes
WHERE
    node_id = ?;

-- name: DeleteCluster :exec
DELETE FROM
    clusters
WHERE
    cluster_id = ?;

-- name: DeleteExerciseByName :exec
DELETE FROM
    exercises
WHERE
    name = ?;

-- name: DeleteExerciseScoped :exec
DELETE FROM
    exercises
WHERE
    exercises.name = ?
    AND group_id IN (
        SELECT
            group_id
        FROM
            groups
        WHERE
            class_id IN (
                SELECT
                    class_id
                FROM
                    classes
                WHERE
                    classes.name = ?
            )
    );

-- name: NukeEverything :exec
DELETE FROM
    clusters;

DELETE FROM
    nodes;

DELETE FROM
    classes;

DELETE FROM
    sqlite_sequence;
