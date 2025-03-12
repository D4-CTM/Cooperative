
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
END

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
END

CREATE OR REPLACE PROCEDURE sp_delete_user(IN v_user_id CHAR(8))
LANGUAGE SQL
BEGIN
    DELETE FROM users WHERE user_id = v_user_id;
END;

