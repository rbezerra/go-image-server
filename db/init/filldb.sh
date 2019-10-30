#!/bin/bash

set -o errexit

readonly REQUIRED_ENV_VARS=(
    "DB_USER"
    "DB_PASSWORD"
    "DB_DATABASE"
    "POSTGRES_USER"
)

main(){
    check_env_vars_set
    init_user_and_db
}

check_env_vars_set(){
    for required_env_var in ${REQUIRED_ENV_VARS[@]}; do
        if [[ -z "${!required_env_var}" ]]; then
            echo "Error:
    Environment variable '$required_env_var' not set.
    Make sure you have the following environment variables set:

              ${REQUIRED_ENV_VARS[@]}

        Aborting."
            exit 1
        fi
    done

}

init_user_and_db(){
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
        CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';
        CREATE DATABASE $DB_DATABASE;
        GRANT ALL PRIVILEGES ON DATABASE $DB_DATABASE TO $DB_USER;
        ALTER USER $DB_USER WITH SUPERUSER;

        \c $DB_DATABASE;

        CREATE SEQUENCE imagem_id_seq START 1;
        CREATE SEQUENCE arquivo_id_seq START 1;


        CREATE TABLE IF NOT EXISTS public.imagem
        (
            id bigint NOT NULL DEFAULT nextval('imagem_id_seq'::regclass),
            uuid text COLLATE pg_catalog."default" NOT NULL,
            descricao text COLLATE pg_catalog."default" NOT NULL,
            CONSTRAINT imagem_pkey PRIMARY KEY (id)
        )
        WITH (
            OIDS = FALSE
        )
        TABLESPACE pg_default;

        ALTER TABLE public.imagem
            OWNER to $DB_USER;

        CREATE TABLE IF NOT EXISTS public.arquivo
        (
            id bigint NOT NULL DEFAULT nextval('arquivo_id_seq'::regclass),
            imagem_id bigint NOT NULL,
            tamanho character varying(10) COLLATE pg_catalog."default" NOT NULL,
            path text COLLATE pg_catalog."default" NOT NULL,
            original boolean NOT NULL,
            CONSTRAINT arquivo_pkey PRIMARY KEY (id),
            CONSTRAINT fk_imagem FOREIGN KEY (imagem_id)
                REFERENCES public.imagem (id) MATCH SIMPLE
                ON UPDATE CASCADE
                ON DELETE CASCADE
        )
        WITH (
            OIDS = FALSE
        )
        TABLESPACE pg_default;

        ALTER TABLE public.arquivo
            OWNER to $DB_USER;

 
EOSQL
}

main "$@"

