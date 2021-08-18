CREATE TABLE IF NOT EXISTS Taxonomy (
	taxa_id INT PRIMARY KEY,
	Kingdom TEXT,
	Phylum TEXT,
	Class TEXT,
	Orders TEXT,
	Family TEXT,
	Genus TEXT,
	Species TEXT,
	Source TEXT
);

CREATE TABLE IF NOT EXISTS Common (
	taxa_id INT,
	Name TEXT,
	Curator TEXT,
	CONSTRAINT fk_taxonomy_common FOREIGN KEY (taxa_id) REFERENCES Taxonomy(taxa_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS Denominators (
	taxa_id INT PRIMARY KEY,
	Noncancer INT,
	CONSTRAINT fk_taxonomy_denominators FOREIGN KEY (taxa_id) REFERENCES Taxonomy(taxa_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS Life_history (
	taxa_id INT PRIMARY KEY,
	female_maturity DOUBLE,
	male_maturity DOUBLE,
	Gestation DOUBLE,
	Weaning DOUBLE,
	Infancy DOUBLE,
	litter_size DOUBLE,
	litters_year DOUBLE,
	interbirth_interval DOUBLE,
	birth_weight DOUBLE,
	weaning_weight DOUBLE,
	adult_weight DOUBLE,
	growth_rate DOUBLE,
	max_longevity DOUBLE,
	metabolic_rate DOUBLE,
	CONSTRAINT fk_taxonomy_lifehistory FOREIGN KEY (taxa_id) REFERENCES Taxonomy(taxa_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS Accounts (
	account_id INT PRIMARY KEY,
	Account TEXT,
	submitter_name TEXT
);

CREATE TABLE IF NOT EXISTS Patient (
	ID INT PRIMARY KEY,
	Sex TEXT,
	Age DECIMAL(6,2),
	Infant TINYINT,
	Castrated TINYINT,
	Wild TINYINT,
	taxa_id INT,
	source_id TEXT,
	source_name TEXT,
	Date TEXT,
	Year INT,
	Comments TEXT,
	CONSTRAINT fk_taxonomy_patient FOREIGN KEY (taxa_id) REFERENCES Taxonomy(taxa_id) ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS Diagnosis (
	ID INT,
	Masspresent TINYINT,
	Hyperplasia TINYINT,
	Necropsy TINYINT,
	Metastasis TINYINT,
	CONSTRAINT fk_patient_diagnosis FOREIGN KEY (ID) REFERENCES Patient(ID) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS Source (
	ID INT,
	service_name TEXT,
	Zoo TINYINT,
	Aza TINYINT,
	Institute TINYINT,
	Approved TINYINT,
	account_id INT,
	CONSTRAINT fk_patient_source FOREIGN KEY (ID) REFERENCES Patient(ID) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_accounts_source FOREIGN KEY (account_id) REFERENCES Accounts(account_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS Tumor (
	ID INT,
	primary_tumor TINYINT,
	Malignant TINYINT,
	Type TEXT,
	Location TEXT,
	CONSTRAINT fk_patient_tumor FOREIGN KEY (ID) REFERENCES Patient(ID) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS Unmatched (
	sourceID TEXT,
	name TEXT,
	sex TEXT,
	age DECIMAL(6,2),
	date TEXT,
	masspresent TINYINT,
	necropsy TINYINT,
	comments TEXT,
	Service TEXT
);

CREATE TABLE IF NOT EXISTS Update_time (
	update_number INT PRIMARY KEY AUTO_INCREMENT,
	Time TEXT
);

CREATE INDEX IX_taxonomy_taxaid ON Taxonomy (taxa_id);
CREATE INDEX IX_lifehistory_taxaid ON Life_history (taxa_id);
CREATE INDEX IX_patient_taxaid ON Patient (taxa_id);
CREATE INDEX IX_patient_id ON Patient (ID);
CREATE INDEX IX_diagnosis_id ON Diagnosis (ID);
CREATE INDEX IX_source_id ON Source (ID);
CREATE INDEX IX_tumor_id ON Tumor (ID);

CREATE OR REPLACE VIEW Records AS
	SELECT
		Patient.ID,
		Patient.Sex,
		Patient.Age AS age_months,
		Patient.Infant,
		Patient.Castrated,
		Patient.Wild,
		Patient.taxa_id,
		Patient.source_id,
		Patient.source_name,
		Patient.Date,
		Patient.Year,
		Patient.Comments,
		Diagnosis.Masspresent,
		Diagnosis.Hyperplasia,
		Diagnosis.Necropsy,
		Diagnosis.Metastasis,
		Tumor.primary_tumor,
		Tumor.Malignant,
		Tumor.Type,
		Tumor.Location,
		Taxonomy.Kingdom,
		Taxonomy.Phylum,
		Taxonomy.Class,
		Taxonomy.Orders,
		Taxonomy.Family,
		Taxonomy.Genus,
		Taxonomy.Species,
		Source.service_name,
		Source.Zoo,
		Source.Aza,
		Source.Institute,
		Source.Approved,
		Life_history.female_maturity,
		Life_history.male_maturity,
		Life_history.Gestation,
		Life_history.Weaning,
		Life_history.Infancy,
		Life_history.litter_size,
		Life_history.litters_year,
		Life_history.interbirth_interval,
		Life_history.birth_weight,
		Life_history.weaning_weight,
		Life_history.adult_weight,
		Life_history.growth_rate,
		Life_history.max_longevity,
		Life_history.metabolic_rate
	FROM Patient
		INNER JOIN Diagnosis on Diagnosis.ID = Patient.ID
		LEFT JOIN Tumor on Tumor.ID = Patient.ID
		INNER JOIN Taxonomy on Taxonomy.taxa_id = Patient.taxa_id
		INNER JOIN Source on Source.ID = Patient.ID
		LEFT JOIN Life_history on Life_history.taxa_id = Patient.taxa_id
	ORDER BY taxa_id
;
