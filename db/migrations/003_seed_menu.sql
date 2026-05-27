-- ============================================================
-- STEP 1 — menu_categories
-- ============================================================
INSERT INTO menu_categories (id, name, sort_order) VALUES
  ('11111111-1111-4111-a111-111111111101', 'Vorspeisen',       1),
  ('11111111-1111-4111-a111-111111111102', 'Suppen & Salate',  2),
  ('11111111-1111-4111-a111-111111111103', 'Steaks',           3),
  ('11111111-1111-4111-a111-111111111104', 'Grillplatten',     4),
  ('11111111-1111-4111-a111-111111111105', 'Hauptspeisen',     5),
  ('11111111-1111-4111-a111-111111111106', 'Fischgerichte',    6),
  ('11111111-1111-4111-a111-111111111107', 'Nudeln & Burger',  7),
  ('11111111-1111-4111-a111-111111111108', 'Desserts',         8),
  ('11111111-1111-4111-a111-111111111109', 'Offene Weine',     9),
  ('11111111-1111-4111-a111-111111111110', 'Flaschenweine',   10)
ON CONFLICT (id) DO NOTHING;


-- ============================================================
-- STEP 2 — menu_items
-- ============================================================

-- VORSPEISEN
INSERT INTO menu_items (id, category_id, name, description, price, is_available, allergens, tags, sort_order) VALUES
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'El Nigo Platte Mix-Vorspeise',  'Knoblauchbrot, Alioli, Putenspieß, Dateln, Chorizo, Hähnchenpflügel, Mini-Mozzarella, Hackfleischbällchen, Fettakäse, Lammkrone (für 2 Personen)', 35.80, true, '{}', ARRAY['sharing'], 1),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Haus gemachte Chorizo Rosario', '', 9.50, true, '{}', '{}', 2),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Mediterrane Oliven',            '', 6.00, true, '{}', ARRAY['vegetarisch'], 3),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Serrano Schinken',              '', 11.90, true, '{}', '{}', 4),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Haus Platte',                   'Käse, Kartoffel, Speck und Kajmak Dip', 17.10, true, '{}', '{}', 5),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Gemüse Mediterran',             'Dazu Knoblauchbrot', 10.40, true, '{}', ARRAY['vegetarisch'], 6),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Gambas in Knoblauchsoße',       'Pikant', 10.80, true, '{}', '{}', 7),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Panierte Champignons',          'Mit Alioli Dip', 8.80, true, '{}', ARRAY['vegetarisch'], 8),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Knoblauchbrot mit Alioli',      '', 5.60, true, '{}', ARRAY['vegetarisch'], 9),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Hähnchenpflügel mariniert',     'In BBQ Soße', 6.90, true, '{}', '{}', 10),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Gebratene Sardellen',           'Mit Alioli Dip', 9.80, true, '{}', '{}', 11),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Octopus Salat',                 'Kartoffeln und Oliven', 14.10, true, '{}', '{}', 12),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Hähnchenspiess',                'In Erdnussbuttersoße und Mandeln', 7.80, true, '{}', '{}', 13),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Gefüllte Aubergine',            'Mit Schafskäse in Tomatensoße', 9.10, true, '{}', ARRAY['vegetarisch'], 14),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Flambierter Schafskäse',        'Mit Chili und Honig', 9.10, true, '{}', ARRAY['vegetarisch'], 15),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Parmesan Kartoffelecken',       'Mit geriebenem Käse', 7.40, true, '{}', ARRAY['vegetarisch'], 16),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111101', 'Datteln im Speckmantel',        'Mit Alioli Dip', 7.90, true, '{}', '{}', 17);

