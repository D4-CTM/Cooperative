
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
RETURN capital * (interest / v_months);

CREATE OR REPLACE FUNCTION fn_calc_pmt(capital NUMERIC(8,2), interest NUMERIC(8,2), v_months INT)
LANGUAGE SQL
RETURNS NUMERIC(8,2)
DETERMINISTIC
RETURN (capital * interest) / (1 - (POW((1 + interest), -v_months)));

CREATE OR REPLACE PROCEDURE sp_insert_payment(IN v_iterations INT, IN v_loan_id CHAR(16), v_pmt NUMERIC(8,2))
LANGUAGE SQL
BEGIN
    DECLARE loan_amount NUMERIC(10,2);
    DECLARE loan_interest NUMERIC(5,2);
    DECLARE loan_periods INT;
    DECLARE loan_date DATE;
    DECLARE IPMT NUMERIC(8,2);    
    DECLARE PPMT NUMERIC(8,2);

	SELECT requested_amount, loan_interest, loan_periods, loan_date 
    INTO loan_amount, loan_interest, loan_periods, loan_date
    FROM loans 
    WHERE loan_id = v_loan_id;	

    SET IPMT = fn_calc_ipmt(loan_amount, loan_interest, loan_periods);
    SET PPMT = v_pmt - IPMT;    
    
    INSERT INTO payments(loan_id, payment_number, deadline, interest_to_pay, capital_to_pay) 
    VALUES(v_loan_id, LPAD(v_iterations, 5, '0'), ADD_MONTHS(loan_date, v_iterations), IPMT, PPMT);
END;

CREATE OR REPLACE TRIGGER tr_create_payments
AFTER INSERT ON loans
REFERENCING NEW AS n
FOR EACH ROW MODE DB2SQL
BEGIN
	DECLARE PMT NUMERIC(8,2);
	DECLARE i INT;
	
	SET PMT = fn_calc_pmt(requested_amount, loan_interest, loan_periods);
	SET i = 0;
	WHILE i < n.loan_periods
		DO CALL sp_insert_payment(i, n.loan_id, PMT);
		SET i = i + 1;
    END WHILE;
END;

