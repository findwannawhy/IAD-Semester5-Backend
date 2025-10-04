-- Users
INSERT INTO users (id, login, password, is_moderator) VALUES
(1, 'ivan_ivanov', 'password123', false),
(2, 'petr_petrov', 'password123', true),
(3, 'anna_sidorova', 'password123', false)
ON CONFLICT (id) DO NOTHING;

-- AcidSolubleSamples
INSERT INTO acid_soluble_samples (id, title, formula, description, relative_molecular_mass, stoichiometric_coefficient, image_url) VALUES
(1, 'Известняк', 'CaCO3', 'Осадочная порода, состоящая преимущественно из кальцита (карбоната кальция).', 100.07, 1, 'izvestnyak.jpg'),
(2, 'Мрамор', 'CaCO3', 'Метаморфическая порода из кальцита, прочная, декоративная, полируемая.', 100.07, 1, 'mramor.jpg'),
(3, 'Металлический цинк', 'Zn', 'Голубовато-белый металл, пластичный, коррозионно-стойкий.', 65.39, 1, 'cink.jpg'),
(4, 'Сода', 'Na2CO3', 'Белый, без запаха, водорастворимый порошок или кристаллы, представляющие собой среднюю соль угольной кислоты', 105.99, 1, 'soda.jpg')
ON CONFLICT (id) DO NOTHING;

-- ImpurityFractionExperiments
-- Note: 'status' is a required field. Let's assume a default status 'draft'.
-- Creator ID is linked to the users table.
INSERT INTO impurity_fraction_experiments (id, creator_id, status, molar_volume, created_at) VALUES
(1, 1, 'draft', 22.4, NOW())
ON CONFLICT (id) DO NOTHING;

-- ExperimentsSamples
INSERT INTO experiments_samples (experiment_id, sample_id, sample_mass, evolved_gas_volume, mass_fraction_percentage) VALUES
(1, 1, 2.35, 0.484, 8),
(1, 2, 1.8, 0.397, 1.5)
ON CONFLICT (experiment_id, sample_id) DO NOTHING;

-- Sync sequences to the current MAX(id) to avoid duplicate key errors on next inserts
SELECT setval(pg_get_serial_sequence('users', 'id'), COALESCE(MAX(id), 0)) FROM users;
SELECT setval(pg_get_serial_sequence('acid_soluble_samples', 'id'), COALESCE(MAX(id), 0)) FROM acid_soluble_samples;
SELECT setval(pg_get_serial_sequence('impurity_fraction_experiments', 'id'), COALESCE(MAX(id), 0)) FROM impurity_fraction_experiments;