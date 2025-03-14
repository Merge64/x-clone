# X.com Clone

[[ CURRENTLY IN DEVELOPMENT BY THE END OF MARCH A STABLE V1 WILL BE RELEASED ]]

A clone of X.com built to enhance our understanding of web development. This project was created as part of our personal learning journey.

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Installation](#installation)
- [Environment Variables Setup](#environment-variables-setup)
- [Usage](#usage)
- [Legal Disclaimer](#legal-disclaimer)
- [How to use the app](#how-to-use-the-app)

## Overview
This is a fully functional clone of **X.com**, designed to simulate key features of the platform, including user registration, transactions, account management, and other e-commerce functionalities. The goal of this project was to improve my understanding of front-end and back-end integration, database management, and web application security.

![image](https://github.com/user-attachments/assets/92ac2324-fc7a-4495-9b1d-501f61408d69)
![image](https://github.com/user-attachments/assets/a1db35b3-d122-4af7-ad7d-e873f184f05f)
![image](https://github.com/user-attachments/assets/b5bd0704-00ae-47fa-84c1-3730e0c94baa)

## Features
- **User Authentication**: Registration, login, and password recovery.
- **User Interactions**: Posting, reposting, commenting, liking, and sending private messages (DMs).

## Technologies Used
- **Frontend**: HTML, CSS, JavaScript (React.js)
- **Backend**: Golang
- **Database**: Postgres
- **Authentication**: JWT (JSON Web Tokens) using Gin.
- **Version Control**: Git, GitHub

## Installation
To get a local copy up and running:

1. Clone the repository and set the required environment-variables as stated in [Environment Variables Setup](#environment-variables-setup):
```bash
git clone https://github.com/Merge64/x-clone.git
```

2. Start the backend by running the main file:
```bash
cd ./x-clone/server
go run main.go
```

3. Navigate into the client directory:
```bash
cd ..
cd ./x-clone/client
```

4. Install dependencies:
```bash
npm install
```

5. Run the client:
```bash
npm run dev
```

## Environment Variables Setup

1. Create a `.env` file in the root directory.
2. Copy the following content into `.env` replacing it with your credentials:

```env
# Database connection for Docker
DATABASE_URL="jdbc:postgresql://<docker-host>:<port>/<database_name>"

# Host configuration
HOST=localhost

# PostgreSQL credentials
POSTGRES_DB=<database_name>
POSTGRES_USER=<username>
POSTGRES_PASSWORD=<password>

# Ports configuration
DATABASE_PORT=5432
SERVER_PORT=8080

# Local database URL
DATABASE_URL_LOCAL="postgresql://<username>:<password>@localhost:<port>/<database_name>"

# Secret key
SECRET=<your_secret_key>

```

3. Add the `.env` file to `.gitignore` to prevent committing sensitive information:
```bash
# Ignore .env files
.env
```

## Usage
This project is intended for educational purposes only. It is designed to help improve understanding of web development and the integration of frontend and backend technologies. **Please note** that this is not a commercial project and is not meant for production use. If you wish to contribute, improve, or extend the project, feel free to create pull requests or open issues to discuss potential changes.

### How to Use the App:
1. After following the installation steps above, open the app in your browser.
2. Create an account by registering on the platform.
3. Explore features such as posting, commenting, liking, and sending direct messages (DMs).
4. Manage your account settings and try out various features available to registered users.

### Legal Disclaimer:
This project is a clone created **for educational purposes only** and **not for commercial use**. All features and designs have been developed to simulate the basic functionality of **X.com** without violating intellectual property rights. The project is not affiliated with or endorsed by X.com or any related companies.
