--
-- PostgreSQL database dump
--

-- Dumped from database version 14.18 (Ubuntu 14.18-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.18 (Ubuntu 14.18-0ubuntu0.22.04.1)

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
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: rideshare_user
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO rideshare_user;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: rideshare_user
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO rideshare_user;

--
-- Name: trip_passengers; Type: TABLE; Schema: public; Owner: rideshare_user
--

CREATE TABLE public.trip_passengers (
    id integer NOT NULL,
    trip_id integer,
    passenger_id integer,
    status character varying(20) DEFAULT 'pending'::character varying,
    joined_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT trip_passengers_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'confirmed'::character varying, 'cancelled'::character varying])::text[])))
);


ALTER TABLE public.trip_passengers OWNER TO rideshare_user;

--
-- Name: trip_passengers_id_seq; Type: SEQUENCE; Schema: public; Owner: rideshare_user
--

CREATE SEQUENCE public.trip_passengers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.trip_passengers_id_seq OWNER TO rideshare_user;

--
-- Name: trip_passengers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: rideshare_user
--

ALTER SEQUENCE public.trip_passengers_id_seq OWNED BY public.trip_passengers.id;


--
-- Name: trips; Type: TABLE; Schema: public; Owner: rideshare_user
--

CREATE TABLE public.trips (
    id integer NOT NULL,
    driver_id integer,
    from_location character varying(255) NOT NULL,
    to_location character varying(255) NOT NULL,
    departure_time timestamp without time zone NOT NULL,
    max_passengers integer NOT NULL,
    current_passengers integer DEFAULT 0,
    price_per_person numeric(10,2) NOT NULL,
    description text,
    status character varying(20) DEFAULT 'active'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT trips_max_passengers_check CHECK ((max_passengers > 0)),
    CONSTRAINT trips_price_per_person_check CHECK ((price_per_person >= (0)::numeric)),
    CONSTRAINT trips_status_check CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'completed'::character varying, 'cancelled'::character varying])::text[])))
);


ALTER TABLE public.trips OWNER TO rideshare_user;

--
-- Name: trips_id_seq; Type: SEQUENCE; Schema: public; Owner: rideshare_user
--

CREATE SEQUENCE public.trips_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.trips_id_seq OWNER TO rideshare_user;

--
-- Name: trips_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: rideshare_user
--

ALTER SEQUENCE public.trips_id_seq OWNED BY public.trips.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: rideshare_user
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    phone character varying(20) NOT NULL,
    profile_image character varying(500),
    is_verified boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.users OWNER TO rideshare_user;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: rideshare_user
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO rideshare_user;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: rideshare_user
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: trip_passengers id; Type: DEFAULT; Schema: public; Owner: rideshare_user
--

ALTER TABLE ONLY public.trip_passengers ALTER COLUMN id SET DEFAULT nextval('public.trip_passengers_id_seq'::regclass);


--
-- Name: trips id; Type: DEFAULT; Schema: public; Owner: rideshare_user
--

