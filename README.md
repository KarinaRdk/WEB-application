# Go Web Application with PostgreSQL and Session Management

## Description
This Go-based web application demonstrates a simple yet effective approach to building a web server that handles HTTP requests, interacts with a PostgreSQL database, and manages user sessions. It includes functionality for user authentication, session management, and CRUD operations on blog entries, showcasing the use of GORM for database operations and Gorilla Sessions for session handling.

## Installation

Steps

1. Clone the repository to your local machine.
2. Navigate to the project directory.
3. Ensure Docker and Docker Compose are installed on your system.
4. Run the application with Docker Compose:

        docker-compose up --build
This command builds the necessary Docker images and starts the containers.
## Usage
Running the Application

After starting the containers with Docker Compose, the application will be accessible at http://localhost:8080 in your web browser.

## Features
User Authentication: Secure login and logout functionality for admin users. Should you wish to tesst it use the following credentials: login: Admin, password: hush-hush

Session Management: Utilizes Gorilla Sessions for session management, providing secure user sessions.

Database Operations: Leverages GORM for efficient PostgreSQL database operations, including Create and Read operations for blog entries.
