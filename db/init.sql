--
-- PostgreSQL database dump
--

-- Dumped from database version 15.4
-- Dumped by pg_dump version 15.4

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
-- Name: segments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.segments (
    id integer NOT NULL,
    slug text NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE public.segments OWNER TO postgres;

--
-- Name: segments_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.segments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.segments_id_seq OWNER TO postgres;

--
-- Name: segments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.segments_id_seq OWNED BY public.segments.id;


--
-- Name: user_segment_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_segment_history (
    user_id integer NOT NULL,
    segment_slug text NOT NULL,
    date_added timestamp without time zone NOT NULL,
    date_removed timestamp without time zone
);


ALTER TABLE public.user_segment_history OWNER TO postgres;

--
-- Name: users_segments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users_segments (
    user_id integer NOT NULL,
    slug text NOT NULL,
    expiration_date date
);


ALTER TABLE public.users_segments OWNER TO postgres;

--
-- Name: segments id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.segments ALTER COLUMN id SET DEFAULT nextval('public.segments_id_seq'::regclass);


--
-- Name: segments segments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.segments
    ADD CONSTRAINT segments_pkey PRIMARY KEY (id);


--
-- Name: segments segments_slug_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.segments
    ADD CONSTRAINT segments_slug_key UNIQUE (slug);


--
-- Name: user_segment_history user_segment_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_segment_history
    ADD CONSTRAINT user_segment_history_pkey PRIMARY KEY (user_id, segment_slug, date_added);


--
-- Name: users_segments users_segments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_segments
    ADD CONSTRAINT users_segments_pkey PRIMARY KEY (user_id, slug);


--
-- PostgreSQL database dump complete
--

