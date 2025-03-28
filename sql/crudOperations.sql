
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
	IN v_last_modification_date TIMESTAMP,
	IN v_admin BOOLEAN)
LANGUAGE SQL
BEGIN
	IF EXISTS (SELECT 1 FROM users WHERE primary_email = v_primary_email) THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'The primary email is already in use by another account!';
	END IF;
	
	SELECT user_id
    INTO v_user_id FROM FINAL TABLE (
	INSERT INTO users (
        password, first_name, second_name, first_lastname, second_lastname, address_house_number, address_street,
        address_avenue, address_city, address_department, address_reference, primary_email,	secondary_email, birth_date,
        hiring_date, created_by, creation_date, modified_by, last_modification_date, is_admin
    ) VALUES (
        v_password, v_first_name, v_second_name, v_first_lastname, v_second_lastname, v_address_house_number, v_address_street,
        v_address_avenue, v_address_city, v_address_department, v_address_reference, v_primary_email, v_secondary_email, v_birth_date,
        v_hiring_date, 'SYSTEM', v_creation_date, v_modified_by, v_last_modification_date, v_admin
	));

    INSERT INTO accounts
        (user_id, account_type)
    VALUES
        (v_user_id, 'CAP'),
        (v_user_id, 'CAR');

    INSERT INTO ACCOUNT_PROFIT(ACCOUNT_ID, PROFIT)
	VALUES(v_user_id || '-CAR', 0);
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
	v_modified_by VARCHAR(101),
	v_last_modification_date TIMESTAMP,
	v_admin BOOLEAN)
LANGUAGE SQL
BEGIN
	IF EXISTS (SELECT 1 FROM users WHERE primary_email = v_primary_email AND user_id != v_user_id) THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'The primary email is already in use by another account!';
	END IF;

    UPDATE users SET
        password = v_password, first_name = v_first_name, second_name = v_second_name, first_lastname = v_first_lastname, second_lastname = v_second_lastname,
        address_house_number = v_address_house_number, address_street = v_address_street, address_avenue = v_address_avenue, address_city = v_address_city,
        address_department = v_address_department, address_reference = v_address_reference, primary_email = v_primary_email, secondary_email = v_secondary_email, 
        birth_date = v_birth_date, modified_by = v_modified_by, last_modification_date = v_last_modification_date, is_admin = v_admin
    WHERE user_id = v_user_id;
END;

CREATE OR REPLACE PROCEDURE sp_delete_user(IN v_user_id CHAR(8))
LANGUAGE SQL
BEGIN
    DELETE FROM users WHERE user_id = v_user_id;
END;

CREATE OR REPLACE PROCEDURE sp_insert_phone(
    IN v_user_id CHAR(8),
    IN v_user_phone_number INT,
    IN v_region_number INT
)
LANGUAGE SQL
BEGIN	
	IF NOT EXISTS (SELECT 1 FROM phone_numbers WHERE user_phone_number = v_user_phone_number) THEN
	    INSERT INTO phone_numbers(user_id, user_phone_number, region_number) VALUES(v_user_id, v_user_phone_number, v_region_number);
	END IF;
END;

CREATE OR REPLACE PROCEDURE sp_delete_phone(IN v_phone_number INT)
LANGUAGE SQL
BEGIN
	IF EXISTS (SELECT 1 FROM phone_numbers WHERE user_phone_number = v_phone_number) THEN
	    DELETE FROM phone_numbers WHERE user_phone_number = v_phone_number;
	END IF;
END;

CREATE OR REPLACE PROCEDURE sp_create_loan(
    OUT v_loan_id CHAR(16),
	OUT v_loan_date DATE,
	OUT v_is_payed BOOLEAN,
	IN v_user_id CHAR(8),
    IN v_loan_periods INT,
	IN v_loan_interest NUMERIC(3,3),
    IN v_requested_amount NUMERIC(18, 2))
LANGUAGE SQL
BEGIN
	DECLARE lv_created_by VARCHAR(101);
	
	IF EXISTS (SELECT 1 FROM loans WHERE user_id = v_user_id AND IS_PAYED = FALSE) THEN
		SIGNAL SQLSTATE '45000'	SET MESSAGE_TEXT = 'Please finish paying the current loan';	
	END IF;
	
	SELECT fn_get_creator_of(v_user_id || '-CAP') 
	INTO lv_created_by
	FROM SYSIBM.SYSDUMMY1;

    SELECT loan_id, loan_date, is_payed
    INTO v_loan_id, v_loan_date, v_is_payed FROM FINAL TABLE ( 
    INSERT INTO loans
        (user_id, loan_periods, loan_interest, requested_amount, created_by, creation_date, modified_by, last_modification)
    VALUES
        (v_user_id, v_loan_periods, v_loan_interest, v_requested_amount, lv_created_by, CURRENT_TIMESTAMP, lv_created_by, CURRENT_TIMESTAMP)
    );
