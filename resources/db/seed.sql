-- Users
INSERT INTO users (id, login, password, is_moderator) VALUES
(1, 'ivan_ivanov', 'password123', true),
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
INSERT INTO impurity_fraction_experiments (id, creator_id, status, molar_volume, created_at, formed_at, finished_at) VALUES
(1, 1, 'draft', 22.4, NOW(), NULL, NULL),
(2, 1, 'formed', 22.4, NOW() - interval '2 day', NOW() - interval '1 day', NULL),
(3, 3, 'finished', 22.4, NOW() - interval '7 day', NOW() - interval '6 day', NOW() - interval '2 day'),
(4, 3, 'formed', 22.4, NOW() - interval '11 day', NOW() - interval '10 day', NULL),
(5, 1, 'deleted', 22.4, NOW() - interval '15 day', NOW() - interval '14 day', NULL),
(6, 1, 'rejected', 22.4, NOW() - interval '5 day', NOW() - interval '4 day', NOW() - interval '1 day')
ON CONFLICT (id) DO NOTHING;

-- ExperimentsSamples
INSERT INTO experiments_samples (experiment_id, sample_id, sample_mass, evolved_gas_volume, mass_fraction_percentage) VALUES
-- Samples for experiment 1 (draft)
(1, 1, 2.35, 0.484, 8),
(1, 2, 1.8, 0.397, 1.5),
-- Samples for experiment 2 (formed)
(2, 3, 10.5, 2.1, 20),
-- Samples for experiment 3 (finished)
(3, 1, 5.0, 1.0, 10),
(3, 4, 12.0, 3.5, 25),
(3, 3, 7.2, 1.8, 15),
-- Samples for experiment 4 (formed)
(4, 2, 3.3, 0.8, 5),
-- Samples for experiment 5 (deleted)
(5, 1, 1.0, 0.2, 5),
-- Samples for experiment 6 (rejected)
(6, 4, 20.0, 5.0, 30),
(6, 3, 15.0, 4.0, 22)
ON CONFLICT (experiment_id, sample_id) DO NOTHING;

-- Sync sequences to the current MAX(id) to avoid duplicate key errors on next inserts
SELECT setval(pg_get_serial_sequence('users', 'id'), COALESCE(MAX(id), 0)) FROM users;
SELECT setval(pg_get_serial_sequence('acid_soluble_samples', 'id'), COALESCE(MAX(id), 0)) FROM acid_soluble_samples;
SELECT setval(pg_get_serial_sequence('impurity_fraction_experiments', 'id'), COALESCE(MAX(id), 0)) FROM impurity_fraction_experiments;
