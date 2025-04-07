# Brevido

Brevido is a web API designed to handle a personal easily searchable list of acronyms that can be shared with others.
This API provides endpoints to create, read, update, and delete acronyms, as well as a admin panel to manage user and organization data.
It is built exclusively as an API service to support seamless integration with front-end applications or other services.

## Philosophy

Brevido is designed to be a simple, easy-to-use, and secure acronym management system.

## Overview

The Brevido API serves as a centralized service to:
- Store and manage acronyms.
- Associate acronyms with users and organizations.
- Provide secure, authenticated access to data.
- Support scalable integration with external systems.

## Features

- **Acronym Management**
  - Create, read, update, and delete acronym records.
  - Search and filter acronyms.
  
- **User Management**
  - User registration and profile management.
  - Authentication and authorization.
  
- **Organization Management**
  - CRUD operations for organization data.
  - Linking users to organizations.
  
- **Security & Compliance**
  - Secure API endpoints with token-based authentication.
  - Input validation and error handling.

- **Documentation & Testing**
  - Auto-generated API documentation.
  - Comprehensive unit and integration tests.

## Roadmap

- [x] **Project Initialization**
  - [x] Set up project structure.
  - [x] Define API endpoints.
  
- [x] **Acronym Module**
  - [x] Implement CRUD endpoints for acronyms.
  - [ ] Add search and filter functionality.
  
- [ ] **User Module**
  - [ ] Create user registration endpoint.
  - [ ] Implement authentication & authorization.
  - [ ] Develop user profile management.
  
- [ ] **Organization Module**
  - [ ] Build CRUD endpoints for organizations.
  - [ ] Enable user-organization linking.
  
- [ ] **Security Enhancements**
  - [ ] Token-based authentication.
  - [ ] Input validation and error handling.
  
- [ ] **Documentation & Testing**
  - [ ] Integrate API documentation (e.g., Swagger/OpenAPI).
  - [ ] Write unit tests for all endpoints.
  - [ ] Perform integration testing.
  
- [ ] **Deployment**
  - [ ] Containerize the application.
  - [ ] Setup CI/CD pipeline.
  - [ ] Deploy to cloud environment.

## Getting Started

### Prerequisites

- Go 1.24+
- Docker (optional, for containerization)
- Git (optional, for cloning the repository)
- Air (optional, for development)

### Installation

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/Burning-Panda/acronyms-system-api.git
   cd acronyms-system-api
   ```

2. **Build the Project:**
   ```bash
   go build -o acronyms-system-api
   ```

3. **Run the Project:**
   ```bash
   ./acronyms-system-api
   ```

4. **Access the Admin Panel:**
   ```bash
   http://localhost:8080/admin
   ```

### Development

With or without Go "Air"

1. **Run the Project without Air:**
   ```bash
   go run main.go
   ```

2. **Run the Project with Air:**
   ```bash
   air
   ```