END;

CREATE OR REPLACE PROCEDURE sp_update_loan(
	v_loan_id CHAR(16),
    v_loan_periods INT,
	v_loan_interest NUMERIC(3,3),
    v_requested_amount NUMERIC(18, 2),
	v_loan_date DATE)
LANGUAGE SQL
BEGIN
    DECLARE lv_modified_by VARCHAR(101);
	DECLARE lv_user_id CHAR(8);

	SELECT user_id
	INTO lv_user_id
	FROM LOANS l
	WHERE l.loan_id = v_loan_id;

	SELECT fn_get_creator_of(lv_user_id || '-CAP') 
	INTO lv_modified_by
	FROM SYSIBM.SYSDUMMY1;

	
	UPDATE loans SET
        loan_periods = v_loan_periods, loan_interest = v_loan_interest, 
        requested_amount = v_requested_amount, loan_date = v_loan_id,
        modified_by = lv_modified_by, last_modification = CURRENT_TIMESTAMP
    WHERE loan_id = v_loan_id;
END;

CREATE OR REPLACE PROCEDURE sp_delete_loan(v_loan_id CHAR(16))
LANGUAGE SQL
BEGIN
    DELETE FROM loans WHERE loan_id = v_loan_id;
END;

CREATE OR REPLACE FUNCTION fn_calculate_remaining_payment(IN v_payment_id INT)
RETURNS NUMERIC(18, 2)
LANGUAGE SQL
NOT DETERMINISTIC
BEGIN
	DECLARE amount_to_pay NUMERIC(18, 2);
	
	SELECT 
		(MAX(p.AMOUNT_TO_PAY) - COALESCE(SUM(t.TRANSACTION_AMMOUNT), 0)) AS AMOUNT_PAYED
	INTO amount_to_pay
	FROM payments p
	LEFT JOIN PAYMENT_TRANSACTIONS  pt
		ON pt.PAYMENT_ID  = p.PAYMENT_ID	
	LEFT JOIN transactions t
		ON pt.TRANSACTION_ID = t.TRANSACTION_ID
	where v_payment_id = p.PAYMENT_ID;
	
	RETURN amount_to_pay;
END;

CREATE OR REPLACE PROCEDURE sp_payment_transaction(
	OUT v_transaction_id VARCHAR(20),
	IN v_account_id CHAR(12),
	IN v_transaction_date DATE,
	IN v_transaction_ammount NUMERIC(18, 2),
	IN v_transaction_comment VARCHAR(255),
	IN v_payment_id INT,
	IN v_loan_id CHAR(16))
LANGUAGE SQL
BEGIN ATOMIC
	DECLARE lv_error_msg VARCHAR(255);
	DECLARE lv_amount_to_pay NUMERIC(18, 2);
	DECLARE lv_acc_owner VARCHAR(101);
	DECLARE lv_remaining_payments INT;
	DECLARE lv_loan_periods INT;

	SELECT fn_calculate_remaining_payment(v_payment_id)
	INTO lv_amount_to_pay
	FROM SYSIBM.SYSDUMMY1;

	IF lv_amount_to_pay IS NULL THEN
		SIGNAL SQLSTATE '45000'	SET MESSAGE_TEXT = 'Could not calculate the amount needed for this payment, please check the sended payment id!';
	END IF;
	
	IF lv_amount_to_pay < v_transaction_ammount THEN
		SET lv_error_msg = v_payment_id || ' - ' || v_loan_id|| ': Amount to pay: ' || lv_amount_to_pay || ' - transaction amount: ' || v_transaction_ammount;
		SIGNAL SQLSTATE '45000'	SET MESSAGE_TEXT = lv_error_msg;
	END IF;
	
	IF (SELECT RIGHT(v_account_id, 3) FROM SYSIBM.SYSDUMMY1) != 'CAP' THEN
		SIGNAL SQLSTATE '45000'	SET MESSAGE_TEXT = 'The only account capable of making payment is the apportation account!!';
	END IF;

	
	CALL sp_retire_money(v_transaction_id,
	v_account_id,
	v_transaction_date,
	v_transaction_ammount,
	v_transaction_comment);

	INSERT INTO payment_transactions(payment_id, transaction_id)
	VALUES (v_payment_id, v_transaction_id);
	
	SELECT COUNT(P.IS_PAYED)
	INTO lv_remaining_payments
	FROM PAYMENTS P
	WHERE P.LOAN_ID = v_loan_id AND P.IS_PAYED = FALSE;
	
	IF (lv_amount_to_pay - v_transaction_ammount) = 0 THEN
		UPDATE payments p
		SET p.IS_PAYED = TRUE
		WHERE p.payment_id = v_payment_id;
	
		SET lv_remaining_payments = lv_remaining_payments - 1;
	END IF;
		
	IF lv_remaining_payments = 0 THEN 
		UPDATE loans SET
			is_payed = TRUE
		WHERE loan_id = v_loan_id;
	END IF;
	
    SELECT fn_get_creator_of(v_account_id) 
	INTO lv_acc_owner
    FROM SYSIBM.SYSDUMMY1;

	
	UPDATE loans
	SET last_modification = CURRENT_TIMESTAMP,
	modified_by = lv_acc_owner
	WHERE LOAN_ID = v_loan_id;
