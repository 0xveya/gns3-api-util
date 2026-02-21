-- name: UpdateNode :exec
UPDATE
    nodes
SET
    protocol = ?,
    auth_user = ?,
    weight = ?,
    max_groups = ?
WHERE
    cluster_id = ?
    AND host = ?
    AND port = ?;

-- name: UpdateClusterDescription :exec
UPDATE
    clusters
SET
    description = ?
WHERE
    cluster_id = ?;
