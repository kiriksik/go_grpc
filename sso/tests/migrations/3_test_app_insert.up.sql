INSERT INTO apps (id, name, secret)
VALUES (1, 'test', 'abobus')
ON CONFLICT DO NOTHING;