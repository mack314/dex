package migrations

// THIS FILE WAS GENERATED BY gen.go DO NOT EDIT

import "github.com/rubenv/sql-migrate"

var PostgresMigrations migrate.MigrationSource = &migrate.MemoryMigrationSource{
	Migrations: []*migrate.Migration{
		{
			Id: "0001_initial_migration.sql",
			Up: []string{
				"-- +migrate Up\nCREATE TABLE IF NOT EXISTS \"authd_user\" (\n       \"id\" text not null primary key,\n       \"email\" text,\n       \"email_verified\" boolean,\n       \"display_name\" text,\n       \"admin\" boolean) ;\n\nCREATE TABLE IF NOT EXISTS \"client_identity\" (\n       \"id\" text not null primary key,\n       \"secret\" bytea,\n       \"metadata\" text);\n\nCREATE TABLE IF NOT EXISTS \"connector_config\" (\n       \"id\" text not null primary key,\n       \"type\" text, \"config\" text) ;\n\nCREATE TABLE IF NOT EXISTS \"key\" (\n       \"value\" bytea not null primary key) ;\n\nCREATE TABLE IF NOT EXISTS \"password_info\" (\n       \"user_id\" text not null primary key,\n       \"password\" text,\n       \"password_expires\" bigint) ;\n\nCREATE TABLE IF NOT EXISTS \"session\" (\n       \"id\" text not null primary key,\n       \"state\" text,\n       \"created_at\" bigint,\n       \"expires_at\" bigint,\n       \"client_id\" text,\n       \"client_state\" text,\n       \"redirect_url\" text, \"identity\" text,\n       \"connector_id\" text,\n       \"user_id\" text, \"register\" boolean) ;\n\nCREATE TABLE IF NOT EXISTS \"session_key\" (\n       \"key\" text not null primary key,\n       \"session_id\" text,\n       \"expires_at\" bigint,\n       \"stale\" boolean) ;\n\nCREATE TABLE IF NOT EXISTS \"remote_identity_mapping\" (\n       \"connector_id\" text not null,\n       \"user_id\" text,\n       \"remote_id\" text not null,\n       primary key (\"connector_id\", \"remote_id\")) ;\n",
			},
		},
		{
			Id: "0002_dex_admin.sql",
			Up: []string{
				"-- +migrate Up\nALTER TABLE client_identity ADD COLUMN \"dex_admin\" boolean;\n\nUPDATE \"client_identity\" SET \"dex_admin\" = false;\n",
			},
		},
		{
			Id: "0003_user_created_at.sql",
			Up: []string{
				"-- +migrate Up\nALTER TABLE authd_user ADD COLUMN \"created_at\" bigint;\n\nUPDATE authd_user SET \"created_at\" = 0;\n",
			},
		},
		{
			Id: "0004_session_nonce.sql",
			Up: []string{
				"-- +migrate Up\nALTER TABLE session ADD COLUMN \"nonce\" text;\n",
			},
		},
		{
			Id: "0005_refresh_token_create.sql",
			Up: []string{
				"-- +migrate Up\nCREATE TABLE refresh_token (\n    id bigint NOT NULL,\n    payload_hash bytea,\n    user_id text,\n    client_id text\n);\n\nCREATE SEQUENCE refresh_token_id_seq\n    START WITH 1\n    INCREMENT BY 1\n    NO MINVALUE\n    NO MAXVALUE\n    CACHE 1;\n\nALTER SEQUENCE refresh_token_id_seq OWNED BY refresh_token.id;\n\nALTER TABLE ONLY refresh_token ALTER COLUMN id SET DEFAULT nextval('refresh_token_id_seq'::regclass);\n\nALTER TABLE ONLY refresh_token\n    ADD CONSTRAINT refresh_token_pkey PRIMARY KEY (id);\n",
			},
		},
		{
			Id: "0006_user_email_unique.sql",
			Up: []string{
				"-- +migrate Up\nALTER TABLE ONLY authd_user\n    ADD CONSTRAINT authd_user_email_key UNIQUE (email);\n",
			},
		},
		{
			Id: "0007_session_scope.sql",
			Up: []string{
				"-- +migrate Up\nALTER TABLE session ADD COLUMN \"scope\" text;\n",
			},
		},
		{
			Id: "0008_users_active_or_inactive.sql",
			Up: []string{
				"-- +migrate Up\nALTER TABLE authd_user ADD COLUMN disabled boolean;\n\nUPDATE authd_user SET \"disabled\" = FALSE;\n",
			},
		},
		{
			Id: "0009_key_not_primary_key.sql",
			Up: []string{
				"-- +migrate Up\nALTER TABLE key ADD COLUMN tmp_value bytea;\nUPDATE KEY SET tmp_value = value;\nALTER TABLE key DROP COLUMN value;\nALTER TABLE key RENAME COLUMN \"tmp_value\" to \"value\";\n",
			},
		},
		{
			Id: "0010_client_metadata_field_changed.sql",
			Up: []string{
				"-- +migrate Up\nUPDATE client_identity\nSET metadata = text(\n    json_build_object(\n        'redirectURLs', json(json(metadata)->>'redirectURLs'),\n        'redirect_uris', json(json(metadata)->>'redirectURLs')\n    )\n )\nWHERE (json(metadata)->>'redirect_uris') IS NULL;\n",
			},
		},
		{
			Id: "0011_case_insensitive_emails.sql",
			Up: []string{
				"-- +migrate Up\n\n-- This migration is a fix for a bug that allowed duplicate emails if they used different cases (see #338).\n-- When migrating, dex will not take the liberty of deleting rows for duplicate cases. Instead it will\n-- raise an exception and call for an admin to remove duplicates manually.\n\nCREATE OR REPLACE FUNCTION raise_exp() RETURNS VOID AS $$\nBEGIN\n     RAISE EXCEPTION 'Found duplicate emails when using case insensitive comparision, cannot perform migration.';\nEND;\n$$ LANGUAGE plpgsql;\n\nSELECT LOWER(email),\n    COUNT(email),\n    CASE\n        WHEN COUNT(email) > 1 THEN raise_exp()\n        ELSE NULL\n    END\nFROM authd_user\nGROUP BY LOWER(email);\n\nUPDATE authd_user SET email = LOWER(email);\n",
			},
		},
		{
			Id: "0012_add_cross_client_authorizers.sql",
			Up: []string{
				"-- +migrate Up\nCREATE TABLE IF NOT EXISTS \"trusted_peers\" (\n       \"client_id\" text not null,\n       \"trusted_client_id\" text not null,\n       primary key (\"client_id\", \"trusted_client_id\")) ;\n",
			},
		},
		{
			Id: "0013_add_public_clients.sql",
			Up: []string{
				"-- +migrate Up\nALTER TABLE client_identity ADD COLUMN \"public\" boolean;\n\nUPDATE \"client_identity\" SET \"public\" = false;\n",
			},
		},
		{
			Id: "0013_add_scopes_to_refresh_tokens.sql",
			Up: []string{
				"-- +migrate Up\nALTER TABLE refresh_token ADD COLUMN \"scopes\" text;\n\nUPDATE refresh_token SET scopes = 'openid profile email offline_access';\n",
			},
		},
	},
}
