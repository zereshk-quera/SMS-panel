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
-- Insert admin user and account if they don't exist
INSERT INTO users (first_name, last_name, phone, email, national_id)
SELECT 'admin',
    'admin',
    'admin',
    'admin',
    'admin'
WHERE NOT EXISTS (
        SELECT 1
        FROM users
        WHERE email = 'admin'
    );
INSERT INTO accounts (
        user_id,
        budget,
        is_active,
        is_admin,
        username,
        password,
        token
    )
SELECT (
        SELECT id
        FROM users
        WHERE email = 'admin'
    ),
    0,
    true,
    true,
    'admin',
    '$2a$10$neJk2ozYYr19o5GseVDYieOu1QHubXeQ08F8dN8eoX2GiDRz09xKi',
    -- hashed password 'admin123'
    'token' -- admin should get new token for login
WHERE NOT EXISTS (
        SELECT 1
        FROM accounts
        WHERE username = 'admin'
    );
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
--
-- Data for Name: sender_numbers; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.sender_numbers (id, number, is_exclusive, is_default)
VALUES (1, '09121234567', false, true);
INSERT INTO public.sender_numbers (id, number, is_exclusive, is_default)
VALUES (2, '09141234567', false, true);
INSERT INTO public.sender_numbers (id, number, is_exclusive, is_default)
VALUES (3, '09151234567', true, false);
--
-- Data for Name: sms_messages; Type: TABLE DATA; Schema: public; Owner: postgres
--

--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: postgres
--

--
-- Data for Name: user_numbers; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.user_numbers (
        id,
        user_id,
        number_id,
        start_date,
        end_date,
        is_available
    )
VALUES (1, 1, 3, '2023-06-23', '2023-06-23', true);
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
- -