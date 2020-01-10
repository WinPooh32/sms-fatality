-- name: GetMessages :many
SELECT * FROM message
WHERE phone = $1;

-- name: ListPhones :many
SELECT * FROM message
ORDER BY phone;

-- name: CreateMessage :exec
INSERT INTO message (
  phone, body
) VALUES (
  $1, $2
);

-- name: DeleteMessage :exec
DELETE FROM message
WHERE id = $1;
