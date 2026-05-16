-- =========================================================
-- TUJUAN SNAPSHOT
-- =========================================================

ALTER TABLE tujuan_opd
ADD COLUMN penetapan_id BIGINT NOT NULL;

ALTER TABLE tujuan_opd
ADD CONSTRAINT fk_tujuan_snapshot
FOREIGN KEY (penetapan_id)
REFERENCES penetapan_opd(id)
ON DELETE CASCADE;

ALTER TABLE tujuan_opd
ADD CONSTRAINT uq_tujuan_snapshot
UNIQUE(
    penetapan_id,
    kode_tujuan_opd
);

CREATE INDEX idx_tujuan_snapshot
ON tujuan_opd(
    penetapan_id
);


ALTER TABLE indikator_tujuan_opd
ADD CONSTRAINT uq_indikator_tujuan
UNIQUE(
    tujuan_opd_id,
    kode_indikator
);

ALTER TABLE target_indikator_tujuan_opd
ADD CONSTRAINT uq_target_tujuan
UNIQUE(
    indikator_tujuan_id,
    kode_target
);


-- =========================================================
-- SASARAN SNAPSHOT
-- =========================================================

ALTER TABLE sasaran_opd
ADD COLUMN penetapan_id BIGINT NOT NULL;

ALTER TABLE sasaran_opd
ADD CONSTRAINT fk_sasaran_snapshot
FOREIGN KEY (penetapan_id)
REFERENCES penetapan_opd(id)
ON DELETE CASCADE;

ALTER TABLE sasaran_opd
ADD CONSTRAINT uq_sasaran_snapshot
UNIQUE(
    penetapan_id,
    kode_sasaran_opd
);

CREATE INDEX idx_sasaran_snapshot
ON sasaran_opd(
    penetapan_id
);


ALTER TABLE indikator_sasaran_opd
ADD CONSTRAINT uq_indikator_sasaran
UNIQUE(
    sasaran_opd_id,
    kode_indikator
);

ALTER TABLE target_indikator_sasaran_opd
ADD CONSTRAINT uq_target_sasaran
UNIQUE(
    indikator_sasaran_id,
    kode_target
);
