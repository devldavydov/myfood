package storage

const (
	// Food
	_sqlCreateTableFood = `
	CREATE TABLE IF NOT EXISTS food (
		userid  INTEGER NOT NULL,
		name    TEXT NOT NULL,
		brand   TEXT,
		cal100  REAL NOT NULL,
		prot100 REAL NOT NULL,
		fat100  REAL NOT NULL,
		carb100 REAL NOT NULL,
		comment TEXT,
		PRIMARY KEY (userid, name, brand)
	) STRICT;
	`

	// Weight
	_sqlCreateTableWeight = `
	CREATE TABLE IF NOT EXISTS weight (
		userid    INTEGER NOT NULL,
		timestamp INTEGER NOT NULL,
		value     REAL NOT NULL,
		PRIMARY KEY (userid, timestamp)
	) STRICT;	
	`
	_sqlCreateWeight = `
	INSERT INTO weight(userid, timestamp, value)
	VALUES($1, $2, $3)
	`
	_sqlFindWeight = `
	SELECT timestamp, value
	FROM weight
	WHERE userid = $1 and timestamp = $2
	`
)
