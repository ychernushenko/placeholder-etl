## Overview
This project demonstrates an ETL (Extract, Transform, Load) process using Go, Docker, and PostgreSQL. The ETL process involves periodically fetching data from an API, saving it to a PostgreSQL database, and then processing the data to save it as raw and transformed files on disk.  

### Process Description
1. Data Extraction:
- The service periodically fetches data from the API endpoint https://jsonplaceholder.typicode.com/posts.  
- The fetched data is saved to a PostgreSQL database.  
2. ETL Process:  
- Another goroutine runs periodically to fetch data from the PostgreSQL database.  
- The fetched data is saved on disk as raw data.  
- The raw data is then read, transformed by selecting and renaming some fields, and saved on disk as processed data.  

### Docker Compose Setup
The solution uses Docker Compose to run two containers:
1. Service Container: Runs the Go service.  
2. PostgreSQL Container: Runs the PostgreSQL database with hardcoded credentials for demonstration purposes.

### Productionization
To productionize this solution on AWS (or analogous cloud provider), consider the following:  

1. Cloud Services:
- API Gateway: Use a cloud API Gateway (e.g., AWS API Gateway) to manage and secure API requests.
- Managed Database: Use a managed PostgreSQL service (e.g., Amazon RDS) for scalability and reliability.
- Object Storage: Use cloud object storage (e.g., Amazon S3) to store raw and processed data files.
- Container Orchestration: Use a container orchestration service (e.g., Amazon ECS or Kubernetes) to manage and scale the service containers.   
  
1. Scalability and Reliability:  
- Auto-scaling: Configure auto-scaling for the service containers based on CPU and memory usage.  
- Load Balancing: Use a load balancer to distribute incoming requests across multiple instances of the service.  
Monitoring and Logging: Implement monitoring and logging using cloud services (e.g., Amazon CloudWatch) to track performance and detect issues.  
- Backup and Recovery: Set up automated backups for the PostgreSQL database and object storage to ensure data durability and recovery.  

By leveraging cloud services and best practices, you can ensure that the ETL process is scalable, reliable, and maintainable in a production environment.   

## Install (MacOS)
brew install docker-compose

## Start
make run

## Metrics
http://localhost:8080/metrics  

## Data
./data (cleaned up before every re-run)

## Logs
./logs (cleaned up before every re-run)

## Tests
make test