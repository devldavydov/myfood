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
	_sqlSetWeight = `
	INSERT into weight(userid, timestamp, value)
	VALUES($1, $2, $3)
	ON CONFLICT (userid, timestamp) DO
	UPDATE SET value = $3
	`
	_sqlWeightList = `
	SELECT timestamp, value
	FROM weight
	WHERE userid = $1 AND 
	    timestamp >= $2 AND
		timestamp <= $3
	ORDER BY timestamp ASC
	`
	_sqlDeleteWeight = `
	DELETE
	FROM weight
	WHERE userid = $1 AND timestamp = $2
	`

	// UserSettings
	_sqlCreateTableUserSettings = `
	CREATE TABLE IF NOT EXISTS user_settings (
		userid    INTEGER NOT NULL,
		cal_limit REAL NOT NULL,
		PRIMARY KEY (userid)
	) STRICT;
	`
	_sqlGetUserSettings = `
	SELECT cal_limit
	FROM user_settings
	WHERE userid = $1
	`
	_sqlSetUserSettings = `
	INSERT INTO user_settings(userid, cal_limit)
	VALUES ($1, $2)
	ON CONFLICT (userid) DO
	UPDATE set cal_limit = $2
	`
)
