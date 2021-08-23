CREATE TABLE IF NOT EXISTS mockdata
(
    int_col integer NOT NULL, 
    float_col double precision,
    string_col text NOT NULL,
    bool_col boolean,
    PRIMARY KEY (int_col)
);
CREATE INDEX IF NOT EXISTS mock_data_string_index ON mockdata (string_col);