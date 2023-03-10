-- name: CreateAccount :one
insert into accounts (
    user_id, 
    category_id, 
    title,
    type,
    description,
    value,
    date
) values (
    $1, $2, $3, $4, $5, $6, $7
) returning *;

-- name: GetAccount :one
select * from accounts 
where id = $1 limit 1;

-- name: GetAccounts :many
select 
    a.id,
    a.user_id, 
    a.category_id, 
    a.title,
    a.type,
    a.description,
    a.value,
    a.date,
    a.created_at,
    c.title as category_title
from accounts a
left join categories c on c.id = a.category_id 
where a.user_id = $1 and a.type = $2
and lower(a.title) like concat('%', lower(@title::text), '%')
and lower(a.description) like concat('%', lower(@description::text), '%')
and a.category_id = coalesce(@category_id, a.category_id)
and a.date = coalesce(@date, a.date);

-- name: GetAccountReports :one
select sum(value) as sum_value from accounts 
where user_id = $1 and type = $2;

-- name: GetAccountGraph :one
select count(*) from accounts 
where user_id = $1 and type = $2;

-- name: UpdateAccount :one
update accounts set title = $2, description = $3, value = $4 where id = $1 returning *;

-- name: DeleteAccount :exec
delete from accounts where id = $1;