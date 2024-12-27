package repository

const (
	createUserQuery = `
		INSERT INTO users 
			(email,username,fullname,password,password_change_at)
		VALUES
			($1,$2,$3,$4,$5)
		RETURNING id,email,username,fullname,password,password_change_at,created_at;
	`

	findUserByIDQuery = `
		SELECT
			id,email,username,fullname,password,password_change_at,created_at
		FROM
			users
		WHERE
			id = $1
	`

	findUserByEmailQuery = `
		SELECT 
			id,email,username,fullname,password,password_change_at,created_at
		FROM 
			users
		WHERE
			email = $1;
	`
)
