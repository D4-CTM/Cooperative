
CREATE OR REPLACE PROCEDURE sp_insert_user( 
	OUT v_user_id CHAR(8),
    IN v_password VARCHAR(25),
	IN v_first_name VARCHAR(50),
	IN v_second_name VARCHAR(50),
	IN v_first_lastname VARCHAR(50),
	IN v_second_lastname VARCHAR(50),
	IN v_address_house_number VARCHAR(50),
	IN v_address_street VARCHAR(50),
	IN v_address_avenue VARCHAR(50),
	IN v_address_city VARCHAR(50),
	IN v_address_department VARCHAR(50),
	IN v_address_reference VARCHAR(150),
	IN v_primary_email VARCHAR(100),
	IN v_secondary_email VARCHAR(100),
	IN v_birth_date DATE,
	IN v_hiring_date DATE,
	IN v_created_by VARCHAR(101),
	IN v_creation_date TIMESTAMP,
	IN v_modified_by VARCHAR(101),
	IN v_last_modification_date TIMESTAMP)
LANGUAGE SQL
BEGIN
    SELECT user_id
    INTO v_user_id FROM FINAL TABLE (
	INSERT INTO users (
        password, first_name, second_name, first_lastname, second_lastname, address_house_number, address_street,
        address_avenue, address_city, address_department, address_reference, primary_email,	secondary_email, birth_date,
        hiring_date, created_by,	creation_date, modified_by,	last_modification_date
    ) VALUES (
        v_password, v_first_name, v_second_name, v_first_lastname, v_second_lastname, v_address_house_number, v_address_street,
        v_address_avenue, v_address_city, v_address_department, v_address_reference, v_primary_email, v_secondary_email, v_birth_date,
        v_hiring_date, v_created_by, v_creation_date, v_modified_by, v_last_modification_date
	));

    INSERT INTO accounts
        (user_id, account_type)
    VALUES
        (v_user_id, 'CAP'),
        (v_user_id, 'CAR');
END;

CREATE OR REPLACE PROCEDURE sp_update_user( 
	v_user_id CHAR(8),
    v_password VARCHAR(25),
	v_first_name VARCHAR(50),
	v_second_name VARCHAR(50),
	v_first_lastname VARCHAR(50),
	v_second_lastname VARCHAR(50),
	v_address_house_number VARCHAR(50),
	v_address_street VARCHAR(50),
	v_address_avenue VARCHAR(50),
	v_address_city VARCHAR(50),
	v_address_department VARCHAR(50),
	v_address_reference VARCHAR(150),
	v_primary_email VARCHAR(100),
	v_secondary_email VARCHAR(100),
	v_birth_date DATE,
	v_hiring_date DATE,
	v_created_by VARCHAR(101),
	v_creation_date TIMESTAMP,
	v_modified_by VARCHAR(101),
	v_last_modification_date TIMESTAMP)
LANGUAGE SQL
BEGIN
    UPDATE users SET
        password = v_password, first_name = v_first_name, second_name = v_second_name, first_lastname = v_first_lastname, second_lastname = v_second_lastname,
        address_house_number = v_address_house_number, address_street = v_address_street, address_avenue = v_address_avenue, address_city = v_address_city,
        address_department = v_address_department, address_reference = v_address_reference, primary_email = v_primary_email, secondary_email = v_secondary_email, 
        birth_date = v_birth_date, hiring_date = v_hiring_date, created_by = v_created_by, creation_date = v_creation_date, modified_by = v_modified_by, 
        last_modification_date = v_last_modification_date
    WHERE user_id = v_user_id;
END;

CREATE OR REPLACE PROCEDURE sp_delete_user(IN v_user_id CHAR(8))
LANGUAGE SQL
BEGIN
    DELETE FROM users WHERE user_id = v_user_id;
END;

CREATE OR REPLACE PROCEDURE sp_grant_admin_to(v_user_id CHAR(8), v_admin BOOLEAN)
LANGUAGE SQL
BEGIN
	UPDATE users SET
	 is_admin = v_admin
	WHERE user_id = v_user_id;
END;

CREATE OR REPLACE PROCEDURE sp_insert_phone(
    IN v_user_id CHAR(8),
    IN v_user_phone_number INT,
    IN v_region_number INT
)
LANGUAGE SQL
BEGIN
    INSERT INTO phone_numbers(user_id, user_phone_number, region_number) VALUES(v_user_id, v_user_phone_number, v_region_number);
END;

CREATE OR REPLACE PROCEDURE sp_update_phone(
    IN v_user_id CHAR(8),
    IN v_user_phone_number INT,
    IN v_region_number INT
)
LANGUAGE SQL
BEGIN
    UPDATE phone_numbers SET 
    region_number = v_region_number 
    WHERE user_id = v_user_id AND user_phone_number = v_user_phone_number;
END;

CREATE OR REPLACE PROCEDURE sp_delete_phone(IN v_phone_number INT)
LANGUAGE SQL
BEGIN
    DELETE FROM phone_numbers WHERE user_phone_number = v_phone_number;
END;

CREATE OR REPLACE PROCEDURE sp_create_loan(
    OUT v_loan_id CHAR(16),
	OUT v_loan_date DATE,
	OUT v_is_payed BOOLEAN,
	IN v_user_id CHAR(8),
    IN v_loan_periods INT,
	IN v_loan_interest NUMERIC(3,3),
    IN v_requested_amount NUMERIC(8,2))
LANGUAGE SQL
BEGIN
	IF EXISTS (SELECT 1 FROM loans WHERE user_id = v_user_id AND IS_PAYED = FALSE) THEN
		SIGNAL SQLSTATE '45000'	SET MESSAGE_TEXT = 'Please finish paying the current loan';	
	END IF;
	
    SELECT loan_id, loan_date, is_payed
    INTO v_loan_id, v_loan_date, v_is_payed FROM FINAL TABLE ( 
    INSERT INTO loans
        (user_id, loan_periods, loan_interest, requested_amount)
    VALUES
        (v_user_id, v_loan_periods, v_loan_interest, v_requested_amount)
    );
END;

CREATE OR REPLACE PROCEDURE sp_update_loan(
	v_loan_id CHAR(16),
    v_loan_periods INT,
	v_loan_interest NUMERIC(3,3),
    v_requested_amount NUMERIC(8,2),
	v_loan_date DATE)
LANGUAGE SQL
BEGIN
    UPDATE loans SET
        loan_periods = v_loan_periods, loan_interest = v_loan_interest, 
        requested_amount = v_requested_amount, loan_date = v_loan_id
    WHERE loan_id = v_loan_id;
END;

CREATE OR REPLACE PROCEDURE sp_mark_loan_as_payed(v_loan_id CHAR(16))
LANGUAGE SQL
BEGIN
	UPDATE loans SET
		is_payed = TRUE
	WHERE loan_id = v_loan_id;
END;

CREATE OR REPLACE PROCEDURE sp_delete_loan(v_loan_id CHAR(16))
LANGUAGE SQL
BEGIN
    DELETE FROM loans WHERE loan_id = v_loan_id;
END;

