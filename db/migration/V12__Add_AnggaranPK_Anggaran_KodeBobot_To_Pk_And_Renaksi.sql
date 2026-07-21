ALTER TABLE pk_individu
    ADD COLUMN anggaran_pk BIGINT NOT NULL DEFAULT 0;

ALTER TABLE renaksi_individu
    ADD COLUMN anggaran BIGINT NOT NULL DEFAULT 0;

-- PELAKSANAAN RENAKSI
ALTER TABLE renaksi_individu_pelaksanaan
    ADD COLUMN kode_pelaksanaan VARCHAR(255);

UPDATE renaksi_individu_pelaksanaan
    SET kode_pelaksanaan = 'PEL-' || id
    WHERE kode_pelaksanaan IS NULL;
