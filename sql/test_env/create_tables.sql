CREATE TABLE Users ( 
	ID int, 
	USERNAME varchar(50),
	PASSWORD varchar(50),
	TOKEN varchar(50),
	TOKEN_CREATED datetime,
	TOKEN_EXPIRY datetime,
	PRIVILEGE_LEVEL int,
	PRIMARY KEY (ID)
);

-- missing BINARY_FILE
CREATE TABLE Binaries (
	ID int,
	PRIMARY KEY (ID)
);

CREATE TABLE Ratings (
	ID int,
	BUS_FACTOR float,
	CORRECTNESS float,
	RAMP_UP float,
	RESPONSIVENESS float,
	LICENSE_SCORE float,
	PINNING_PRACTICE float,
	PULL_REQUEST float,
	NET_SCORE float,
	PRIMARY KEY (ID)
);

CREATE TABLE Registry (
	ID int,
	NAME varchar (50),
	RATING_PK int,
	AUTHOR_PK int, 
	URL varchar(255),
	BINARY_PK int,
	VERSION varchar (15),
	UPLOADED datetime,
	IS_EXTERNAL boolean,
	PRIMARY KEY (ID),
	FOREIGN KEY(RATING_PK) REFERENCES Ratings(ID),
	FOREIGN KEY(AUTHOR_PK) REFERENCES Users(ID),
	FOREIGN KEY(BINARY_PK) REFERENCES Binaries(ID)
);