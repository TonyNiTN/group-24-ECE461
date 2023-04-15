-- Fill Users table
-- what does privilege level look like?

-- user with valid token
INSERT INTO Users (ID, USERNAME, PASSWORD, TOKEN, TOKEN_CREATED, TOKEN_EXPIRY, PRIVILEDGE_LEVEL) VALUES (
    1, 
    'valid-user',
    'strongpassword',
    'token1',
    '2023-04-14 02:20:20',
    '2024-04-14 02:20:20',
    1
);

-- user with expired token
INSERT INTO Users (ID, USERNAME, PASSWORD, TOKEN, TOKEN_CREATED, TOKEN_EXPIRY, PRIVILEDGE_LEVEL) VALUES (
    2, 
    'invalid-user',
    'strongpassword',
    'token2',
    '2022-04-14 02:20:20',
    '2023-04-14 02:20:20',
    1
);

-- Fill Binaries table # missing actual package

INSERT INTO Binaries (ID) VALUES (
    1
);

INSERT INTO Binaries (ID) VALUES (
    2
);

-- Fill Ratings table

-- Good rating
INSERT INTO Ratings (ID, BUS_FACTOR, CORRECTNESS, RAMP_UP, RESPONSIVENESS, LICENSE_SCORE, PINNING_PRACTICE, PULL_REQUEST, NET_SCORE) VALUES (
	1,
	1.0,
    1.0,
	1.0,
	1.0,
	1.0,
	1.0,
	1.0,
	1.0
);

-- Bad rating
INSERT INTO Ratings (ID, BUS_FACTOR, CORRECTNESS, RAMP_UP, RESPONSIVENESS, LICENSE_SCORE, PINNING_PRACTICE, PULL_REQUEST, NET_SCORE) VALUES (
	2,
	0.0,
    0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0,
	0.0
);

-- Fill Registry table

INSERT INTO Registry (ID, NAME, RATING_PK, AUTHOR_PK, URL, BINARY_PK, VERSION, UPLOADED, IS_EXTERNAL) VALUES (
	1,
	'package1-name',
	1,
	1, 
	'github.com/packit461/packit23',
	1,
	'2.0',
	'2023-04-14 02:20:20',
	False
);

INSERT INTO Registry (ID, NAME, RATING_PK, AUTHOR_PK, URL, BINARY_PK, VERSION, UPLOADED, IS_EXTERNAL) VALUES (
	2,
	'package2-name',
	2,
	2, 
	'github.com/19chonm/461_1_23',
	2,
	'1.5',
	'2023-04-14 02:20:20',
	False
);