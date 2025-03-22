
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS account_profit;
DROP TABLE IF EXISTS phone_numbers;
DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS payment_transactions;
DROP TABLE IF EXISTS closures;
DROP TABLE IF EXISTS payouts;
DROP TABLE IF EXISTS closure_transactions;
DROP TABLE IF EXISTS closure_payments;
DROP TABLE IF EXISTS liquidations;
DROP TABLE IF EXISTS liquidation_payments;
DROP TABLE IF EXISTS liquidation_transactions;

-- usuarios
CREATE TABLE IF NOT EXISTS users (
	user_id CHAR(8) UNIQUE NOT NULL, --Generated via trigger & sequence ('AF-' || LPAD(nexval('user_seq'), 5, '0'))
	password VARCHAR(25) NOT NULL CHECK (LENGTH(password) > 5),
	first_name VARCHAR(50) NOT NULL,
	second_name VARCHAR(50),
	first_lastname VARCHAR(50) NOT NULL,
	second_lastname VARCHAR(50),
	address_house_number VARCHAR(50),
	address_street VARCHAR(50),
	address_avenue VARCHAR(50),
	address_city VARCHAR(50),
	address_department VARCHAR(50),
	address_reference VARCHAR(150),
	primary_email VARCHAR(100) NOT NULL,
	secondary_email VARCHAR(100),
	birth_date DATE,
    is_active BOOLEAN DEFAULT TRUE,
    is_admin BOOLEAN DEFAULT FALSE,
	hiring_date DATE NOT NULL DEFAULT CURRENT_DATE,
	created_by VARCHAR(101),
	creation_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified_by VARCHAR(101),
	last_modification_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS phone_numbers (
	user_id CHAR(8) NOT NULL,
	user_phone_number INT NOT NULL UNIQUE,
	region_number INT,
	PRIMARY KEY(user_id, user_phone_number),
 	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- prestamos
CREATE TABLE IF NOT EXISTS loans (
	user_id CHAR(8) NOT NULL,
	loan_id CHAR(16) UNIQUE NOT NULL, -- Generated via trigger(user_id + "-PT" + LPAD(nexval(loan_seq), 5, '0')) - done
	loan_periods INT CHECK (loan_periods <= 12),
	loan_interest NUMERIC(3,3) NOT NULL DEFAULT 0.15 CHECK (loan_interest BETWEEN 0 AND 1),
    requested_amount NUMERIC(8,2) NOT NULL CHECK (REQUESTED_AMOUNT >= 120),
	loan_date DATE NOT NULL DEFAULT CURRENT_DATE,
	is_payed BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (loan_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- pagos
CREATE TABLE IF NOT EXISTS payments (
	payment_id INT NOT NULL UNIQUE GENERATED ALWAYS AS IDENTITY,
	loan_id CHAR(16) NOT NULL,
	payment_number char(5) NOT NULL, -- Generated via trigger - done
	deadline DATE NOT NULL,
	interest_to_pay NUMERIC(8,2),
	capital_to_pay NUMERIC(8,2),
	amount_to_pay NUMERIC(8,2) NOT NULL GENERATED ALWAYS AS (interest_to_pay + capital_to_pay),
	is_payed BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (loan_id, payment_number),
	FOREIGN KEY (loan_id) REFERENCES loans(loan_id)
);

-- cuenta
CREATE TABLE IF NOT EXISTS accounts (
	user_id CHAR(8) NOT NULL,
	account_type char(3) NOT NULL CHECK(account_type IN ('CAP', 'CAR')),
	account_id char(12) UNIQUE NOT NULL GENERATED ALWAYS AS (user_id || '-' || account_type),
	balance NUMERIC(8,2) NOT NULL DEFAULT 0 CHECK(balance >= 0),
	created_by VARCHAR(101) NOT NULL,
	creation_date DATE NOT NULL DEFAULT CURRENT_DATE,
	modified_by VARCHAR(101),
	last_modification_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY(account_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Ganancias de los dividendos a nombre de la cuenta
CREATE TABLE IF NOT EXISTS account_profit (
	account_id CHAR(12) UNIQUE NOT NULL,
	profit DECIMAL(8, 2) NOT NULL DEFAULT 0,
	PRIMARY KEY (account_id),
	FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

-- abonos
CREATE TABLE IF NOT EXISTS transactions (
	account_id CHAR(12) NOT NULL,
	transaction_id VARCHAR(20) UNIQUE NOT NULL,
	transaction_date DATE NOT NULL DEFAULT CURRENT_DATE,
	transaction_ammount NUMERIC(8,2) NOT NULL CHECK(transaction_ammount > 0),
	transaction_comment VARCHAR(255),
	PRIMARY KEY (transaction_id),
	FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

-- Cierre mensual
CREATE TABLE IF NOT EXISTS closures (
	closure_id INT UNIQUE NOT NULL GENERATED ALWAYS AS IDENTITY,
	closure_month INT NOT NULL,
	closure_year INT NOT NULL,
	description VARCHAR(255) NOT NULL,
	PRIMARY KEY(closure_id)
);

-- transacciones de pago (realizadas por el usuario)
CREATE TABLE IF NOT EXISTS payment_transactions (
	payment_id INT NOT NULL,
	transaction_id VARCHAR(20) NOT NULL,
	PRIMARY KEY (transaction_id, payment_id),
	FOREIGN KEY (payment_id) REFERENCES payments(payment_id),
	FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id)	
);

-- Transacciones de cierre
CREATE TABLE IF NOT EXISTS closure_transactions (
	closure_id INT NOT NULL,
	transaction_id VARCHAR(20) NOT NULL,
	PRIMARY KEY (transaction_id, closure_id),
	FOREIGN KEY (closure_id) REFERENCES closures(closure_id),
	FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id)
);

-- pagos de cierre
CREATE TABLE IF NOT EXISTS closure_payments (
	closure_id INT NOT NULL,
	payment_id INT NOT NULL,
	PRIMARY KEY (payment_id, closure_id),
	FOREIGN KEY (closure_id) REFERENCES closures(closure_id),
	FOREIGN KEY (payment_id) REFERENCES payments(payment_id)
);

-- Dividendos/ganancias
CREATE TABLE IF NOT EXISTS payouts (
	payout_id INT NOT NULL UNIQUE GENERATED ALWAYS AS IDENTITY,
	closure_id INT,
	account_id CHAR(12),
	account_balance NUMERIC(8,2),
	apportation_percentage INT,
	account_profit NUMERIC(8,2),
	PRIMARY KEY (payout_id),
	FOREIGN KEY (closure_id) REFERENCES closures(closure_id),
	FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

CREATE TABLE IF NOT EXISTS liquidations (
	liquidation_id INT NOT NULL UNIQUE GENERATED ALWAYS AS IDENTITY,
	account_id CHAR(12) NOT NULL,
	liquidation_type CHAR(1) NOT NULL CHECK (liquidation_type IN ('T', 'P')),
	retirement_date DATE NOT NULL DEFAULT CURRENT_DATE,
	total_money DECIMAL(8, 2),
	PRIMARY KEY(liquidation_id),
	FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

CREATE TABLE IF NOT EXISTS liquidation_payments (
	liquidation_id INT NOT NULL,
	payment_id INT NOT NULL,
	PRIMARY KEY (liquidation_id),
	FOREIGN KEY (liquidation_id) REFERENCES liquidations(liquidation_id),
	FOREIGN KEY (payment_id) REFERENCES payments(payment_id)
);

CREATE TABLE IF NOT EXISTS liquidation_transactions (
	liquidation_id INT NOT NULL,
	transaction_id VARCHAR(20) NOT NULL,
	PRIMARY KEY (transaction_id),
	FOREIGN KEY (liquidation_id) REFERENCES liquidations,
	FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id)
);

