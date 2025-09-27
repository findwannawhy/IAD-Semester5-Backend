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
INSERT INTO experiments (id, creator_id, status, molar_volume, created_at) VALUES
(1, 1, 'draft', 22.4, NOW())
ON CONFLICT (id) DO NOTHING;

-- ExperimentItems
INSERT INTO experiment_items (id, experiment_id, material_id, material_mass, gas_volume, mass_fraction_percentage) VALUES
(1, 1, 1, 2.35, 0.484, 8),
(2, 1, 2, 1.8, 0.397, 1.5)
ON CONFLICT (id) DO NOTHING;

-- Reset sequences to avoid conflicts with new data
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users), true);
SELECT setval('materials_id_seq', (SELECT MAX(id) FROM materials), true);
SELECT setval('experiments_id_seq', (SELECT MAX(id) FROM experiments), true);
SELECT setval('experiment_items_id_seq', (SELECT MAX(id) FROM experiment_items), true);