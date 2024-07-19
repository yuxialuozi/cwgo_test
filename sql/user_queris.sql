
-- name: insert:exec
--api.get: /api/Users/create
INSERT INTO Users (Id, Username, Age, City, Banned, RoleId, Email, DefaultAge)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);


-- name: DeleteByID:exec
--api.delete: /api/Users/:id
DELETE FROM Users WHERE Id = ?;


-- name: UpdateUsersById:exec
--api.put: /api/Users/:id
UPDATE Users SET  Id = ?, Username = ?, Age = ?, City = ?, Banned = ?, RoleId = ?, Email = ?, DefaultAge = ?
WHERE Id = ?;


-- name: GetUsersById:many
--api.get: /api/Users/Id/:value
SELECT * FROM Users
WHERE Id = ?;


-- name: DeleteUsersById:exec
--api.delete: /api/Users/Id/:value
DELETE FROM Users
WHERE Id = ?;


-- name: GetUsersByUsername:many
--api.get: /api/Users/Username/:value
SELECT * FROM Users
WHERE Username = ?;


-- name: DeleteUsersByUsername:exec
--api.delete: /api/Users/Username/:value
DELETE FROM Users
WHERE Username = ?;


-- name: GetUsersByAge:many
--api.get: /api/Users/Age/:value
SELECT * FROM Users
WHERE Age = ?;


-- name: DeleteUsersByAge:exec
--api.delete: /api/Users/Age/:value
DELETE FROM Users
WHERE Age = ?;


-- name: GetUsersByCity:many
--api.get: /api/Users/City/:value
SELECT * FROM Users
WHERE City = ?;


-- name: DeleteUsersByCity:exec
--api.delete: /api/Users/City/:value
DELETE FROM Users
WHERE City = ?;


-- name: GetUsersByBanned:many
--api.get: /api/Users/Banned/:value
SELECT * FROM Users
WHERE Banned = ?;


-- name: DeleteUsersByBanned:exec
--api.delete: /api/Users/Banned/:value
DELETE FROM Users
WHERE Banned = ?;


-- name: GetUsersByRoleId:many
--api.get: /api/Users/RoleId/:value
SELECT * FROM Users
WHERE RoleId = ?;


-- name: DeleteUsersByRoleId:exec
--api.delete: /api/Users/RoleId/:value
DELETE FROM Users
WHERE RoleId = ?;


-- name: GetUsersByEmail:many
--api.get: /api/Users/Email/:value
SELECT * FROM Users
WHERE Email = ?;


-- name: DeleteUsersByEmail:exec
--api.delete: /api/Users/Email/:value
DELETE FROM Users
WHERE Email = ?;


-- name: GetUsersByDefaultAge:many
--api.get: /api/Users/DefaultAge/:value
SELECT * FROM Users
WHERE DefaultAge = ?;


-- name: DeleteUsersByDefaultAge:exec
--api.delete: /api/Users/DefaultAge/:value
DELETE FROM Users
WHERE DefaultAge = ?;

