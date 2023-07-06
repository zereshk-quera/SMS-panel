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
---- Insert admin user and account if they don't exist
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
INSERT INTO public.users (
        id,
        first_name,
        last_name,
        phone,
        email,
        national_id
    )
VALUES (
        2,
        'Adel',
        'Mohammadzadeh',
        '09131234567',
        'adel.mohamadzadeph@gmail.com',
        '0611066998'
    );
--
-- Data for Name: accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.accounts (
        id,
        user_id,
        username,
        budget,
        password,
        token,
        is_active
    )
VALUES (
        2,
        2,
        'adel',
        440,
        '$2a$10$kMc4TRt0i1WCIdABCPAivuuV1SKY2G82HExrJJntcjKud5B/ZsjY.',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODc1OTY0NTUsImlkIjoxfQ.v4tmAt0PqQ9bYitkYMMZzBIFuH9adOhSy6aGU1HA3TA',
        true
    );
--
-- Data for Name: budget; Type: TABLE DATA; Schema: public; Owner: postgres
--

--
-- Data for Name: configuration; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.configuration (id, name, value)
VALUES (1, 'group sms', 10);
--
-- Data for Name: phone_books; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.phone_books (id, account_id, name)
VALUES (1, 1, 'phonebook_1');
--
-- Data for Name: phone_book_numbers; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.phone_book_numbers (id, phone_book_id, prefix, name, phone, username)
VALUES (1, 1, NULL, 'ali', '09191234567', 'ali');
--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.schema_migrations (version, dirty)
VALUES (10, false);
--
-- Data for Name: sender_numbers; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.sender_numbers (id, number, is_exclusive, is_default)
VALUES (1, '09121234567', false, true);
INSERT INTO public.sender_numbers (id, number, is_exclusive, is_default)
VALUES (2, '09141234567', false, true);
INSERT INTO public.sender_numbers (id, number, is_exclusive, is_default)
VALUES (3, '09151234567', true, false);
INSERT INTO public.sender_numbers (id, number, is_exclusive, is_default)
VALUES (4, '09161234567', true, false);
--
-- Data for Name: sms_messages; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.sms_messages (
        id,
        sender,
        recipient,
        message,
        schedule,
        delivery_report,
        created_at,
        account_id
    )
VALUES (
        7,
        '09151234567',
        '09191234567',
        'Yooo',
        NULL,
        'Message sent field',
        '2023-06-24 09:05:49.047283',
        1
    );
INSERT INTO public.sms_messages (
        id,
        sender,
        recipient,
        message,
        schedule,
        delivery_report,
        created_at,
        account_id
    )
VALUES (
        8,
        '09141234567',
        '09191234567',
        'Yooo',
        NULL,
        'Message sent field',
        '2023-06-24 09:16:47.999489',
        1
    );
INSERT INTO public.sms_messages (
        id,
        sender,
        recipient,
        message,
        schedule,
        delivery_report,
        created_at,
        account_id
    )
VALUES (
        9,
        '09141234567',
        '09191234567',
        'Yooo',
        NULL,
        'Message sent field',
        '2023-06-24 09:17:27.13585',
        1
    );
INSERT INTO public.sms_messages (
        id,
        sender,
        recipient,
        message,
        schedule,
        delivery_report,
        created_at,
        account_id
    )
VALUES (
        10,
        '09141234567',
        '09191234567',
        'Yooo',
        NULL,
        'Message sent successfully',
        '2023-06-24 09:17:42.347483',
        1
    );
INSERT INTO public.sms_messages (
        id,
        sender,
        recipient,
        message,
        schedule,
        delivery_report,
        created_at,
        account_id
    )
VALUES (
        11,
        '09151234567',
        '09191234567',
        'Yooo',
        NULL,
        'Message sent successfully',
        '2023-06-24 09:26:17.325593',
        1
    );
INSERT INTO public.sms_messages (
        id,
        sender,
        recipient,
        message,
        schedule,
        delivery_report,
        created_at,
        account_id
    )
VALUES (
        12,
        '09151234567',
        '09191234567',
        'Yooo',
        NULL,
        'Message sent successfully',
        '2023-06-24 09:42:08.563243',
        1
    );
INSERT INTO public.sms_messages (
        id,
        sender,
        recipient,
        message,
        schedule,
        delivery_report,
        created_at,
        account_id
    )
VALUES (
        13,
        '09141234567',
        '09191234567',
        'Yooo',
        NULL,
        'Message sent successfully',
        '2023-06-24 09:42:28.548462',
        1
    );
INSERT INTO public.sms_messages (
        id,
        sender,
        recipient,
        message,
        schedule,
        delivery_report,
        created_at,
        account_id
    )
VALUES (
        14,
        '09151234567',
        '09131234567',
        'Hello, World!',
        NULL,
        'Message sent successfully',
        '2023-06-24 11:20:49.579756',
        1
    );
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
        is_available,
        subscription_package_id
    )
VALUES (1, 1, 3, '2023-06-23', '2023-06-23', true, NULL);
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

SELECT pg_catalog.setval('public.configuration_id_seq', 1, true);
--
-- Name: phone_book_numbers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.phone_book_numbers_id_seq', 1, true);
--
-- Name: phone_books_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.phone_books_id_seq', 1, true);
--
-- Name: sender_numbers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sender_numbers_id_seq', 4, true);
--
-- Name: sms_messages_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sms_messages_id_seq', 14, true);
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