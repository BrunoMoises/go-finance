-- name: CreateCategory :one
insert into categories (
    user_id, 
    title,
    type,
    description
) values (
    $1, $2, $3, $4
) returning *;

-- name: GetCategory :one
select * from categories 
where id = $1 limit 1;

-- name: GetCategories :many
select * from categories
where user_id = $1 and type = $2
and title like concat('%', @title::text, '%')
and description like concat('%', @description::text, '%');

-- name: UpdateCategory :one
update categories 
set title = $2, description = $3 
where id = $1 returning *;

-- name: DeleteCategory :exec
delete from categories
where id = $1;