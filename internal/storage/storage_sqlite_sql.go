package storage

const (
	// Food
	_sqlCreateTableFood = `
	CREATE TABLE IF NOT EXISTS food (
		key     TEXT NOT NULL,
		name    TEXT NOT NULL,
		brand   TEXT,
		cal100  REAL NOT NULL,
		prot100 REAL NOT NULL,
		fat100  REAL NOT NULL,
		carb100 REAL NOT NULL,
		comment TEXT,
		PRIMARY KEY (key)
	) STRICT;
	`
	_sqlGetFood = `
	SELECT key, name, brand, cal100, prot100, fat100, carb100, comment
	FROM food
	WHERE key = $1		
	`
	_sqlGetFoodList = `
	SELECT key, name, brand, cal100, prot100, fat100, carb100, comment
	FROM food
	ORDER BY name
	`
	_sqFindFood = `
	SELECT key, name, brand, cal100, prot100, fat100, carb100, comment
	FROM food
	WHERE
	    go_upper(key) like '%' || $1 || '%' OR
	    go_upper(name) like '%' || $1 || '%' OR
		go_upper(brand) like '%' || $1 || '%' OR
		go_upper(comment) like '%' || $1 || '%'
	ORDER BY name
	LIMIT 10
	`
	_sqlSetFood = `
	INSERT INTO food(key, name, brand, cal100, prot100, fat100, carb100, comment)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (key) DO
	UPDATE SET name = $2, brand = $3, cal100 = $4, prot100 = $5, fat100 = $6, carb100 = $7, comment = $8
	`
	_sqlDeleteFood = `
	DELETE
	FROM food
	WHERE key = $1
	`

	// Journal
	_sqlCreateTableJournal = `
	CREATE TABLE IF NOT EXISTS journal (
		timestamp INTEGER NOT NULL
		meal      INTEGER NOT NULL
		foodkey   TEXT NOT NULL,
		foodlbl   TEXT NOT NULL,
        weight    REAL NOT NULL,
		cal       REAL NOT NULL,
		prot      REAL NOT NULL,
		fat       REAL NOT NULL,
		carb      REAL NOT NULL,
		PRIMARY KEY (timestamp, meal, foodkey)
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
	VALUES ($1, $2, $3)
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
