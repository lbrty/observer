# 🗄️ Initial Database Schema for DBA Expert

This document outlines the starting schema for the humanitarian tracking system. It adheres to the architectural layering and design conventions agreed across agents.

## 📚 Conventions

* Primary keys use **ULID**, generated in the application layer.
* All schemas use **SQLAlchemy Core**.
* No database-generated defaults unless explicitly necessary.
* Foreign keys must enforce `ON DELETE SET NULL` or `CASCADE` based on usage.
* `subjects` is a unified abstraction for `people` and `pets`, with a `type` field.

## 🧱 Core Tables

### `offices`

* `id`: ULID, primary key
* `name`: Text, required
* `place_id`: FK to `places.id`, nullable

### `users`

* `id`: ULID, primary key
* `email`: unique, required
* `full_name`: optional
* `password_hash`: required
* `is_active`: bool, default true
* `is_confirmed`: bool, default false
* `office_id`: FK to `offices.id`, nullable
* `created_at`, `updated_at`: timestamps

### `projects`

* `id`: ULID, primary key
* `name`: required
* `description`: optional
* `owner_id`: FK to `users.id`

### `categories`

* `id`: ULID, primary key
* `name`: unique, required

## 🧑‍🤝‍🧑 Subject Abstraction

### `subjects`

* `id`: ULID, primary key
* `type`: enum `'person' | 'pet'`
* `name`: required
* `status`: optional
* `email`, `phone_number`, `additional_phone`: optional
* `birth_date`: optional
* `sex`: enum `'male' | 'female' | 'unknown'`, optional
* `external_id`, `reference_id`: optional
* `tags`: array of text
* `project_id`, `category_id`, `consultant_id`, `office_id`: FK, nullable
* `created_at`, `updated_at`: timestamps

### `pets`

* `subject_id`: ULID, PK, FK to `subjects.id`
* `registration_id`: optional

> All pets are stored in `subjects` with `type='pet'`. Supplementary fields (e.g. `registration_id`) go into `pets`.

## 📁 Documents

### `documents`

* `id`: ULID, primary key
* `name`, `path`, `mimetype`: required
* `size`: integer
* `encryption_key`: optional
* `owner_id`: FK to `users.id`
* `project_id`: FK to `projects.id`
* `created_at`: timestamp

## 🔐 Permissions & Roles

### `user_roles`

* `user_id`: FK to `users.id`
* `role`: enum `'admin' | 'consultant' | 'staff' | 'guest'`
* PK: `(user_id, role)`

### `user_permissions`

* `user_id`: FK to `users.id`
* `project_id`: FK to `projects.id`
* `action`: enum `'create' | 'read' | 'update' | 'delete' | 'invite' | 'read_documents' | 'read_personal_info'`
* PK: `(user_id, project_id, action)`

### `user_tokens`

* `code`: primary key
* `user_id`: FK to `users.id`
* `type`: enum `'invite' | 'confirmation' | 'reset'`
* `expires_at`: timestamp

## 🧾 Audit Logs

### `audit_logs`

* `id`: ULID, primary key
* `ref`: Text, required
* `action`: optional
* `user_id`: FK to `users.id`
* `data`: JSONB
* `created_at`: timestamp

## 🧑‍⚕️ Support Records

### `support_records`

* `id`: ULID, primary key
* `description`: optional
* `type`: enum `'humanitarian' | 'legal' | 'medical' | 'general'`
* `subject_id`: FK to `subjects.id`
* `consultant_id`: FK to `users.id`
* `age_group`: optional
* `project_id`: FK to `projects.id`
* `created_at`: timestamp

## 🌍 Geography

### `countries`, `states`, `places`

* `places`: linked to `states` and `countries`
* `place_type`: enum `'city' | 'town' | 'village'`
* All place-related names and codes are indexed (lowercased)

## 🧳 Migration History

### `migration_history`

* `id`: ULID, primary key
* `subject_id`: FK to `subjects.id`
* `project_id`: FK to `projects.id`
* `from_place_id`, `current_place_id`: FK to `places.id`
* `migration_date`, `created_at`: timestamps

## ☑️ Notes

* All `created_at` and `updated_at` values are set in Python
* ULIDs are passed explicitly during inserts
* All constraints and indices must be reviewed by the DBA Expert
* Changes to this schema must go through the `migrations/` system and be coordinated with live data compatibility