END;

CREATE OR REPLACE PROCEDURE sp_retire_money(
	INOUT v_transaction_id VARCHAR(20),
	IN v_account_id CHAR(12),
	IN v_transaction_date DATE,
	IN v_transaction_ammount NUMERIC(18, 2),
	IN v_transaction_comment VARCHAR(255)
)
LANGUAGE SQL
BEGIN ATOMIC
	DECLARE user_cant INT;
	DECLARE lv_acc_owner VARCHAR(101);

	IF v_transaction_ammount <= 0 THEN
		SIGNAL SQLSTATE '45000'	SET MESSAGE_TEXT = 'The retire amount must be greater than 0!';	
	END IF;
	
  	SELECT COALESCE(COUNT(t.account_id) + 1, 1) 
	INTO user_cant FROM transactions t
  	WHERE t.account_id = v_account_id;

  	SET v_transaction_id = v_account_id || '-' || CAST(user_cant AS VARCHAR(5));

    INSERT INTO transactions(transaction_id,  account_id, transaction_date, transaction_ammount, transaction_comment)
    VALUES(v_transaction_id, v_account_id, v_transaction_date, v_transaction_ammount, v_transaction_comment);
	
	IF (SELECT balance FROM accounts WHERE account_id = v_account_id) < v_transaction_ammount THEN 
		SIGNAL SQLSTATE '45000'	SET MESSAGE_TEXT = 'Not enough money on the account balance to retire!';	
	END IF;
	
    SELECT fn_get_creator_of(v_account_id) 
	INTO lv_acc_owner
    FROM SYSIBM.SYSDUMMY1;

	
	UPDATE accounts
	SET balance = balance - v_transaction_ammount,
	LAST_MODIFICATION_DATE = CURRENT_TIMESTAMP,
	MODIFIED_BY = lv_acc_owner
	WHERE account_id  = v_account_id;
END;

-- deposit money, I forgot TO CHANGE th name
CREATE OR REPLACE PROCEDURE sp_create_transaction(    
	OUT v_transaction_id VARCHAR(20),
	IN v_account_id CHAR(12),
	IN v_transaction_date DATE,
	IN v_transaction_ammount NUMERIC(18, 2),
	IN v_transaction_comment VARCHAR(255)
)
LANGUAGE SQL
BEGIN ATOMIC
	DECLARE user_cant INT;
	DECLARE lv_acc_owner VARCHAR(101);

	IF v_transaction_ammount <= 0 THEN
		SIGNAL SQLSTATE '45000'	SET MESSAGE_TEXT = 'The deposit amount must be greater than 0!';	
	END IF;
	
  	SELECT COALESCE(COUNT(t.account_id) + 1, 1) 
	INTO user_cant FROM transactions t
  	WHERE t.account_id = v_account_id;

  	SET v_transaction_id = v_account_id || '-' || CAST(user_cant AS VARCHAR(5));

    INSERT INTO transactions(transaction_id, account_id, transaction_date, transaction_ammount, transaction_comment)
    VALUES(v_transaction_id, v_account_id, v_transaction_date, v_transaction_ammount, v_transaction_comment);
	
    SELECT fn_get_creator_of(v_account_id) 
	INTO lv_acc_owner
    FROM SYSIBM.SYSDUMMY1;

    
	UPDATE accounts
	SET balance = balance + v_transaction_ammount,
	LAST_MODIFICATION_DATE = CURRENT_TIMESTAMP,
	MODIFIED_BY = lv_acc_owner
	WHERE account_id  = v_account_id;
