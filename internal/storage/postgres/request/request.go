package request

const (
	SaveUser = "INSERT INTO users (email, pass_hash) VALUES ($1, $2) RETURNING id"
	User     = "SELECT * FROM users WHERE email = $1"
	App      = "SELECT * FROM apps WHERE id = $1"
	IsAdmin  = "SELECT * FROM admins WHERE id = $1"
)
