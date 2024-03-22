CREATE TABLE IF NOT EXISTS logs (
    Id Int64, ProjectId Int64, Name String, Description String, Priority Int64, Removed Boolean, EventTime DateTime
) ENGINE = Log;