-- SUPPEN & SALATE
INSERT INTO menu_items (id, category_id, name, description, price, is_available, allergens, tags, sort_order) VALUES
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111102', 'Tomatencremesuppe',       '', 7.20, true, '{}', ARRAY['suppe','vegetarisch'], 1),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111102', 'Fischsuppe',              '', 7.20, true, '{}', ARRAY['suppe'], 2),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111102', 'Hühnersuppe',             '', 6.50, true, '{}', ARRAY['suppe'], 3),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111102', 'Salat mit Hähnchen',      'Bunter Salat mit panierten Hähnchenbruststreifen und Aliolisoße', 16.50, true, '{}', ARRAY['salat'], 4),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111102', 'Salat mit Rinderstreifen','Bunter Salat, Knoblauchbrot und BBQ Soße', 21.50, true, '{}', ARRAY['salat'], 5),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111102', 'Vegetarischer Salat',     'Bunter Salat mit Gemüse, gegrillten Oliven, Schafskäse und Avocado Dip', 15.90, true, '{}', ARRAY['salat','vegetarisch'], 6);

-- STEAKS
INSERT INTO menu_items (id, category_id, name, description, price, is_available, allergens, tags, sort_order) VALUES
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111103', 'Filetsteak 180g',    '', 26.10, true, '{}', ARRAY['filet'], 1),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111103', 'Filetsteak 250g',    '', 35.50, true, '{}', ARRAY['filet'], 2),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111103', 'Filetsteak 300g',    '', 40.60, true, '{}', ARRAY['filet'], 3),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111103', 'Rumpsteak 180g',     '', 22.40, true, '{}', ARRAY['rumpsteak'], 4),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111103', 'Rumpsteak 250g',     '', 29.70, true, '{}', ARRAY['rumpsteak'], 5),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111103', 'Rumpsteak 300g',     '', 33.30, true, '{}', ARRAY['rumpsteak'], 6),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111103', 'Rib-Eye Steak 300g', '', 31.10, true, '{}', ARRAY['ribeye'], 7),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111103', 'Rib-Eye Steak 400g', '', 38.20, true, '{}', ARRAY['ribeye'], 8),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111103', 'Rib-Eye Steak 500g', '', 42.60, true, '{}', ARRAY['ribeye'], 9);

-- GRILLPLATTEN
INSERT INTO menu_items (id, category_id, name, description, price, is_available, allergens, tags, sort_order) VALUES
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111104', 'El Nigo Grillplatte',        'Rindersteak, Cordon Bleu, gefüllte Pljeskavica, Speck & Raznjici, Geschmortes Gemüse, Pommes & Djuvecreis (für 2 Personen)', 46.80, true, '{}', ARRAY['sharing'], 1),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111104', 'El Nigo Gourmet Grillplatte', 'Rindersteak, Schnitzel (Kalb) Wiener Art, Hähnchen-Raznjici, Lammkrone, Hähnchensteak, Geschmortes Gemüse, Pommes & Djuvecreis (für 2 Personen)', 52.60, true, '{}', ARRAY['sharing'], 2),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111104', 'Surf & Turf',                '2 Filet Mignon, 3 Scampi serviert mit Tagliatelle in einer Tomatensoße', 37.20, true, '{}', '{}', 3);

