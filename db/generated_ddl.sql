create database loan_management

create table loan_schema.loans
(
    loan_id             serial
        primary key,
    principal_amount    numeric(15, 2)                          not null
        constraint chk_principal_amount_positive
            check (principal_amount > (0)::numeric),
    interest_rate       numeric(5, 2)            default 10 not null
        constraint chk_interest_rate_positive
            check (interest_rate >= (0)::numeric),
    term_weeks          integer                                 not null
        constraint chk_term_weeks_positive
            check (term_weeks > 0),
    weekly_payment      numeric(15, 2)                          not null
        constraint chk_weekly_payment_positive
            check (weekly_payment > (0)::numeric),
    outstanding_balance numeric(15, 2)                          not null
        constraint chk_outstanding_balance_non_negative
            check (outstanding_balance >= (0)::numeric),
    delinquent          boolean                  default false,
    late_fee_percentage numeric(5, 2)            default 0.00,
    max_late_fee        numeric(15, 2),
    grace_period_days   integer                  default 3,
    created_at          timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at          timestamp with time zone default CURRENT_TIMESTAMP
);


create table loan_schema.paymentschedules
(
    schedule_id       serial
        primary key,
    loan_id           integer        not null
        references loan_schema.loans,
    week_number       integer        not null
        constraint chk_week_number_positive
            check (week_number > 0),
    due_amount        numeric(15, 2) not null
        constraint chk_due_amount_positive
            check (due_amount > (0)::numeric),
    due_date          date           not null,
    paid              boolean                  default false,
    payment_date      date,
    late_fee_applied  boolean                  default false,
    late_fee_amount   numeric(15, 2)           default 0.00
        constraint chk_late_fee_amount_non_negative
            check (late_fee_amount >= (0)::numeric),
    grace_period_days integer                  default 3
        constraint chk_grace_period_days_positive
            check (grace_period_days >= 0),
    grace_period_end  date generated always as ((due_date + ((grace_period_days)::double precision * '1 day'::interval))) stored,
    created_at        timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at        timestamp with time zone default CURRENT_TIMESTAMP
);

create table loan_schema.payments
(
    payment_id     serial
        primary key,
    loan_id        integer                                                       not null
        references loan_schema.loans,
    payment_amount numeric(15, 2)                                                not null
        constraint chk_payment_amount_positive
            check (payment_amount > (0)::numeric),
    payment_date   timestamp with time zone default CURRENT_TIMESTAMP,
    week_number    integer                                                       not null
        constraint chk_week_number_positive
            check (week_number > 0),
    status         varchar(20)              default 'pending'::character varying not null
);


create table loan_schema.auditlog
(
    audit_id    serial
        primary key,
    table_name  varchar(50) not null,
    operation   varchar(10) not null,
    record_id   integer     not null,
    old_data    jsonb,
    new_data    jsonb,
    modified_by varchar(100),
    modified_at timestamp with time zone default CURRENT_TIMESTAMP
);


create function loan_schema.update_updated_at_column() returns trigger
    language plpgsql
as
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$;

alter function loan_schema.update_updated_at_column() owner to arizal;

create trigger update_loans_updated_at
    before update
    on loan_schema.loans
    for each row
    execute procedure loan_schema.update_updated_at_column();

create trigger update_paymentschedules_updated_at
    before update
    on loan_schema.paymentschedules
    for each row
    execute procedure loan_schema.update_updated_at_column();

create trigger update_payments_updated_at
    before update
    on loan_schema.payments
    for each row
    execute procedure loan_schema.update_updated_at_column();

create function loan_schema.log_audit_changes() returns trigger
    language plpgsql
as
$$
DECLARE
pk_name TEXT;
    pk_value INT;
BEGIN
    -- Retrieve the primary key column name for the table
SELECT a.attname
INTO pk_name
FROM pg_index i
         JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
WHERE i.indrelid = TG_RELID AND i.indisprimary
    LIMIT 1;

-- Handle DELETE operation
IF (TG_OP = 'DELETE') THEN
        EXECUTE format('SELECT $1.%I', pk_name) INTO pk_value USING OLD;
INSERT INTO loan_schema.AuditLog (table_name, operation, record_id, old_data, new_data, modified_by)
VALUES (TG_TABLE_NAME, TG_OP, pk_value, row_to_json(OLD), NULL, SESSION_USER);
RETURN OLD;

-- Handle UPDATE operation
ELSIF (TG_OP = 'UPDATE') THEN
        EXECUTE format('SELECT $1.%I', pk_name) INTO pk_value USING OLD;
INSERT INTO loan_schema.AuditLog (table_name, operation, record_id, old_data, new_data, modified_by)
VALUES (TG_TABLE_NAME, TG_OP, pk_value, row_to_json(OLD), row_to_json(NEW), SESSION_USER);
RETURN NEW;

-- Handle INSERT operation
ELSIF (TG_OP = 'INSERT') THEN
        EXECUTE format('SELECT $1.%I', pk_name) INTO pk_value USING NEW;
INSERT INTO loan_schema.AuditLog (table_name, operation, record_id, old_data, new_data, modified_by)
VALUES (TG_TABLE_NAME, TG_OP, pk_value, NULL, row_to_json(NEW), SESSION_USER);
RETURN NEW;
END IF;
END;
$$;

-- alter function loan_schema.log_audit_changes() owner to arizal;

create trigger loans_audit_trigger
    after insert or update or delete
                    on loan_schema.loans
                        for each row
                        execute procedure loan_schema.log_audit_changes();

create trigger paymentschedules_audit_trigger
    after insert or update or delete
                    on loan_schema.paymentschedules
                        for each row
                        execute procedure loan_schema.log_audit_changes();

create trigger payments_audit_trigger
    after insert or update or delete
                    on loan_schema.payments
                        for each row
                        execute procedure loan_schema.log_audit_changes();


-- Index, create index when necessary
-- 15. Create Index
-- uncomment this if column not frequently queried
-- create index when you really need it
-- CREATE INDEX IF NOT EXISTS idx_payments_loan_id ON loan_schema.Payments(loan_id);
-- CREATE INDEX IF NOT EXISTS idx_paymentschedules_loan_id ON loan_schema.PaymentSchedules(loan_id);
-- CREATE INDEX idx_due_date ON PaymentSchedules(due_date);
-- CREATE INDEX idx_grace_period_end ON PaymentSchedules(grace_period_end);

-- index for cron check delinquency
CREATE INDEX idx_due_date_paid ON loan_schema.paymentschedules (due_date, paid);

-- index for cron late fee checker
CREATE INDEX idx_grace_period_end_paid ON loan_schema.paymentschedules (grace_period_end, paid, late_fee_applied);

