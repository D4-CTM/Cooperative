
DROP SEQUENCE user_seq;
DROP SEQUENCE loan_seq;

CREATE SEQUENCE user_seq;
CREATE SEQUENCE loan_seq;

CREATE OR REPLACE TRIGGER tr_user_id
BEFORE INSERT ON users
REFERENCING NEW AS n
FOR EACH ROW MODE DB2SQL
BEGIN
	DECLARE next_val BIGINT;  
  SET next_val = NEXT VALUE FOR user_seq;

  SET n.user_id = 'AF-' || LPAD(next_val, 5, '0');
END;

CREATE OR REPLACE TRIGGER tr_loan_id
BEFORE INSERT ON loans
REFERENCING NEW AS n
FOR EACH ROW MODE DB2SQL
BEGIN
  DECLARE next_val BIGINT;  
  SET next_val = NEXT VALUE FOR loan_seq;

  SET n.loan_id = n.user_id || '-PT'  || LPAD(next_val, 5, '0');
END;

CREATE OR REPLACE FUNCTION fn_calc_IPMT(capital NUMERIC(8,2), interest NUMERIC(8,2), v_months INT)
LANGUAGE SQL
RETURNS NUMERIC(8,2)
DETERMINISTIC
RETURN capital * CAST(interest / v_months AS DOUBLE);

CREATE OR REPLACE FUNCTION fn_calc_pmt(capital NUMERIC(8,2), interest NUMERIC(6,6), v_months INT)
LANGUAGE SQL
RETURNS NUMERIC(8,2)
DETERMINISTIC
RETURN (capital * interest) / CAST(1 - (POW((1 + interest), ((-1)*v_months))) AS NUMERIC(6,6));

CREATE OR REPLACE TRIGGER tr_create_payments
AFTER INSERT ON loans
REFERENCING NEW AS n
FOR EACH ROW
BEGIN ATOMIC
	DECLARE PMT NUMERIC(8,2);
	DECLARE capital NUMERIC(8,2);
	DECLARE interest_rate NUMERIC(6,6);
    DECLARE IPMT NUMERIC(8,2);    
    DECLARE PPMT NUMERIC(8,2);
	DECLARE i INT;
	
	SET capital = n.requested_amount;
	SET interest_rate = CAST(n.loan_interest/n.loan_periods AS NUMERIC(6,6));
	SET PMT = fn_calc_pmt(n.requested_amount, interest_rate, n.loan_periods);
	SET i = 0;
	WHILE i < loan_periods DO
	  SET IPMT = fn_calc_ipmt(capital, n.loan_interest, n.loan_periods);
	  SET PPMT = PMT - IPMT;

	  INSERT INTO payments(loan_id, payment_number, deadline, interest_to_pay, capital_to_pay)
	  VALUES(loan_id, LPAD((i+1), 5, '0'), ADD_MONTHS(n.loan_date, (i+1)), IPMT, PPMT);
	  SET capital = capital - PPMT;
    SET i = i + 1;
  END WHILE;
END;

CREATE OR REPLACE TRIGGER tr_get_creator_of_account
BEFORE INSERT ON accounts
REFERENCING NEW AS n
FOR EACH ROW MODE DB2SQL
BEGIN
  DECLARE user_first_name VARCHAR(50);
  DECLARE user_first_lastname VARCHAR(50);  
	
  SELECT u.first_name, u.first_lastname INTO user_first_name, user_first_lastname
  FROM users u WHERE u.user_id = n.user_id;
  
    SET n.created_by = user_first_name || ' ' || user_first_lastname;
    SET n.modified_by = n.created_by;
END;