-- HAUPTSPEISEN
INSERT INTO menu_items (id, category_id, name, description, price, is_available, allergens, tags, sort_order) VALUES
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'El Nigo Schnitzel',            'Mit Bratkartoffeln und Frankfurter Grüner Soße', 16.30, true, '{}', ARRAY['schnitzel'], 1),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Wiener Schnitzel',             'Vom Kalb mit Bratkartoffeln und Preiselbeersoße', 23.50, true, '{}', ARRAY['schnitzel','kalb'], 2),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Jägerschnitzel',               'Mit Bratkartoffeln', 17.50, true, '{}', ARRAY['schnitzel'], 3),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Cordon Bleu',                  'Vom Schwein gefüllt mit Schinken und Käse, dazu Pommes', 19.70, true, '{}', ARRAY['schnitzel'], 4),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Kalbsschnitzel Natur',         'In einer Jägersoße, dazu Bratkartoffeln', 27.30, true, '{}', ARRAY['schnitzel','kalb'], 5),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Paniertes Hähnchenbrustfilet', 'Dazu Pommes', 16.30, true, '{}', ARRAY['hähnchen'], 6),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Hähnchenschnitzel',            'Vom Grill mit Djuvecreis und Pommes', 18.60, true, '{}', ARRAY['hähnchen','schnitzel'], 7),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Hähnchenkeule Grill',          'Ohne Knochen mit Djuvecreis und Pommes', 18.40, true, '{}', ARRAY['hähnchen'], 8),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Hähnchenbrustfilet',           'Überbacken mit Spinat, Mozzarella dazu Rosmarinkartoffel', 21.30, true, '{}', ARRAY['hähnchen'], 9),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'El Nigo Rumpsteak',            'Vom Lavasteingrill mit Kräuterbutter und Pommes', 28.80, true, '{}', ARRAY['rind'], 10),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'El Nigo Grillteller',          'Rindersteak, Hähnchensteak, Schweinesteak und Speck, dazu geschmortes Gemüse, Pommes & Kräuterbutter', 21.90, true, '{}', ARRAY['mix'], 11),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Rumpsteak Gefüllt',            'Mit Schafskäse und Serrano-Schinken, dazu Pommes und Kräuterbutter', 31.50, true, '{}', ARRAY['rind'], 12),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Raznjici',                     'Vom Schwein mit Djuvecreis und Pommes', 17.50, true, '{}', ARRAY['schwein'], 13),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Pljeskavica Spezial',          'Gefüllt mit Schafskäse und Serrano-Schinken, dazu Djuvecreis und Pommes', 18.60, true, '{}', '{}', 14),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Pljeskavica Grill',            'Mit Speck und Kräuterbutter, dazu Djuvecreis und Pommes', 17.50, true, '{}', '{}', 15),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Rinder Pfefferspieß',          'Mit grüner Pfeffersoße und Bratkartoffeln', 25.90, true, '{}', ARRAY['rind'], 16),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'El Nigo Pfanne',               'Div. Filets in Rotweinsaße, Pflaumen, Champignons mit Bratkartoffeln', 27.30, true, '{}', ARRAY['rind'], 17),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Duroc Schweine-Kotelett',      'Vom Lavasteingrill mit Rosmarinkartoffeln, Sauercremesoße und Kräuterbutter', 22.90, true, '{}', ARRAY['schwein'], 18),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Wildschein Spareribs',         'Mariniert und Honig & BBQ überbacken, dazu Pommes und BBQ-Soße', 23.90, true, '{}', ARRAY['schwein'], 19),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Lammkrone Grill',              'Rosmarinkartoffeln und geschmortes Gemüse', 34.80, true, '{}', ARRAY['lamm'], 20),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111105', 'Ladies-Steak',                 'Filetsteak ~120–140g mit Bernaissoße und Folienkartoffel', 23.60, true, '{}', ARRAY['filet'], 21);

-- FISCHGERICHTE
INSERT INTO menu_items (id, category_id, name, description, price, is_available, allergens, tags, sort_order) VALUES
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111106', 'Lachsfilet vom Grill',    'Dazu Runzelkartoffeln und Gemüse', 23.40, true, '{}', '{}', 1),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111106', 'Tintenfisch vom Grill',   'Dazu Kartoffel und Spinat', 23.90, true, '{}', '{}', 2),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111106', 'Tintenfisch frittiert',   'Mit Pommes und Aliolisauce', 21.40, true, '{}', '{}', 3),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111106', 'Riesengarnelen / Scampi', 'Mit Tagliatelle in Tomatensoße', 25.70, true, '{}', '{}', 4),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111106', 'El Nigo Seafood Platte',  'Gambas in Knoblauchsoße, Tintenfisch, Octopussalat, Sardellen, Alioli (für 2 Personen)', 38.30, true, '{}', ARRAY['sharing'], 5);