END;

CREATE OR REPLACE PROCEDURE sp_change_transaction_comment(
    IN v_transaction_id VARCHAR(20),
	IN v_transaction_comment VARCHAR(255)
)
LANGUAGE SQL
BEGIN
    UPDATE transactions 
    SET transaction_comment = v_transaction_comment
    WHERE transaction_id = v_transaction_id;
END;

CREATE OR REPLACE PROCEDURE sp_make_remaining_transactions(
	IN v_month INT,
	IN v_year INT,
	IN v_closure_id INT
)
LANGUAGE SQL
BEGIN ATOMIC
	DECLARE lv_account_id CHAR(12);
	DECLARE lv_deposited_amount NUMERIC(18,2);
	DECLARE lv_transaction_id VARCHAR(20);

	DECLARE crsr_transactions_finished INT DEFAULT 0;
	DECLARE lv_accounts_cursor CURSOR
	FOR SELECT account_id, deposited_amount
	FROM (
		SELECT a.ACCOUNT_ID, COALESCE(SUM(t.TRANSACTION_AMMOUNT), 0) AS deposited_amount
		FROM accounts a
		LEFT JOIN TRANSACTIONS t
		ON a.ACCOUNT_ID  = t.ACCOUNT_ID
		JOIN USERS u 
		ON u.USER_ID = a.USER_ID
		WHERE 
			((MONTH(t.TRANSACTION_DATE) = v_month AND 
			YEAR(t.TRANSACTION_DATE) = v_year) OR
			t.ACCOUNT_ID IS NULL) AND 
			u.IS_ACTIVE = TRUE
		GROUP BY a.ACCOUNT_ID)
	WHERE DEPOSITED_AMOUNT < 200;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET crsr_transactions_finished = 1;
	
	OPEN lv_accounts_cursor;

    FETCH lv_accounts_cursor INTO lv_account_id, lv_deposited_amount;
    
    WHILE crsr_transactions_finished = 0 DO   			
		CALL sp_create_transaction(lv_transaction_id, lv_account_id, CURRENT_DATE, (200 - lv_deposited_amount), 'Transaction made during the monthly closure!');
    
    	UPDATE accounts
    	SET LAST_MODIFICATION_DATE = CURRENT_TIMESTAMP,
		MODIFIED_BY = 'System'
    	WHERE account_id = lv_account_id;
		
		INSERT INTO CLOSURE_TRANSACTIONS(CLOSURE_ID, TRANSACTION_ID)
		VALUES(v_closure_id, lv_transaction_id);
		
	    FETCH lv_accounts_cursor INTO lv_account_id, lv_deposited_amount;
    END WHILE;
	
	CLOSE lv_accounts_cursor;
END;

CREATE OR REPLACE PROCEDURE sp_make_monthly_payment(
	IN v_month INT,
	IN v_year INT,
	IN v_closure_id INT	
)
LANGUAGE SQL
BEGIN ATOMIC
	DECLARE lv_account_id CHAR(12);
	DECLARE lv_amount_to_pay NUMERIC(18,2);
	DECLARE lv_account_balance NUMERIC(18,2);
	DECLARE lv_transaction_id VARCHAR(20);
	DECLARE lv_payment_id INT;
	DECLARE lv_loan_id CHAR(16);
	DECLARE lv_user_id CHAR(8);
	DECLARE crsr_accounts_finished INT DEFAULT 0;

	DECLARE lv_accounts_cursor CURSOR
	FOR SELECT p.PAYMENT_ID, P.AMOUNT_TO_PAY, l.LOAN_ID, u.USER_ID
	FROM PAYMENTS p 
	JOIN LOANS l
	ON l.LOAN_ID = p.LOAN_ID
	JOIN USERS u
	ON u.USER_ID = l.USER_ID
	WHERE 
		MONTH(p.DEADLINE) = v_month AND
		YEAR(p.DEADLINE) = v_year AND 
		p.IS_PAYED = FALSE AND
		u.IS_ACTIVE = TRUE;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET crsr_accounts_finished = 1;

	OPEN lv_accounts_cursor;

    FETCH lv_accounts_cursor INTO lv_payment_id, lv_amount_to_pay, lv_loan_id, lv_user_id;
        
    pay_loop: WHILE crsr_accounts_finished = 0 DO
	    SELECT a.balance, a.account_id
	    INTO lv_account_balance, lv_account_id
		FROM ACCOUNTS a
		WHERE a.ACCOUNT_TYPE = 'CAP' AND
			lv_user_id = a.USER_ID;
    
    	IF lv_account_balance = 0 THEN
		    FETCH lv_accounts_cursor INTO lv_payment_id, lv_amount_to_pay, lv_loan_id, lv_user_id;
    		ITERATE pay_loop;
    	END IF;
    
		IF lv_account_balance > lv_amount_to_pay THEN
			SET lv_account_balance = lv_amount_to_pay;
		END IF;
			
    	CALL sp_payment_transaction(
		lv_transaction_id,
		lv_account_id,
		CURRENT_DATE,
		lv_account_balance,
		'Transaction made during the monthly closure to cover this month loan payment',
		lv_payment_id,
		lv_loan_id);
    	
    	UPDATE accounts
    	SET LAST_MODIFICATION_DATE = CURRENT_TIMESTAMP,
		MODIFIED_BY = 'System'
    	WHERE account_id = lv_account_id;

    	UPDATE loans
    	SET LAST_MODIFICATION = CURRENT_TIMESTAMP,
		MODIFIED_BY = 'System'
    	WHERE loan_id = lv_loan_id;
    	
		INSERT INTO CLOSURE_TRANSACTIONS(CLOSURE_ID, TRANSACTION_ID)
		VALUES(v_closure_id, lv_transaction_id);
		
		INSERT INTO CLOSURE_PAYMENTS(PAYMENT_ID, CLOSURE_ID)
		VALUES(lv_payment_id, v_closure_id);
		
	    FETCH lv_accounts_cursor INTO lv_payment_id, lv_amount_to_pay, lv_loan_id, lv_user_id;
    END WHILE;
	
	CLOSE lv_accounts_cursor;
