CREATE TABLE link_ads
(
    link_id BIGINT NOT NULL REFERENCES links (id) ON DELETE CASCADE,
    ad_id   BIGINT
        CONSTRAINT link_ads_ad_id_key NOT NULL REFERENCES ads (id) ON DELETE CASCADE,
    PRIMARY KEY (link_id, ad_id)
);
