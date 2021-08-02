CREATE TABLE Accounts (
	account_id INT PRIMARY KEY,
	Account TEXT,
	submitter_name TEXT
);

CREATE TABLE Common (
	taxa_id INT,
	Name TEXT,
	Curator TEXT,
	CONSTRAINT fk_taxonomy_common FOREIGN KEY (taxa_id) REFERENCES Taxonomy (taxa_id) ON DELETE CASCADE ON UPDATE CASCADE,
);

CREATE TABLE Denominators (
	taxa_id INT PRIMARY KEY,
	Noncancer INT,
	CONSTRAINT fk_taxonomy_denominators FOREIGN KEY (taxa_id) REFERENCES Taxonomy (taxa_id) ON DELETE CASCADE ON UPDATE CASCADE,
);

CREATE TABLE Diagnosis (
	ID INT,
	Masspresent TINYINT,
	Hyperplasia TINYINT,
	Necropsy TINYINT,
	Metastasis TINYINT,
	CONSTRAINT fk_patient_diagnosis FOREIGN KEY (ID) REFERENCES Patient (ID) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE Life_history (
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
	CONSTRAINT fk_taxonomy_lifehistory FOREIGN KEY (taxa_id) REFERENCES Taxonomy (taxa_id) ON DELETE CASCADE ON UPDATE CASCADE,
);

CREATE TABLE Patient (
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
	CONSTRAINT fk_taxonomy_patient FOREIGN KEY (taxa_id) REFERENCES Taxonomy (taxa_id) ON UPDATE CASCADE,
);

CREATE TABLE Source (
	ID INT,
	service_name TEXT,
	Zoo TINYINT,
	Aza TINYINT,
	Institute TINYINT,
	Approved TINYINT,
	account_id INT,
	CONSTRAINT fk_patient_source FOREIGN KEY (ID) REFERENCES Patient (ID) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT fk_accounts_source FOREIGN KEY (account_id) REFERENCES Accounts (account_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE Taxonomy (
	taxa_id INT PRIMARY KEY,
	Kingdom TEXT,
	Phylum TEXT,
	Class TEXT,
	Orders TEXT,
	Family TEXT,
	Genus TEXT,
	Species TEXT,
	Source TEXT,
	INDEX IX_taxonomy_species (taxa_id)

);

CREATE TABLE Tumor (
	ID INT,
	primary_tumor TINYINT,
	Malignant TINYINT,
	Type TEXT,
	Location TEXT,
	CONSTRAINT fk_patient_tumor FOREIGN KEY (ID) REFERENCES Patient (ID) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE Unmatched (
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

CREATE TABLE Update_time (
	update_number INT PRIMARY KEY AUTO_INCREMENT,
	Time TEXT
);
