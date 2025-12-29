-- Генерация 100,000 образцов для демонстрации индексов
-- Названия и формулы генерируются случайно

INSERT INTO acid_soluble_samples (title, formula, description, deleted, relative_molecular_mass, stoichiometric_coefficient)
SELECT 
    'Sample_' || gs || '_' || (ARRAY['Oxide', 'Sulfate', 'Chloride', 'Nitrate', 'Carbonate', 'Hydroxide', 'Phosphate', 'Acetate'])[1 + floor(random() * 8)::int],
    (ARRAY['H2SO4', 'HCl', 'HNO3', 'H3PO4', 'CH3COOH', 'NaOH', 'KOH', 'NH4OH', 'Ca(OH)2', 'Mg(OH)2'])[1 + floor(random() * 10)::int],
    'Автогенерированный образец #' || gs || '. Используется для тестирования пагинации и индексов.',
    false,
    ROUND((50 + random() * 200)::numeric, 2),  -- relative_molecular_mass: 50-250
    ROUND((1 + random() * 3)::numeric, 2)      -- stoichiometric_coefficient: 1-4
FROM generate_series(1, 100000) AS gs;

-- Обновляем последовательность ID
SELECT setval(pg_get_serial_sequence('acid_soluble_samples', 'id'), COALESCE(MAX(id), 0)) FROM acid_soluble_samples;

-- Вывод статистики
SELECT 'Всего записей: ' || COUNT(*) FROM acid_soluble_samples;
