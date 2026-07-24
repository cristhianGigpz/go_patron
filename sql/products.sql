--
-- PostgreSQL database dump
--

-- Dumped from database version 16.9
-- Dumped by pg_dump version 16.9

-- Started on 2026-07-24 11:42:48

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 216 (class 1259 OID 42056)
-- Name: products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.products (
    id bigint NOT NULL,
    name text,
    description text,
    price numeric,
    stock bigint
);


ALTER TABLE public.products OWNER TO postgres;

--
-- TOC entry 215 (class 1259 OID 42055)
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.products_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.products_id_seq OWNER TO postgres;

--
-- TOC entry 4787 (class 0 OID 0)
-- Dependencies: 215
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- TOC entry 4634 (class 2604 OID 42059)
-- Name: products id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- TOC entry 4781 (class 0 OID 42056)
-- Dependencies: 216
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (id, name, description, price, stock) FROM stdin;
1	Laptop Lenovo IdeaPad	Laptop para trabajo y estudio	2499.9	12
2	Mouse Logitech M185	Mouse inalámbrico	59.9	35
3	Teclado Mecánico Redragon	Teclado mecánico RGB	189.9	18
4	Monitor LG 24 pulgadas	Monitor Full HD	699.9	10
5	Audífonos Sony WH-CH520	Audífonos Bluetooth	179.9	22
6	Webcam Logitech C920	Cámara Full HD para videollamadas	329.9	8
7	Disco SSD Kingston 480GB	Unidad de almacenamiento SSD	159.9	25
8	Memoria USB Kingston 64GB	Memoria USB 3.0	29.9	50
9	Router TP-Link Archer C6	Router Wi-Fi de doble banda	219.9	14
10	Cable HDMI 2 metros	Cable HDMI de alta velocidad	24.9	40
11	Tablet Samsung Galaxy Tab A9	Tablet de 8.7 pulgadas	749.9	7
12	Smartphone Motorola G54	Teléfono móvil 5G	899.9	9
13	Cargador USB-C 25W	Cargador de carga rápida	69.9	30
14	Power Bank Xiaomi 10000mAh	Batería portátil	99.9	16
15	Parlante JBL Go 3	Parlante Bluetooth portátil	159.9	11
16	Silla Ergonómica	Silla para oficina	549.9	5
17	Escritorio 120cm	Escritorio de melamina	399.9	6
18	Mochila para Laptop	Mochila impermeable	119.9	20
19	Impresora Multifuncional Epson	Impresora con conexión Wi-Fi	629.9	4
20	Regulador de Voltaje	Regulador de voltaje de 1000VA	129.9	13
\.


--
-- TOC entry 4788 (class 0 OID 0)
-- Dependencies: 215
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.products_id_seq', 20, true);


--
-- TOC entry 4636 (class 2606 OID 42063)
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


-- Completed on 2026-07-24 11:42:48

--
-- PostgreSQL database dump complete
--

