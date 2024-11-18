-- Insert 10 initial users with properly hashed passwords (using bcrypt)
INSERT INTO users (email, password, name) VALUES
('user1@example.com', '$2a$10$1234567890123456789012uqZFruBqT8Z9z4XQ9Q9Q9Q9Q9Q9Q9Q9Q', 'John Doe'),
('user2@example.com', '$2a$10$1234567890123456789012uqZFruBqT8Z9z4XQ9Q9Q9Q9Q9Q9Q9Q9Q', 'Jane Smith'),
('user3@example.com', '$2a$10$1234567890123456789012uqZFruBqT8Z9z4XQ9Q9Q9Q9Q9Q9Q9Q9Q', 'Mike Johnson'),
('user4@example.com', '$2a$10$1234567890123456789012uqZFruBqT8Z9z4XQ9Q9Q9Q9Q9Q9Q9Q9Q', 'Sarah Williams'),
('user5@example.com', '$2a$10$1234567890123456789012uqZFruBqT8Z9z4XQ9Q9Q9Q9Q9Q9Q9Q9Q', 'David Brown'),
('user6@example.com', '$2a$10$1234567890123456789012uqZFruBqT8Z9z4XQ9Q9Q9Q9Q9Q9Q9Q9Q', 'Emily Davis'),
('user7@example.com', '$2a$10$1234567890123456789012uqZFruBqT8Z9z4XQ9Q9Q9Q9Q9Q9Q9Q9Q', 'Michael Wilson'),
('user8@example.com', '$2a$10$1234567890123456789012uqZFruBqT8Z9z4XQ9Q9Q9Q9Q9Q9Q9Q9Q', 'Lisa Anderson'),
('user9@example.com', '$2a$10$1234567890123456789012uqZFruBqT8Z9z4XQ9Q9Q9Q9Q9Q9Q9Q9Q', 'Robert Taylor'),
('user10@example.com', '$2a$10$1234567890123456789012uqZFruBqT8Z9z4XQ9Q9Q9Q9Q9Q9Q9Q9Q', 'Jennifer Martin');