END;

CREATE OR REPLACE PROCEDURE sp_generate_payouts(IN v_closure_id INT)
LANGUAGE SQL
BEGIN ATOMIC
	DECLARE lv_apportation INT;
	DECLARE lv_account_profit NUMERIC(18, 2);
	DECLARE lv_account_balance NUMERIC(18, 2);
	DECLARE lv_account_id CHAR(12);
	DECLARE lv_total_profit NUMERIC(18, 2);
	DECLARE lv_total_balance NUMERIC(18,2);
	DECLARE lv_user_id CHAR(8);	
	DECLARE TEST VARCHAR(255);

	DECLARE crsr_accounts_apportation_finished INT DEFAULT 0;

	DECLARE lv_accounts_cursor CURSOR
	FOR SELECT A.BALANCE, A.ACCOUNT_ID
	FROM USERS U
	JOIN ACCOUNTS A
	ON A.USER_ID = U.USER_ID
	WHERE U.IS_ACTIVE = TRUE AND 
		  A.ACCOUNT_TYPE = 'CAP' AND
		  A.BALANCE > 0;

	DECLARE CONTINUE HANDLER FOR NOT FOUND SET crsr_accounts_apportation_finished = 1;

	SELECT SUM(A.BALANCE)
	INTO lv_total_balance
	FROM USERS U
	JOIN ACCOUNTS A
	ON A.USER_ID = U.USER_ID
	WHERE U.IS_ACTIVE = TRUE AND A.ACCOUNT_TYPE = 'CAP';

	SELECT SUM(p.INTEREST_TO_PAY)
	INTO lv_total_profit
	FROM CLOSURES c
	JOIN PAYMENTS p 
	ON MONTH(p.DEADLINE) = c.CLOSURE_MONTH AND YEAR(p.DEADLINE) = c.CLOSURE_YEAR
	WHERE p.IS_PAYED = TRUE;

	IF lv_total_profit <= 0 THEN 
		RETURN ;
	END IF;
	
	OPEN lv_accounts_cursor;

    FETCH lv_accounts_cursor INTO lv_account_balance, lv_account_id;
    
    WHILE crsr_accounts_apportation_finished = 0 DO
		SET lv_apportation = (lv_account_balance/lv_total_balance) * 10000.0000;
		SET lv_account_profit = lv_total_profit*CAST(lv_apportation/10000.0000 AS NUMERIC(8, 4));
       
	    INSERT INTO PAYOUTS(CLOSURE_ID, ACCOUNT_ID, ACCOUNT_BALANCE, APPORTATION_PERCENTAGE, ACCOUNT_PROFIT )
	    VALUES(v_closure_id, lv_account_id, lv_account_balance, lv_apportation, lv_account_profit);

    	UPDATE accounts
    	SET LAST_MODIFICATION_DATE = CURRENT_TIMESTAMP,
		MODIFIED_BY = 'System'
    	WHERE account_id = lv_account_id;
	    
	    SELECT USER_ID
	    INTO lv_user_id
	    FROM ACCOUNTS a 
	    WHERE a.ACCOUNT_ID = lv_account_id;
	    
	    UPDATE ACCOUNT_PROFIT ap
	    SET ap.PROFIT = ap.PROFIT + lv_account_profit
		WHERE ap.ACCOUNT_ID = (lv_user_id || '-CAR');
	   	
	    FETCH lv_accounts_cursor INTO lv_account_balance, lv_account_id;    
	END WHILE;
	
	CLOSE lv_accounts_cursor;
