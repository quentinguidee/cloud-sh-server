-- region: SERVER
DROP TABLE IF EXISTS servers CASCADE;
-- endregion

-- region: AUTH
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS auth_github CASCADE;
-- endregion

-- region: STORAGE
DROP TABLE IF EXISTS buckets CASCADE;
DROP TABLE IF EXISTS buckets_nodes CASCADE;
DROP TABLE IF EXISTS buckets_to_node CASCADE;
DROP TABLE IF EXISTS buckets_access CASCADE;
DROP TABLE IF EXISTS buckets_nodes_associations CASCADE;
DROP TABLE IF EXISTS buckets_nodes_user_specific_data CASCADE;
-- endregion
