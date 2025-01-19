INSERT INTO public.users (created_at, updated_at, avatar, deleted_at, username, email, password, secret)
VALUES
    (NOW(), NOW(), 'https://example.com/avatar1.png', NULL, 'user1', 'user1@example.com', 'password1', 'secret1'),
    (NOW(), NOW(), 'https://example.com/avatar2.png', NULL, 'user2', 'user2@example.com', 'password2', 'secret2'),
    (NOW(), NOW(), 'https://example.com/avatar3.png', NULL, 'user3', 'user3@example.com', 'password3', 'secret3'),
    (NOW(), NOW(), 'https://example.com/avatar4.png', NULL, 'user4', 'user4@example.com', 'password4', 'secret4'),
    (NOW(), NOW(), 'https://example.com/avatar5.png', NULL, 'user5', 'user5@example.com', 'password5', 'secret5');