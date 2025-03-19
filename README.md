# Financial Assistance Scheme API

## Overview
This API, built with **Golang and Gin**, manages financial assistance schemes, applicants, and applications. It allows system administrators to register applicants, view eligible schemes, and submit applications.

## Features
- **CRUD operations** for applicants, schemes, and applications.
- **Eligibility verification** based on applicant and scheme.
- **RESTful API design** using Gin framework.1eee1eee1eee
- **PostgreSQL support** for data storage.

## Prerequisites
- Install [Go](https://golang.org/doc/install) (at least v1.18)
- Install [PostgreSQL](https://www.postgresql.org/)

---

## Setup Instructions

### 1. Clone the Repository
```sh
git clone https://github.com/duweigg/FASMS-BE.git
cd FASMS-BE
```

### 2. Install Dependencies
```sh
go mod tidy
```
### 3. create database "FASMS"
in postgresql, create a database named "FASMS"

### 4. Configure Environment Variables
Create a `.env` file:
```env
PORT=8000
DB_URL="host=localhost user=postgres password=admin dbname=FASMS port=5432 sslmode=disable"
```
or copy the file .env.example. and rename it as .env

### 5. Run Database Migrations
Using GORM, initialize the database:
```sh
go run migrate/migrate.go
```

### 5. Start the Server
```sh
go run main.go
```

### 6. API Documentation
- **Postman Collection**: Import the provided JSON file. and add envvironmnet varible 
```
domain:"http://localhost:8000"
```

---

## API Endpoints

| Method | Endpoint | Description | remarks |
|--------|----------|-------------|---------|
| `GET` | `/api/applicants` | Retrieve all applicants | will need query param of page (default as 0) and page_size (default as 10). the first page is page 0 |
| `POST` | `/api/applicants` | Create a new applicant | allow batch creatation. Please refer the payload in postman file |
| `PUT` | `/api/applicants/{id}` | update existing applicant | The logic will compare the applicant's data, as well as households' data, so need to post the entire applicant data with households data including their UUIDs |
| `DELETE` | `/api/applicants/{id}` | delete existing applicant | this will soft delete the applicant as well as his households, and updated related application record to "need review" status |
| `GET` | `/api/schemes` | Retrieve all schemes | will need query param of page (default as 0) and page_size (default as 10). the first page is page 0 |
| `GET` | `/api/schemes/eligible?applicant={id}` | Retrieve eligible schemes for an applicant | In order to be eligible, applicant must satisify all the criteria groups, each criteria group is considered as satisified if any of the criteria within the criteria groupo is satisified |
| `POST` | `/api/schemes` | create new schemes | allow batch creatation. Please refer the payload in postman file |
| `PUT` | `/api/schemes/{id}` | update existing schemes | The logic will compare the scheme's data, as well as all its criteria and benefits data, so need to post the entire scheme data with  criteria and benefits data including their UUIDs |
| `DELETE` | `/api/schemes/{id}` | delete existing schemes | this will soft delete the scheme as well as its criteria and benefits, and updated related application record to "need review" status |
| `GET` | `/api/applications` | Retrieve all applications | will need query param of page (default as 0) and page_size (default as 10). the first page is page 0 |
| `POST` | `/api/applications` | Submit a new application | Please refer the payload in postman file |
| `PUT` | `/api/applications/{id}` | update existing application | Please refer the payload in postman file |
| `DELETE` | `/api/applications/{id}` | Delete existing application | this will soft delete the application |

For full API details, check the **Postman Collection**.

---

## Database Schema

### Applicants Table
![alt text](<FASMS - public.png>)

### Application constant 
| column | Value | Meaning |
|--------|----------|-------------|
| employ status | 1 | unemployed |
| employ status | 2 | employed |
| employ status | 3 | in school |
| employ status | 99 | no employment requiment |
| sex | 1 | male |
| sex | 2 | female |
| sex | 99 | no sex requiment |
| relation | 1 | children |
| relation | 2 | spouse |
| relation | 3 | parents |
| relation | 99 | no relation requiment |
| application_status | 1 | submitted |
| application_status | 2 | approved |
| application_status | 3 | rejected |
| application_status | 4 | need review |
---

## Testing

### API Testing
Use Postman to test API requests.

---