END;

CREATE OR REPLACE PROCEDURE sp_generate_monthly_closure(
	IN v_month INT,
	IN v_year INT,
	IN description VARCHAR(255)
)
LANGUAGE SQL
BEGIN ATOMIC
	DECLARE lv_closure_id INT;
	DECLARE ERROR_MSG VARCHAR(255);

	IF NOT EXISTS (SELECT 1 FROM closures WHERE closure_month = v_month AND closure_year = v_year) THEN
		SELECT closure_id
		INTO lv_closure_id
		FROM FINAL TABLE (
			INSERT INTO closures(closure_month, closure_year, description)
			VALUES(v_month, v_year, description));
	ELSE
		SET ERROR_MSG = 'Already made the closure for this month! (' || v_month || '/' || v_year || ')';
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = ERROR_MSG;
	END IF;

	CALL sp_make_remaining_transactions(v_month, v_year, lv_closure_id);
	CALL sp_make_monthly_payment(v_month, v_year, lv_closure_id);
	CALL sp_generate_payouts(lv_closure_id);
END;

CREATE OR REPLACE PROCEDURE sp_change_closure_comment(
	IN v_month INT,
	IN v_year INT,
	IN v_neo_description VARCHAR(255)
)
LANGUAGE SQL
BEGIN
	UPDATE closures
	SET description = v_neo_description
	WHERE closure_month = v_month AND closure_year = v_year;
END;

CREATE OR REPLACE PROCEDURE sp_retire_profit(
	OUT v_transaction_id VARCHAR(20),
	IN v_account_id CHAR(12),
	IN v_transaction_date DATE,
	IN v_transaction_ammount NUMERIC(18, 2),
	IN v_transaction_comment VARCHAR(255)
)
LANGUAGE SQL
BEGIN ATOMIC
	DECLARE user_cant INT;
	DECLARE lv_balance DECIMAL(18, 2);
	DECLARE lv_profit DECIMAL(18, 2);
	DECLARE lv_acc_owner VARCHAR(101);

    SELECT BALANCE, PROFIT
    INTO lv_balance, lv_profit
	FROM ACCOUNTS a
	JOIN ACCOUNT_PROFIT ap
	ON a.ACCOUNT_ID = ap.ACCOUNT_ID
	WHERE a.ACCOUNT_ID = v_account_id;

	IF v_transaction_ammount <= 0 THEN
		SIGNAL SQLSTATE '45000'	SET MESSAGE_TEXT = 'The transaction amount must be greater than 0!';	
	END IF;

  	SELECT COALESCE(COUNT(t.account_id) + 1, 1) 
	INTO user_cant FROM transactions t
  	WHERE t.account_id = v_account_id;

  	SET v_transaction_id = v_account_id || '-' || CAST(user_cant AS VARCHAR(5));

    INSERT INTO transactions(transaction_id, account_id, transaction_date, transaction_ammount, transaction_comment)
    VALUES(v_transaction_id, v_account_id, v_transaction_date, v_transaction_ammount, v_transaction_comment);
        
    IF v_transaction_ammount > (lv_balance + lv_profit) THEN 
		SIGNAL SQLSTATE '45000'	SET MESSAGE_TEXT = 'Not enough money register into your account to make this transaction!';
    END IF;

    -- Just an edge case, but tecnically, this SHOULD never trigger!
    IF (v_transaction_ammount - lv_balance) > lv_profit THEN
		SIGNAL SQLSTATE '45000'	SET MESSAGE_TEXT = 'Not enough money gain with dividends to make this transaction!';
    END IF;
    
    SELECT fn_get_creator_of(v_account_id) 
	INTO lv_acc_owner
    FROM SYSIBM.SYSDUMMY1;

    UPDATE ACCOUNTS
    SET BALANCE = 0,
    LAST_MODIFICATION_DATE = CURRENT_TIMESTAMP,
	MODIFIED_BY = lv_acc_owner
    WHERE ACCOUNT_ID = v_account_id;

	UPDATE account_profit
	SET profit = profit - (v_transaction_ammount - lv_balance)
	WHERE account_id  = v_account_id;
