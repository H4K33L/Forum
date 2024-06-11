# Forum

As part of our final year project for the 2023-2024 bachelor's degree in computer science, our mission is to put into practice all the skills and knowledge acquired during our first year.

This project, of utmost importance, involves creating a forum website that integrates all the functionalities, tools, and languages we have studied. It aims to demonstrate our ability to design and develop a complete web application, applying best practices in software development, databases, security, and project management.

This specification will detail the objectives, functional and technical requirements, as well as the tools and technologies used to carry out this ambitious project.

## Project Overview

This project consists of creating a web forum that allows:

- Communication between users
- Association of categories with messages
- Appreciation (likes) and disapproval (dislikes) of messages and comments
- Filtering of messages

## SQLite

To store our forum's data (such as users, messages, comments, etc.), we will use the SQLite database library.

SQLite is a popular choice as embedded database software for local/client storage in application software such as web browsers. It allows the creation of a database and its management using queries.

To structure our database and improve performance, we will consult the entity-relationship diagram and create one based on our own database.

We must use at least one `SELECT` query, one `CREATE` query, and one `INSERT` query. For more information on SQLite, we can consult the [SQLite page](https://www.sqlite.org).

## Authentication

In this section, the client must be able to register as a new user on the forum by entering their credentials. We must also create a login session to access the forum and be able to add messages and comments.

We can consider filtering by categories as sub-forums. A sub-forum is a section of an online forum dedicated to a specific topic.

We must use cookies to allow each user to have only one open session. Each of these sessions must contain an expiration date. It is up to us to decide the cookie's lifespan. Using UUID is a bonus task.

### Instructions for User Registration:

- Request an email.
- If the email is already taken, return an error response.
- Request a username.
- Request a password.
- The password must be encrypted when stored (bonus task).
- The forum must check if the provided email is in the database and if all credentials are correct. It will verify if the password matches the provided one, and if not, return an error response.

## Communication

To allow users to communicate with each other, they must be able to create messages and comments.

- Only registered users can create messages and comments.
- When creating a message, users can associate it with one or more categories.
- The implementation and choice of categories are at your discretion.
- Messages and comments must be visible to all users (registered or not).
- Unregistered users can only view messages and comments.

## Likes and Dislikes

- Only registered users can like or dislike messages and comments.
- The number of likes and dislikes must be visible to all users (registered or not).

## Filter

We must implement a filtering mechanism that allows users to filter displayed messages by:

- Categories
- Created messages
- Liked messages

Note that the last two points are only available for registered users and must refer to the logged-in user.

## Instructions

- Use SQLite.
- Handle website errors and HTTP statuses.
- Handle all kinds of technical errors.
- The code must follow best practices.
- It is recommended to have test files for unit tests.

## Authorized Packages

- All standard Go packages are allowed.
- `sqlite3`
- `bcrypt`
- `UUID`

We must not use any frontend library or framework like React, Angular, Vue, etc.

## Learning Outcomes

This project will help us learn:

- Web basics: HTML, HTTP, Sessions, and Cookies
- SQL language
- Database manipulation
- Basics of encryption
