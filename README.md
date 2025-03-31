# Gator (Blog AggroGATOR) 

A CLI tool for managing posts from various rss feeds. 

### Prerequisites

Before you can run the program, ensure you have the following installed:

    PostgreSQL: Required to store and retrieve data.
    Go: Version 1.24 or higher is required.

## Installation

**Step 1: Install Golang** 

Please follow the instructions located here: https://golang.org/dl/. 

**Step 2: Install Gator**

To install the Gator CLI tool, use the following go install command:

`go install github.com/AleksZieba/gator@latest`

**Step 3: Set up the Configuration File** 

Create a .gatorc 

`{
  "db_url": "postgresql://USERNAME:PASSWORD@localhost:5432/DBNAME?sslmode=disable",
  "current_user_name": "YOUR_USER_NAME"
}` 

Please be sure to replace `USERNAME`, `PASSWORD`, `DBNAME` and `YOUR_USER_NAME` with the appropriate values. 

## Available Commands 
- login: Log in to the system.
- register: Register a new user.
- reset: Reset the user’s password (this will reset the database for the user).
- users: List all users.
- agg: Aggregate data.
- addfeed: Add a new feed.
- feeds: View all available feeds.
- follow: Follow a feed.
- following: View feeds you're following.
- unfollow: Unfollow a feed.
- browse: Browse posts from the feeds you follow.
