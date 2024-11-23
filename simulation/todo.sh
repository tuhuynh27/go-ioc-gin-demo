#!/bin/bash

# Check if mode argument is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 <mode>"
    echo "Modes: create, clean"
    exit 1
fi

MODE=$1

# Function to generate random title and description
generate_todo() {
    local priorities=("LOW" "MEDIUM" "HIGH")
    local titles=("Meeting" "Review" "Project" "Task" "Planning" "Research" "Development" "Testing" "Documentation" "Deployment")
    local actions=("Complete" "Start" "Review" "Update" "Prepare" "Analyze" "Implement" "Test" "Document" "Deploy")
    
    local random_title="${actions[$((RANDOM % 10))]} ${titles[$((RANDOM % 10))]}"
    local random_desc="Description for: $random_title"
    local random_priority=${priorities[$((RANDOM % 3))]}
    
    echo "{\"title\":\"$random_title\",\"description\":\"$random_desc\",\"priority\":\"$random_priority\"}"
}

if [ "$MODE" = "create" ]; then
    # Create 100 todos
    for i in {1..100}; do
        todo_data=$(generate_todo)
        
        # Send POST request to create todo
        response=$(curl -s -X POST \
            -H "Content-Type: application/json" \
            -d "$todo_data" \
            http://localhost:8080/todos)
            
        # Check if request was successful
        if [[ $response == *"id"* ]]; then
            echo "Created todo #$i successfully"
        else
            echo "Failed to create todo #$i"
            # Print the error response
            echo "Error response: $response"
            # Print the request payload for debugging
            echo "Request payload: $todo_data"
        fi
        
        # Add small delay to prevent overwhelming the server
        sleep 0.1
    done

    echo "Completed generating 100 todos"

elif [ "$MODE" = "clean" ]; then
    # Get all todos
    todos=$(curl -s -X GET http://localhost:8080/todos)
    
    # Extract IDs and delete each todo
    echo "$todos" | grep -o '"id":[0-9]*' | grep -o '[0-9]*' | while read -r id; do
        response=$(curl -s -X DELETE http://localhost:8080/todos/$id)
        if [ -z "$response" ]; then
            echo "Deleted todo #$id successfully"
        else
            echo "Failed to delete todo #$id"
            echo "Error response: $response"
        fi
        sleep 0.1
    done

    echo "Completed cleaning up todos"

else
    echo "Invalid mode: $MODE"
    echo "Valid modes are: create, clean"
    exit 1
fi
