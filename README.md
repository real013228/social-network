# Social Network GraphQL API

This project implements a simple social network system where users can add posts and comments, similar to platforms like Habr or Reddit. The system provides functionality for creating posts, adding and viewing comments in a hierarchical structure, and allows users to interact with the content in real time via GraphQL Subscriptions.

## Features

### Post System
- **View Posts**: Users can fetch a list of available posts.
- **View a Post and Comments**: Users can view individual posts along with their associated comments.
- **Restrict Comments**: The post owner can disable commenting on their posts.

### Comment System
- **Hierarchical Comments**: Users can comment on posts and other comments, allowing for unlimited nested replies.
- **Comment Length Limit**: Each comment is restricted to a maximum of 2000 characters.
- **Comment Pagination**: Implemented pagination for efficient fetching of comments.
  
### GraphQL Subscriptions (Optional)
- **Real-time Updates**: Users subscribed to a post can receive notifications asynchronously when new comments are added without needing to refresh.

## Technologies

- **Programming Language**: Go
- **GraphQL**: API communication is handled using GraphQL for querying, mutations, and subscriptions.
- **Database**: Configurable to use either:
  - **In-memory storage**
  - **PostgreSQL** database
- **Containerization**: The service is Dockerized, allowing easy deployment as a Docker image.
- **Testing**: Unit tests are provided to cover implemented functionalities.

## Getting Started

### Prerequisites

Before running this project, ensure you have the following installed:

- [Go](https://golang.org/doc/install)
- [Docker](https://www.docker.com/products/docker-desktop)
- [PostgreSQL](https://www.postgresql.org/download/) (if using PostgreSQL for data persistence)

### Installation

1. Clone this repository:

   ```bash
   git clone https://github.com/real013228/social-network.git
   cd social-network
   ```

2. Build the Docker image:

   ```bash
   docker build -t social-network-api .
   ```

3. Run the Docker container with in-memory storage:

   ```bash
   docker run -p 8080:8080 \
   -e DB_OPTION=inmemory social-network-api
   ```

   To use PostgreSQL, pass the following environment variables:

   ```bash
   docker run -p 8080:8080 \
   -e POSTGRES_HOST=<POSTGRES_HOST> \
   -e POSTGRES_PORT=<POSTGRES_PORT> \
   -e POSTGRES_USER=<POSTGRES_USER> \
   -e POSTGRES_PASSWORD=<POSTGRES_PASSWORD> \
   -e POSTGRES_DBNAME=<POSTGRES_DBNAME> \
   social-network-api
   ```
   Or just use already set database in render.com for playground, use .env file:
   ```bash
   docker run -p 8080:8080 social-network-api
   ```
   
5. Access the GraphQL Playground at `http://localhost:8080/playground` to explore the API.

### Example GraphQL Queries

- **Fetch all posts:**
  ```graphql
  query posts($filter :PostsFilter) {
    posts(filter: $filter) {
      posts {
        id
        title
        description
        authorID
      }
    }
  }
  ```
  
- **Fetch all posts with comments within it:**
  ```graphql
  query posts($filter :PostsFilter) {
    posts(filter: $filter) {
      posts {
        id
        title
        description
        comments {
          id
          text
        }
        authorID
      }
    }
  }
  ```


- **Add a new post:**  
  ```graphql
   mutation CreatePost($input: CreatePostInput!) {
      createPost(input: $input) {
        post {
          id
          title
          description
          authorID
        }
      }
    }
  ```
for more details, explore the graphql documentation, that application provides
