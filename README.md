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


# Screenshots: 

<img width="750" alt="Screenshot 2024-03-10 at 10 54 39 PM" src="https://github.com/KarinaRdk/WEB-application/assets/101858179/4507ac51-40b9-457d-9355-919f1cc8f721">

<img width="749" alt="Screenshot 2024-03-10 at 10 54 51 PM" src="https://github.com/KarinaRdk/WEB-application/assets/101858179/f3676233-1366-4b2c-9427-a8ff53a06045">

<img width="745" alt="Screenshot 2024-03-10 at 10 55 11 PM" src="https://github.com/KarinaRdk/WEB-application/assets/101858179/1d0bd57f-82db-49a8-a0e5-83b2807e229a">

<img width="833" alt="Screenshot 2024-03-10 at 10 48 13 PM" src="https://github.com/KarinaRdk/WEB-application/assets/101858179/d784942c-3be8-4753-b86d-809d1c8d7708">

<img width="810" alt="Screenshot 2024-03-10 at 10 47 40 PM" src="https://github.com/KarinaRdk/WEB-application/assets/101858179/252df704-3000-401a-afae-40f3e25e63cb">

