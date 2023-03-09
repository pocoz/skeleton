CREATE TABLE IF NOT EXISTS statistics
(
    advert_id         varchar(255),
    advert_name       varchar(255),
    advertiser_name   varchar(255),
    campaign_name     varchar(255),
    views             float,
    clicks            float,
    action_date       date,
    advert_start_date date,
    advert_end_date   date,
    PRIMARY KEY (advert_id, advert_name, action_date)
    );