END;

CREATE OR REPLACE PROCEDURE sp_create_partial_liquidation(
	OUT v_liquidation_id INT,
	IN v_account_id CHAR(12),
	IN v_total_money DECIMAL(8, 2), --retired money
	IN v_liquidation_date DATE,
	IN v_transaction_comment VARCHAR(255)
)
LANGUAGE SQL
BEGIN ATOMIC
	DECLARE lv_loan_id VARCHAR(16);
	DECLARE lv_transaction_id VARCHAR(20);
	DECLARE lv_year INT;
	DECLARE lv_month INT;

	SELECT EXTRACT(MONTH FROM v_liquidation_Date), EXTRACT(YEAR FROM v_liquidation_Date)
	INTO lv_month, lv_year
	FROM SYSIBM.SYSDUMMY1;
	
	IF (lv_month != 6 AND lv_month != 12) THEN 
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'This option is valid only during the months of June and December!';
	END IF;
	
	IF NOT EXISTS (SELECT 1 FROM CLOSURES WHERE CLOSURE_MONTH = lv_month AND CLOSURE_YEAR = lv_year) THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'To make a partial liquidation first is needed to make the monthly closure!';
	END IF;
	
	IF v_total_money <= (SELECT BALANCE FROM ACCOUNTS WHERE account_id = v_account_id) THEN
		CALL sp_retire_money(lv_transaction_id, v_account_id, v_liquidation_date, v_total_money, v_transaction_comment);
	ELSE
		CALL sp_retire_profit(lv_transaction_id, v_account_id, v_liquidation_date, v_total_money, v_transaction_comment);		
	END IF;

	SELECT liquidation_id
	INTO v_liquidation_id FROM FINAL TABLE (
		INSERT INTO liquidations(account_id, liquidation_type, retirement_date, total_money)
		VALUES(v_account_id, 'P', v_liquidation_date, v_total_money)
	);
	
	INSERT INTO LIQUIDATION_TRANSACTIONS(liquidation_id, transaction_id)
	VALUES(v_liquidation_id, lv_transaction_id);
END;

