Simple gRPC Exploration
This project is to explore how to implement simple gRPC API with golang

📌 Overview
This project want to explore how gRPC is being implemented for API to enrich my knowledge with gRPC that can be more performant than other form of API
This project focus on simple CRUD (Create, Read, Update, Delete) for user data and provide cursor-based pagination for better performance on query

🚀 Key Features
Feature A: CRUD individual user data

Feature B: Get all user data with cursor-based pagination

Feature C: Provide seeder for initial testing

🛠 Tech Stack
Language: Golang

Database: In memory

Tools: gRPC

📥 Installation
Clone the repository:

run this code:
git clone https://github.com/Farabi99/Simple-gRPC-Exploration.git
cd Simple-gRPC-Exploration

Install dependencies:

run this code:
go mod tidy

💻 Usage
For running the server:
go run main.go

After you're add / edit the protobuf you have to run this:
protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative     proto/user.proto

the server running on port 50051

🧪 Testing
run this code:
go run Test-client/main.go
📜 License
This project is licensed under the MIT License.
