-- name: GetAllAcronyms :many
SELECT
    id,
    uuid,
    short_form,
    long_form,
    description,
    created_at,
    updated_at
FROM
    acronyms;

-- name: GetAcronym :one
SELECT
    id,
    uuid,
    short_form,
    long_form,
    description,
    created_at,
    updated_at
FROM
    acronyms
WHERE
    id = ?;

-- name: SearchAcronyms :many
SELECT
    id,
    uuid,
    short_form,
    long_form,
    description,
    created_at,
    updated_at
FROM
    acronyms
WHERE
    short_form LIKE ?
    OR long_form LIKE ?;

-- name: CreateAcronym :one
INSERT INTO
    acronyms (uuid, short_form, long_form, description)
VALUES
    (?, ?, ?, ?) RETURNING id,
    uuid,
    short_form,
    long_form,
    description,
    created_at,
    updated_at;

-- name: UpdateAcronym :one
UPDATE acronyms
SET
    short_form = ?,
    long_form = ?,
    description = ?
WHERE
    id = ? RETURNING id,
    uuid,
    short_form,
    long_form,
    description,
    created_at,
    updated_at;