CREATE OR REPLACE PROCEDURE sp_create_total_liquidation(
	OUT v_liquidation_id INT,
	IN v_account_id CHAR(12),
	IN v_liquidation_date DATE,
	IN v_transaction_comment VARCHAR(255)
)
LANGUAGE SQL
BEGIN ATOMIC
	DECLARE lv_loan_id VARCHAR(16);
	DECLARE lv_remaining_loan NUMERIC(18,2);
	DECLARE lv_transaction_id VARCHAR(20);
	DECLARE lv_cap_balance NUMERIC(18,2);
	DECLARE lv_car_balance NUMERIC(18,2);
	DECLARE lv_profit NUMERIC(18,2);
	DECLARE lv_month INT;

	SELECT
		MAX(CASE WHEN a.ACCOUNT_TYPE = 'CAP' THEN a.BALANCE END),
		MAX(CASE WHEN a.ACCOUNT_TYPE = 'CAR' THEN a.BALANCE END),
		MAX(ap.PROFIT)
	INTO lv_cap_balance, lv_car_balance, lv_profit
	FROM USERS u 
	JOIN ACCOUNTS a 
	ON u.USER_ID = a.USER_ID 
	LEFT JOIN ACCOUNT_PROFIT ap 
	ON ap.ACCOUNT_ID = a.ACCOUNT_ID
	WHERE a.USER_ID = LEFT(v_account_id, 8);
	
	SELECT l.LOAN_ID
	INTO lv_loan_id
	FROM LOANS l 
	WHERE l.USER_ID = LEFT(v_account_id, 8)
	AND l.IS_PAYED = FALSE;

	SET lv_car_balance = lv_car_balance + lv_profit;
	
	-- If the sum of all is greater than 0, then it, does in fact, has money
	IF (lv_cap_balance + lv_car_balance) > 0 THEN
		SELECT liquidation_id
		INTO v_liquidation_id FROM FINAL TABLE (
			INSERT INTO liquidations(account_id, liquidation_type, retirement_date, total_money)
			VALUES(v_account_id, 'T', v_liquidation_date, lv_car_balance + lv_cap_balance)
		);
	END IF;
	
	-- If it has money on the savings account, we take it and deposit it into the apportation accounts
	IF lv_car_balance > 0 THEN
		CALL sp_retire_profit(
		lv_transaction_id,
		LEFT(v_account_id, 8) || '-CAR',
		v_liquidation_date,
		lv_car_balance,
		'Transaction made during total liquidation!');
				
		INSERT INTO LIQUIDATION_TRANSACTIONS(liquidation_id, transaction_id)
		VALUES(v_liquidation_id, lv_transaction_id);
			
		CALL sp_create_transaction(lv_transaction_id,
		LEFT(v_account_id, 8) || '-CAP',
		v_liquidation_date,
		lv_car_balance,
		'Transfered from savings accounts to pay the loan!');
	
		INSERT INTO LIQUIDATION_TRANSACTIONS(liquidation_id, transaction_id)
		VALUES(v_liquidation_id, lv_transaction_id);
	END IF;
	
	SET lv_cap_balance = lv_cap_balance + lv_car_balance;
	
	IF lv_loan_id IS NULL AND lv_cap_balance <= 0 THEN
		UPDATE USERS 
		SET IS_ACTIVE = FALSE
		WHERE USER_ID = LEFT(v_account_id, 8);	

		RETURN;
	END IF;
	
	-- We pay our loans with our money, it uses the apportations acc as it's the only one who can actually pay loans!
	IF lv_loan_id IS NOT NULL THEN
		SELECT SUM(fn_calculate_remaining_payment(payment_id)) 
		INTO lv_remaining_loan
		FROM PAYMENTS p 
		WHERE p.LOAN_ID = lv_loan_id;
		
		IF lv_remaining_loan > lv_cap_balance THEN
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Not enough money in account to finish paying the loan!';
		END IF;

		CALL sp_pay_loan(lv_cap_balance,
		lv_remaining_loan,
		lv_loan_id,
		v_account_id,
		v_liquidation_date,
		v_liquidation_id);
	END IF;	

	-- We check again to see if we still have money left, so we can retire it
	IF lv_cap_balance > 0 THEN	
		CALL sp_retire_money(
		lv_transaction_id, 
		LEFT(v_account_id, 8) || '-CAP', 
		v_liquidation_date, 
		lv_cap_balance, 
		v_transaction_comment);
		
		INSERT INTO LIQUIDATION_TRANSACTIONS(liquidation_id, transaction_id)
		VALUES(v_liquidation_id, lv_transaction_id);
	END IF;
		
	UPDATE USERS 
	SET IS_ACTIVE = FALSE
	WHERE USER_ID = LEFT(v_account_id, 8);	
END;

CREATE OR REPLACE PROCEDURE sp_pay_loan(
	INOUT v_cap_balance NUMERIC(18,2),
	IN lv_remaining_loan NUMERIC(18,2),
	IN v_loan_id VARCHAR(16),
	IN v_account_id CHAR(12),
	IN v_liquidation_date DATE,
	IN v_liquidation_id INT
)
LANGUAGE SQL
BEGIN ATOMIC
	DECLARE lv_remaining_payment NUMERIC(18,2);
	DECLARE lv_transaction_id VARCHAR(20);
	DECLARE lv_payment_id INT;

	DECLARE crsr_payments_finish INT DEFAULT 0;
	DECLARE crsr_payments CURSOR
	FOR SELECT p.PAYMENT_ID, fn_calculate_remaining_payment(payment_id)
		FROM PAYMENTS p 
		WHERE p.LOAN_ID = v_loan_id AND p.IS_PAYED = FALSE;
	
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET crsr_payments_finish = 1;
		
	OPEN crsr_payments;
	
	FETCH crsr_payments INTO lv_payment_id, lv_remaining_payment;
	
	WHILE crsr_payments_finish = 0 DO  
	
		CALL sp_payment_transaction(lv_transaction_id, 
	    LEFT(v_account_id, 8) || '-CAP',
	    CURRENT_DATE,
	    lv_remaining_payment,
	    'Payment made during total liquidation!',
	    lv_payment_id,
	    v_loan_id);	

		INSERT INTO LIQUIDATION_TRANSACTIONS (liquidation_id, transaction_id)
		VALUES(v_liquidation_id, lv_transaction_id);
		
		INSERT INTO LIQUIDATION_PAYMENTS (liquidation_id, payment_id)
		VALUES(v_liquidation_id, lv_payment_id);
	
		FETCH crsr_payments INTO lv_payment_id, lv_remaining_payment;
	END WHILE;
	
	CLOSE crsr_payments;
	
	SET v_cap_balance = lv_remaining_loan - v_cap_balance;
END;