-- NUDELN & BURGER
INSERT INTO menu_items (id, category_id, name, description, price, is_available, allergens, tags, sort_order) VALUES
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111107', 'Tagliarini Aglio e Olio', 'Mit Gambas, Cherry-Tomaten, Knoblauch und Chilisoße', 17.80, true, '{}', ARRAY['pasta'], 1),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111107', 'Tagliatelle Polo',        'Mit Hähnchenstreifen in Tomatensoße und Pesto Genovese, mit geriebenem Käse', 17.10, true, '{}', ARRAY['pasta'], 2),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111107', 'Tagliatelle Lachs',       'Mit einer Sahne-Tomatensoße', 17.80, true, '{}', ARRAY['pasta'], 3),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111107', 'Hamburger',               'Mit hausgemachter Soße, Gurke, Tomaten und Salat. Pommes und BBQ Soße', 15.90, true, '{}', ARRAY['burger'], 4),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111107', 'Burger Special',          'Mit paniertem Schafskäse, Tomaten, Gurke und Salat. Pommes und BBQ Soße', 17.60, true, '{}', ARRAY['burger'], 5),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111107', 'Burger Vegie',            'Mit hausgemachtem Hummus, Falafel, Tomaten, Gurke, Salat und Pommes', 15.90, true, '{}', ARRAY['burger','vegetarisch'], 6),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111107', 'Burger Hähnchen',         'Käse, Avocado und Pommes', 16.70, true, '{}', ARRAY['burger','hähnchen'], 7);

-- DESSERTS
INSERT INTO menu_items (id, category_id, name, description, price, is_available, allergens, tags, sort_order) VALUES
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111108', 'Palacinke',      'Wahlweise Nutella, Eis oder Heiße Liebe', 7.20, true, '{}', '{}', 1),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111108', 'Schoko Soufflé', 'Vanilleeis und Vanillesoße', 8.50, true, '{}', '{}', 2),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111108', 'Apfelstrudel',   'Vanilleeis und Vanillesoße', 8.50, true, '{}', '{}', 3),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111108', 'Semifredo',      'Hausgemacht', 7.20, true, '{}', '{}', 4),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111108', 'Kugel Eis',      'Schokolade, Vanille, Erdbeere', 1.70, true, '{}', '{}', 5);

-- OFFENE WEINE (price = glass 0,2l)
INSERT INTO menu_items (id, category_id, name, description, price, is_available, allergens, tags, sort_order) VALUES
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111109', 'Chardonnay Dika',         'Slawonien. Harmonischer Weißwein mit Aromen von Orangen, Banane und Buttercreme. 0,2l — Flasche €26', 7.00, true, '{}', ARRAY['weisswein'], 1),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111109', 'g.tocka by Galic',        'Slawonien. Aromen von grünem Apfel, reifer Quitte und Kamille. Straffe Struktur. 0,2l — Flasche €30', 8.00, true, '{}', ARRAY['weisswein'], 2),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111109', 'Rosé Premium Kutjevo',    'Eleganter Rosé mit Aromen von Erdbeere, Granatapfel und feiner Mineralität. 0,2l — Flasche €30', 8.00, true, '{}', ARRAY['rose'], 3),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111109', 'Cabernet Sauvignon Dika', 'Slawonien. Kraftvoller Cabernet mit Noten von schwarzen Johannisbeeren und Eichenholz. 0,2l — Flasche €26', 7.00, true, '{}', ARRAY['rotwein'], 4),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111109', 'Terra Rossa Laguna',      'Istrien. Würziger Rotwein mit Waldfrüchten und reifen Kirschen. Dicht und strukturiert. 0,2l — Flasche €30', 8.00, true, '{}', ARRAY['rotwein'], 5),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111109', 'Plavac Skaramuca',        'Halbinsel Peljesac. Rotwein mit Aromen von dunklen Beeren, Pflaumen und Kakao. 0,2l — Flasche €28', 7.00, true, '{}', ARRAY['rotwein'], 6);

