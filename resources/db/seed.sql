-- Users
INSERT INTO users (id, login, password, is_moderator) VALUES
(1, 'ivan_ivanov', 'password123', false),
(2, 'petr_petrov', 'password123', true),
(3, 'anna_sidorova', 'password123', false)
ON CONFLICT (id) DO NOTHING;

-- Materials
INSERT INTO materials (id, title, formula, description, relative_molecular_mass, stoichiometric_coefficient, image_url) VALUES
(1, 'Известняк', 'CaCO3', 'Осадочная порода, состоящая преимущественно из кальцита (карбоната кальция).', 100.07, 1, 'izvestnyak.jpg'),
(2, 'Мрамор', 'CaCO3', 'Метаморфическая порода из кальцита, прочная, декоративная, полируемая.', 100.07, 1, 'mramor.jpg'),
(3, 'Металлический цинк', 'Zn', 'Голубовато-белый металл, пластичный, коррозионно-стойкий.', 65.39, 1, 'cink.jpg'),
(4, 'Сода', 'Na2CO3', 'Белый, без запаха, водорастворимый порошок или кристаллы, представляющие собой среднюю соль угольной кислоты', 105.99, 1, 'soda.jpg')
ON CONFLICT (id) DO NOTHING;

-- Experiments
-- Note: 'status' is a required field. Let's assume a default status 'draft'.
-- Creator ID is linked to the users table.
INSERT INTO experiments (id, creator_id, moderator_id, status, molar_volume, created_at, formed_at, finished_at) VALUES
(1, 1, NULL, 'draft', 22.4, NOW(), NULL, NULL),
(2, 1, NULL, 'formed', 22.4, NOW() - interval '2 day', NOW() - interval '1 day', NULL),
(3, 3, 2, 'finished', 22.4, NOW() - interval '7 day', NOW() - interval '6 day', NOW() - interval '2 day'),
(4, 3, NULL, 'formed', 22.4, NOW() - interval '11 day', NOW() - interval '10 day', NULL),
(5, 1, NULL, 'deleted', 22.4, NOW() - interval '15 day', NOW() - interval '14 day', NULL),
(6, 1, 2, 'rejected', 22.4, NOW() - interval '5 day', NOW() - interval '4 day', NOW() - interval '1 day')
ON CONFLICT (id) DO UPDATE SET
    creator_id = EXCLUDED.creator_id,
    moderator_id = EXCLUDED.moderator_id,
    status = EXCLUDED.status,
    molar_volume = EXCLUDED.molar_volume,
    created_at = EXCLUDED.created_at,
    formed_at = EXCLUDED.formed_at,
    finished_at = EXCLUDED.finished_at;

-- ExperimentItems
INSERT INTO experiment_items (id, experiment_id, material_id, material_mass, gas_volume, mass_fraction_percentage) VALUES
-- Items for experiment 1 (draft)
(1, 1, 1, 2.35, 0.484, 8),
(2, 1, 2, 1.8, 0.397, 1.5),
-- Items for experiment 2 (formed)
(3, 2, 3, 10.5, 2.1, 20),
-- Items for experiment 3 (moderated)
(4, 3, 1, 5.0, 1.0, 10),
(5, 3, 4, 12.0, 3.5, 25),
(6, 3, 3, 7.2, 1.8, 15),
-- Items for experiment 4 (old formed)
(7, 4, 2, 3.3, 0.8, 5),
-- Items for experiment 5 (deleted)
(8, 5, 1, 1.0, 0.2, 5),
-- Items for experiment 6 (rejected)
(9, 6, 4, 20.0, 5.0, 30),
(10, 6, 3, 15.0, 4.0, 22)
ON CONFLICT (id) DO UPDATE SET
    experiment_id = EXCLUDED.experiment_id,
    material_id = EXCLUDED.material_id,
    material_mass = EXCLUDED.material_mass,
    gas_volume = EXCLUDED.gas_volume,
    mass_fraction_percentage = EXCLUDED.mass_fraction_percentage;

-- Reset sequences to avoid conflicts with new data
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users), true);
SELECT setval('materials_id_seq', (SELECT MAX(id) FROM materials), true);
SELECT setval('experiments_id_seq', (SELECT MAX(id) FROM experiments), true);
SELECT setval('experiment_items_id_seq', (SELECT MAX(id) FROM experiment_items), true);
