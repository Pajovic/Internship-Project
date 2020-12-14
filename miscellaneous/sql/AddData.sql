-- Kompanije
INSERT INTO public.companies
(id, "name", ismain)
VALUES('bb099b7e-03d1-11eb-adc1-0242ac120002', 'Kompanija C', false);
INSERT INTO public.companies
(id, "name", ismain)
VALUES('5f46ed8c-03d1-11eb-adc1-0242ac120002', 'Kompanija B', false);
INSERT INTO public.companies
(id, "name", ismain)
VALUES('1051258e-88e9-4168-a889-5c222eaa152c', 'Kompanija A', true);

-- Kompanija B
INSERT INTO public.employees
(id, firstname, lastname, idc, c, r, u, d)
VALUES('8a14994b-2438-46f0-8bca-604a39d92591', 'Pera', 'Peric', '5f46ed8c-03d1-11eb-adc1-0242ac120002', true, true, true, true);
INSERT INTO public.employees
(id, firstname, lastname, idc, c, r, u, d)
VALUES('8397075f-06e7-475e-931f-e9bd49647c6d', 'Mika', 'Mikic', '5f46ed8c-03d1-11eb-adc1-0242ac120002', false, true, false, false);
INSERT INTO public.employees
(id, firstname, lastname, idc, c, r, u, d)
VALUES('c4b8f7b2-0a1f-4bbc-a525-8bb93257b63f', 'Zora', 'Zoric', '5f46ed8c-03d1-11eb-adc1-0242ac120002', false, true, false, false);

-- Proizvodi kompanije B
INSERT INTO public.products
(id, "name", price, quantity, idc)
VALUES('8d0cd2e1-30c3-4a9e-a862-f04bc764a4f9', 'Elektromotor', 550, 5, '5f46ed8c-03d1-11eb-adc1-0242ac120002');
INSERT INTO public.products
(id, "name", price, quantity, idc)
VALUES('6267b35e-fb6e-4380-8dd6-1b4ec3b8cec7', 'Energetski kabel', 100, 8, '5f46ed8c-03d1-11eb-adc1-0242ac120002');
INSERT INTO public.products
(id, "name", price, quantity, idc)
VALUES('5097a3b3-626e-479a-a501-00077be2b816', 'Osigurac', 20, 120, '5f46ed8c-03d1-11eb-adc1-0242ac120002');
INSERT INTO public.products
(id, "name", price, quantity, idc)
VALUES('e50fe62d-2650-4fc9-a65e-2eb70a692901', 'Kuciste', 200, 20, '5f46ed8c-03d1-11eb-adc1-0242ac120002');

-- Kompanija C
INSERT INTO public.employees
(id, firstname, lastname, idc, c, r, u, d)
VALUES('c1cbc0e1-b2e1-4321-adc6-2bef48e516df', 'Sloba', 'Stankovic', 'bb099b7e-03d1-11eb-adc1-0242ac120002', true, true, true, true);
INSERT INTO public.employees
(id, firstname, lastname, idc, c, r, u, d)
VALUES('5921c5ef-818b-4b93-bd57-f85e9c8f3df2', 'Jovana', 'Jovanovic', 'bb099b7e-03d1-11eb-adc1-0242ac120002', true, true, true, false);
INSERT INTO public.employees
(id, firstname, lastname, idc, c, r, u, d)
VALUES('1e2371c9-0eca-49ee-8072-4fd8cf2db22a', 'Nemanja', 'Nemanjic', 'bb099b7e-03d1-11eb-adc1-0242ac120002', false, true, false, false);

-- Proizvodi kompanije C
INSERT INTO public.products
(id, "name", price, quantity, idc)
VALUES('a347930a-fc20-491b-9709-5a5184be64bc', 'Ves masina', 2000, 5, 'bb099b7e-03d1-11eb-adc1-0242ac120002');
INSERT INTO public.products
(id, "name", price, quantity, idc)
VALUES('a9ed5b74-4778-4a6f-ad06-dd865b029406', 'Bojler', 500, 4, 'bb099b7e-03d1-11eb-adc1-0242ac120002');
INSERT INTO public.products
(id, "name", price, quantity, idc)
VALUES('1a2705ab-ffd2-4390-9c6e-445c017904d3', 'Sporet', 1200, 2, 'bb099b7e-03d1-11eb-adc1-0242ac120002');
INSERT INTO public.products
(id, "name", price, quantity, idc)
VALUES('e9fdda54-14bc-4580-855c-c2a77a93c2b3', 'Sudo masina', 800, 7, 'bb099b7e-03d1-11eb-adc1-0242ac120002');

-- Kompanija B dijeli podatke o proizvodima sa kompanijom C
INSERT INTO public.external_access_rights
(id, idsc, idrc, r, u, d, approved)
VALUES('fa71c166-980f-4a17-aa7c-b85df4be8989', '5f46ed8c-03d1-11eb-adc1-0242ac120002', 'bb099b7e-03d1-11eb-adc1-0242ac120002', true, true, true, true);
INSERT INTO public.external_access_rights
(id, idsc, idrc, r, u, d, approved)
VALUES('de7cc1b1-d858-4bd6-92cf-abf2274731ac', '5f46ed8c-03d1-11eb-adc1-0242ac120002', 'bb099b7e-03d1-11eb-adc1-0242ac120002', true, false, false, true);

-- Operatori
INSERT INTO public.operators
(id, "name")
VALUES(1, '>');
INSERT INTO public.operators
(id, "name")
VALUES(2, '>=');
INSERT INTO public.operators
(id, "name")
VALUES(3, '<');
INSERT INTO public.operators
(id, "name")
VALUES(4, '<=');

-- Propertiji
INSERT INTO public.properties
(id, "name")
VALUES(1, 'quantity');
INSERT INTO public.properties
(id, "name")
VALUES(2, 'price');

-- Ogranicenja
INSERT INTO public.access_constraints
(id, idear, operator_id, property_id, property_value)
VALUES('f87b85d8-9037-41a5-8d2b-6861cde17c18', 'fa71c166-980f-4a17-aa7c-b85df4be8989', 2, 1, 10);
INSERT INTO public.access_constraints
(id, idear, operator_id, property_id, property_value)
VALUES('8120ea1b-5823-4100-8bd5-80f9cb0db831', 'de7cc1b1-d858-4bd6-92cf-abf2274731ac', 3, 1, 10);