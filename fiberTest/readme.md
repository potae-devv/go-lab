# Fiber API Test

This is a simple REST API built with the Fiber web framework for Go.

## Project Structure

The project contains a basic book management API with the following components:
- `main.go`: Contains the API implementation with book CRUD operations
- Book struct with ID, Title, and Author fields

## Development

### Prerequisites

- Go installed on your machine
- Node.js and npm (for nodemon)

### Installing nodemon

If you don't have nodemon installed globally, you can install it with:

```bash
npm install -g nodemon
```

### Running with Hot Reload

To run the application with hot reload (automatically restart when code changes):

```bash
nodemon --watch . --ext go --exec go run . --signal SIGTERM
```