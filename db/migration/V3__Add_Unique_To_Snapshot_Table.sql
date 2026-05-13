-- =========================================================
-- TUJUAN SNAPSHOT
-- =========================================================

ALTER TABLE tb_tujuan_penetapan_opd
ADD COLUMN penetapan_id BIGINT NOT NULL;

ALTER TABLE tb_tujuan_penetapan_opd
ADD CONSTRAINT fk_tujuan_snapshot
FOREIGN KEY (penetapan_id)
REFERENCES penetapan_opd(id)
ON DELETE CASCADE;

ALTER TABLE tb_tujuan_penetapan_opd
ADD UNIQUE(penetapan_id, kode_tujuan_opd);

CREATE INDEX idx_tujuan_snapshot
ON tb_tujuan_penetapan_opd(penetapan_id);

-- =========================================================
-- INDIKATOR TUJUAN SNAPSHOT
-- =========================================================

ALTER TABLE tb_indikator_tujuan_penetapan_opd
ADD COLUMN penetapan_id BIGINT NOT NULL;

ALTER TABLE tb_indikator_tujuan_penetapan_opd
ADD CONSTRAINT fk_indikator_tujuan_snapshot
FOREIGN KEY (penetapan_id)
REFERENCES penetapan_opd(id)
ON DELETE CASCADE;

ALTER TABLE tb_indikator_tujuan_penetapan_opd
ADD UNIQUE(penetapan_id, kode_indikator, tahun_aktif);

CREATE INDEX idx_indikator_tujuan_snapshot
ON tb_indikator_tujuan_penetapan_opd(penetapan_id);

-- =========================================================
-- TARGET TUJUAN SNAPSHOT
-- =========================================================

ALTER TABLE tb_target_indikator_tujuan_penetapan_opd
ADD COLUMN penetapan_id BIGINT NOT NULL;

ALTER TABLE tb_target_indikator_tujuan_penetapan_opd
ADD CONSTRAINT fk_target_tujuan_snapshot
FOREIGN KEY (penetapan_id)
REFERENCES penetapan_opd(id)
ON DELETE CASCADE;

CREATE INDEX idx_target_tujuan_snapshot
ON tb_target_indikator_tujuan_penetapan_opd(penetapan_id);

-- =========================================================
-- SASARAN SNAPSHOT
-- =========================================================

ALTER TABLE tb_sasaran_penetapan_opd
ADD COLUMN penetapan_id BIGINT NOT NULL;

ALTER TABLE tb_sasaran_penetapan_opd
ADD CONSTRAINT fk_sasaran_snapshot
FOREIGN KEY (penetapan_id)
REFERENCES penetapan_opd(id)
ON DELETE CASCADE;

ALTER TABLE tb_sasaran_penetapan_opd
ADD UNIQUE(penetapan_id, kode_sasaran_opd);

CREATE INDEX idx_sasaran_snapshot
ON tb_sasaran_penetapan_opd(penetapan_id);

-- =========================================================
-- INDIKATOR SASARAN SNAPSHOT
-- =========================================================

ALTER TABLE tb_indikator_sasaran_penetapan_opd
ADD COLUMN penetapan_id BIGINT NOT NULL;

ALTER TABLE tb_indikator_sasaran_penetapan_opd
ADD CONSTRAINT fk_indikator_sasaran_snapshot
FOREIGN KEY (penetapan_id)
REFERENCES penetapan_opd(id)
ON DELETE CASCADE;

ALTER TABLE tb_indikator_sasaran_penetapan_opd
ADD UNIQUE(penetapan_id, kode_indikator, tahun_aktif);


CREATE INDEX idx_indikator_sasaran_snapshot
ON tb_indikator_sasaran_penetapan_opd(penetapan_id);

-- =========================================================
-- TARGET SASARAN SNAPSHOT
-- =========================================================

ALTER TABLE tb_target_indikator_sasaran_penetapan_opd
ADD COLUMN penetapan_id BIGINT NOT NULL;

ALTER TABLE tb_target_indikator_sasaran_penetapan_opd
ADD CONSTRAINT fk_target_sasaran_snapshot
FOREIGN KEY (penetapan_id)
REFERENCES penetapan_opd(id)
ON DELETE CASCADE;

CREATE INDEX idx_target_sasaran_snapshot
ON tb_target_indikator_sasaran_penetapan_opd(penetapan_id);