ALTER TABLE ONLY public.trips ALTER COLUMN id SET DEFAULT nextval('public.trips_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: rideshare_user
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: rideshare_user
--

COPY public.schema_migrations (version, dirty) FROM stdin;
3	f
\.


--
-- Data for Name: trip_passengers; Type: TABLE DATA; Schema: public; Owner: rideshare_user
--

COPY public.trip_passengers (id, trip_id, passenger_id, status, joined_at) FROM stdin;
1	1	1	confirmed	2025-06-18 18:35:58.422765
2	1	2	confirmed	2025-06-18 18:43:33.498504
3	1	3	confirmed	2025-06-18 23:45:06.994225
4	4	1	confirmed	2025-06-21 13:35:16.878372
\.


--
-- Data for Name: trips; Type: TABLE DATA; Schema: public; Owner: rideshare_user
--

COPY public.trips (id, driver_id, from_location, to_location, departure_time, max_passengers, current_passengers, price_per_person, description, status, created_at, updated_at) FROM stdin;
1	1	Hall 13 , IIT Kanpur	Kanpur Central	2025-06-20 12:27:00	4	3	450.00	\N	active	2025-06-18 17:57:53.613932	2025-06-18 23:45:06.994225
2	1	hall 12,iit kanpur	kanpur central	2025-06-20 08:39:00	3	0	450.00	\N	cancelled	2025-06-19 14:09:27.432858	2025-06-19 22:29:19.22784
3	1	Hall 13 , IIT Kanpur	Kanpur Central	2025-06-21 14:45:00	3	0	500.00	\N	cancelled	2025-06-19 20:15:43.4127	2025-06-19 22:29:34.888288
4	2	Hall 13 , IIT Kanpur	Kanpur Central	2025-06-22 06:06:00	3	1	300.00	\N	active	2025-06-21 11:36:23.489311	2025-06-21 13:35:16.878372
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: rideshare_user
--

COPY public.users (id, name, email, password, phone, profile_image, is_verified, created_at, updated_at) FROM stdin;
1	Vivek	vivekwarkad017@gmail.com	$2a$10$pOY.QPb./HVmhb.Y1Mr7BOamFnikv3I8P3lNxDnmFJOu20kaRxVTa	9226283331	\N	f	2025-06-17 22:28:11.726317	2025-06-17 22:28:11.726317
2	Vivek	vivekwarkad0@gmail.com	$2a$10$nrDsUFZ9P34U24eTnGHaU.q4ThFItfElQIlYoKk1hZ.1M1/09lbjy	9226283331	\N	f	2025-06-18 17:58:29.29443	2025-06-18 17:58:29.29443
3	Vivek	vivekmw24@iitk.ac.in	$2a$10$u0S2bcNtd3lXTQUjRfTDgehH8RAGI6s3StHy/Ig3TUi8kFEhJkoT.	9226283331	\N	f	2025-06-18 23:44:48.583086	2025-06-18 23:44:48.583086
\.


--
-- Name: trip_passengers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: rideshare_user
--

SELECT pg_catalog.setval('public.trip_passengers_id_seq', 4, true);


--
-- Name: trips_id_seq; Type: SEQUENCE SET; Schema: public; Owner: rideshare_user
--

SELECT pg_catalog.setval('public.trips_id_seq', 4, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: rideshare_user
--

SELECT pg_catalog.setval('public.users_id_seq', 3, true);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: rideshare_user
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: trip_passengers trip_passengers_pkey; Type: CONSTRAINT; Schema: public; Owner: rideshare_user
--

ALTER TABLE ONLY public.trip_passengers
    ADD CONSTRAINT trip_passengers_pkey PRIMARY KEY (id);


--
-- Name: trip_passengers trip_passengers_trip_id_passenger_id_key; Type: CONSTRAINT; Schema: public; Owner: rideshare_user
--

ALTER TABLE ONLY public.trip_passengers
    ADD CONSTRAINT trip_passengers_trip_id_passenger_id_key UNIQUE (trip_id, passenger_id);


--
-- Name: trips trips_pkey; Type: CONSTRAINT; Schema: public; Owner: rideshare_user
--

ALTER TABLE ONLY public.trips
    ADD CONSTRAINT trips_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: rideshare_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: rideshare_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_trip_passengers_joined_at; Type: INDEX; Schema: public; Owner: rideshare_user
--

CREATE INDEX idx_trip_passengers_joined_at ON public.trip_passengers USING btree (joined_at);


--
-- Name: idx_trip_passengers_status; Type: INDEX; Schema: public; Owner: rideshare_user
--

CREATE INDEX idx_trip_passengers_status ON public.trip_passengers USING btree (status);


--
-- Name: idx_trip_passengers_trip; Type: INDEX; Schema: public; Owner: rideshare_user
--

CREATE INDEX idx_trip_passengers_trip ON public.trip_passengers USING btree (trip_id);


--
-- Name: idx_trip_passengers_user; Type: INDEX; Schema: public; Owner: rideshare_user
--

CREATE INDEX idx_trip_passengers_user ON public.trip_passengers USING btree (passenger_id);


--
-- Name: idx_trips_created_at; Type: INDEX; Schema: public; Owner: rideshare_user
--

CREATE INDEX idx_trips_created_at ON public.trips USING btree (created_at);


--
-- Name: idx_trips_departure; Type: INDEX; Schema: public; Owner: rideshare_user
--

CREATE INDEX idx_trips_departure ON public.trips USING btree (departure_time);


--
-- Name: idx_trips_driver; Type: INDEX; Schema: public; Owner: rideshare_user
--

CREATE INDEX idx_trips_driver ON public.trips USING btree (driver_id);


--
-- Name: idx_trips_locations; Type: INDEX; Schema: public; Owner: rideshare_user
--

CREATE INDEX idx_trips_locations ON public.trips USING btree (from_location, to_location);


--
-- Name: idx_trips_status; Type: INDEX; Schema: public; Owner: rideshare_user
--

CREATE INDEX idx_trips_status ON public.trips USING btree (status);


--
-- Name: idx_users_created_at; Type: INDEX; Schema: public; Owner: rideshare_user
--

CREATE INDEX idx_users_created_at ON public.users USING btree (created_at);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: rideshare_user
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: trips update_trips_updated_at; Type: TRIGGER; Schema: public; Owner: rideshare_user
--

CREATE TRIGGER update_trips_updated_at BEFORE UPDATE ON public.trips FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: users update_users_updated_at; Type: TRIGGER; Schema: public; Owner: rideshare_user
--

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: trip_passengers trip_passengers_passenger_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rideshare_user
--

ALTER TABLE ONLY public.trip_passengers
    ADD CONSTRAINT trip_passengers_passenger_id_fkey FOREIGN KEY (passenger_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: trip_passengers trip_passengers_trip_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rideshare_user
--

ALTER TABLE ONLY public.trip_passengers
    ADD CONSTRAINT trip_passengers_trip_id_fkey FOREIGN KEY (trip_id) REFERENCES public.trips(id) ON DELETE CASCADE;


--
-- Name: trips trips_driver_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: rideshare_user
--

ALTER TABLE ONLY public.trips
    ADD CONSTRAINT trips_driver_id_fkey FOREIGN KEY (driver_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

GRANT USAGE ON SCHEMA public TO vivek;


--
-- Name: FUNCTION update_updated_at_column(); Type: ACL; Schema: public; Owner: rideshare_user
--

GRANT ALL ON FUNCTION public.update_updated_at_column() TO vivek;


--
-- Name: TABLE schema_migrations; Type: ACL; Schema: public; Owner: rideshare_user
--

GRANT SELECT ON TABLE public.schema_migrations TO vivek;


--
-- Name: TABLE trip_passengers; Type: ACL; Schema: public; Owner: rideshare_user
--

GRANT SELECT ON TABLE public.trip_passengers TO vivek;


--
-- Name: SEQUENCE trip_passengers_id_seq; Type: ACL; Schema: public; Owner: rideshare_user
--

GRANT SELECT,USAGE ON SEQUENCE public.trip_passengers_id_seq TO vivek;


--
-- Name: TABLE trips; Type: ACL; Schema: public; Owner: rideshare_user
--

GRANT SELECT ON TABLE public.trips TO vivek;


--
-- Name: SEQUENCE trips_id_seq; Type: ACL; Schema: public; Owner: rideshare_user
--

GRANT SELECT,USAGE ON SEQUENCE public.trips_id_seq TO vivek;


--
-- Name: TABLE users; Type: ACL; Schema: public; Owner: rideshare_user
--

GRANT SELECT ON TABLE public.users TO vivek;


--
-- Name: SEQUENCE users_id_seq; Type: ACL; Schema: public; Owner: rideshare_user
--

GRANT SELECT,USAGE ON SEQUENCE public.users_id_seq TO vivek;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT SELECT,USAGE ON SEQUENCES  TO vivek;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT SELECT ON TABLES  TO vivek;


--
-- PostgreSQL database dump complete
--