-- FLASCHENWEINE (price = bottle 0,75l)
INSERT INTO menu_items (id, category_id, name, description, price, is_available, allergens, tags, sort_order) VALUES
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111110', 'Stina Bijeli Cuvée',              'Dicht strukturiert, frisch-fruchtiger Geschmack mit Aromen von getrockneten Feigen und Aprikosen.', 27.00, true, '{}', ARRAY['weisswein'], 1),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111110', 'G*Tocka by Galic',                'Außergewöhnliche Frische, fruchtige Aromen vom grünen Apfel sowie Noten von Weinbergfirsich und Kamille.', 26.00, true, '{}', ARRAY['weisswein'], 2),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111110', 'Grasevina Galic',                 'Feine fruchtige Aromen von Ananas, grünen Äpfeln und Weinbergfirsich. Aromen von Kamille und Heu.', 35.00, true, '{}', ARRAY['weisswein'], 3),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111110', 'Xtriana Malvazija Organic Veralda','Kristalline Farbe. Frisch-fruchtiger Duft nach Traubenpfirsich, Sommerbirne, Basilikum und Holunder.', 45.00, true, '{}', ARRAY['weisswein'], 4),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111110', 'Malvazija Kozlovic',              'Blumige Note, Grapefruit und Zitrone. Am Gaumen trocken und volle Frische. Funkelnd, milden Noten.', 39.00, true, '{}', ARRAY['weisswein'], 5),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111110', 'Zlatna Vrbnicka Zlathina',        'Blassgelb mit angenehmen Akazienaroma. Passt zu Meeresfrüchten und leichten Vorspeisen.', 28.00, true, '{}', ARRAY['weisswein'], 6),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111110', 'Posip Stina Majstor',             'Komplexer, zarter und eleganter Wein, hauch von Pfirsich und Apfel. Gleichgewicht Frische und Säure.', 55.00, true, '{}', ARRAY['weisswein'], 7),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111110', 'Rosé Kutjevo Premium',            'Sanftes Lachsros. Beerenfrüchte, Preiselbeeren und frische Kräuter. Kraftvoll und seidig.', 35.00, true, '{}', ARRAY['rose'], 8),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111110', 'Stina Crni Cuvée',                'Lebendige tiefrote Farbe mit violetten Reflexen. Am Gaumen elegant mit außerordentlicher Komplexität.', 31.00, true, '{}', ARRAY['rotwein'], 9),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111110', 'Plavac Skaramuca',                'Intensives dunkles Rot. Am Gaumen schwarze Kirschen mit Noten von dunklen Beeren. Sehr vollmundig.', 27.00, true, '{}', ARRAY['rotwein'], 10),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111110', 'Korlat Cabernet Sauvignon',       'Tiefe dunkle Farbe, vollmundig, harmonisch. Am Gaumen Pflaume, Heidelbeere und schwarze Johannisbeeren.', 59.00, true, '{}', ARRAY['rotwein'], 11),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111110', 'Dingac Skarmuca',                 'Tiefrote Farbe. Am Gaumen Kirsche, Kaffee, Vanille, Trüffel und dalmatinischen Graskitzel.', 32.00, true, '{}', ARRAY['rotwein'], 12),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111110', 'Plavac Stina Majstor',            'Saftiger und enorm fruchtiger Plavac. Aromen von schwarzen Kirschen, Brombeere, Himbeeren. Mundfüllend.', 69.00, true, '{}', ARRAY['rotwein'], 13),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111110', 'Prosper Peljesac',                'Leuchtend rubinrote Farbe. Üppiges Bukket mit fruchtigen Aromen.', 32.00, true, '{}', ARRAY['dessertwein'], 14),
  (gen_random_uuid(), '11111111-1111-4111-a111-111111111110', 'Prosek Jakov Sibenik',            'Duftet nach getrockneten Früchten und Honig, leicht ausgeprägte Säure. Bernsteingelb bis dunkelblau.', 37.00, true, '{}', ARRAY['dessertwein'], 15);
