aws-server-es1=(@test\\.com|@test\\.co)$; # mongodb delete queries require you use regex
digital-ocean-nyc1=DELETE FROM users WHERE email LIKE "%test.com" OR email LIKE "%test.co";
lenode-es1=DELETE FROM users WHERE email LIKE "%test.com" OR email LIKE "%test.co";
aws-server-ws1=DELETE FROM users WHERE email LIKE "%test.com" OR email LIKE "%test.co";
