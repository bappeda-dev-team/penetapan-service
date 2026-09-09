ALTER TABLE indikator_pk
    ADD COLUMN kode_indikator_sasaran_opd VARCHAR(255);

ALTER TABLE target_pk
    ADD COLUMN kode_target_sasaran_opd VARCHAR(255);
