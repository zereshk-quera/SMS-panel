--
-- PostgreSQL database dump
--

-- Dumped from database version 14.8 (Ubuntu 14.8-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.8 (Ubuntu 14.8-0ubuntu0.22.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.users (id, first_name, last_name, phone, email, national_id) VALUES (1, 'Adel', 'Mohammadzadeh', '09131234567', 'adel.mohamadzadeph@gmail.com', '0611066998');


--
-- Data for Name: accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.accounts (id, user_id, username, budget, password, token, is_active) VALUES (1, 1, 'adel', 0, '$2a$10$kMc4TRt0i1WCIdABCPAivuuV1SKY2G82HExrJJntcjKud5B/ZsjY.', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODc1MjUzNjIsImlkIjoxfQ.WvXK4ugzTEn6fQhaKRgW7jiy9irpSqkwhNh11X2Tg90', true);


--
-- Data for Name: budget; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: configuration; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: phone_books; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: phone_book_numbers; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.schema_migrations (version, dirty) VALUES (10, false);


--
-- Data for Name: sender_numbers; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.sender_numbers (id, number, is_exclusive, is_default) VALUES (1, '09121234567', false, true);
INSERT INTO public.sender_numbers (id, number, is_exclusive, is_default) VALUES (2, '09141234567', false, true);
INSERT INTO public.sender_numbers (id, number, is_exclusive, is_default) VALUES (3, '09151234567', true, false);


--
-- Data for Name: sms_messages; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: user_numbers; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.user_numbers (id, user_id, number_id, start_date, end_date, is_available) VALUES (1, 1, 3, '2023-06-23', '2023-06-23', true);


--
-- Name: accounts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.accounts_id_seq', 1, true);


--
-- Name: budget_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.budget_id_seq', 1, false);


--
-- Name: configuration_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.configuration_id_seq', 1, false);


--
-- Name: phone_book_numbers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.phone_book_numbers_id_seq', 1, false);


--
-- Name: phone_books_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.phone_books_id_seq', 1, false);


--
-- Name: sender_numbers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sender_numbers_id_seq', 3, true);


--
-- Name: sms_messages_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sms_messages_id_seq', 1, false);


--
-- Name: transactions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.transactions_id_seq', 1, false);


--
-- Name: user_numbers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_numbers_id_seq', 1, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 1, true);


--
-- PostgreSQL database dump complete
--

