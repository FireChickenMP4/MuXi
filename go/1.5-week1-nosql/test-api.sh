#!/bin/bash
set -e

BASE="http://localhost:8080"
SERVER_PID=""
COMPOSE_DIR="$(cd "$(dirname "$0")" && pwd)"

cleanup() {
    echo ""
    echo "Cleaning up..."
    if [ -n "$SERVER_PID" ] && kill -0 "$SERVER_PID" 2>/dev/null; then
        kill "$SERVER_PID" 2>/dev/null
        echo "  [Go server stopped]"
    fi
    docker compose -f "$COMPOSE_DIR/docker-compose.yml" down 2>/dev/null
    echo "  [Docker containers stopped]"
}

trap cleanup EXIT INT TERM

start_server() {
    echo "Starting Docker containers..." 
    docker compose -f "$COMPOSE_DIR/docker-compose.yml" up -d 2>&1

    echo -n "Waiting for databases..."
    local mongo_ok=false pg_ok=false
    for i in $(seq 1 30); do
        $mongo_ok || docker exec week1-mongo mongosh --quiet --eval "db.runCommand({ping:1})" 2>/dev/null && mongo_ok=true
        $pg_ok || (docker exec week1-pg pg_isready -U postgres 2>/dev/null | grep -q accepting) && pg_ok=true
        $mongo_ok && $pg_ok && break
        echo -n "."
        sleep 1
    done
    echo " ready!"

    echo "Starting Go server..."
    go run "$COMPOSE_DIR" &
    SERVER_PID=$!
    sleep 4

    curl -sf "$BASE/api/mongo/posts" > /dev/null 2>&1 || {
        echo "Server failed to start"
        exit 1
    }
    echo "Server is ready!"
}

test_api() {
    local name="$1" method="$2" url="$3" body="$4"

    echo ""
    echo "=== $name ==="
    echo "> $method $url"
    if [ -n "$body" ]; then
        echo "> Body: $body"
        response=$(curl -sf "$BASE$url" -X "$method" -H "Content-Type: application/json" -d "$body" 2>/dev/null)
    else
        response=$(curl -sf "$BASE$url" -X "$method" -H "Content-Type: application/json" 2>/dev/null)
    fi

    if [ -z "$response" ]; then
        echo "Error: no response" >&2
        return 1
    fi
    echo "$response" | python3 -m json.tool 2>/dev/null || echo "$response"
}

test_db() {
    local db="$1"

    echo ""
    echo "========================================="
    echo "  Testing $db API"
    echo "========================================="

    post=$(test_api "Create Post" "POST" "/api/$db/posts" '{
        "title": "Welcome to NoSQL Week",
        "content": "This is a post about MongoDB and PostgreSQL comparison.",
        "author": "Alice",
        "extensions": {
            "tags": ["nosql", "database", "go"],
            "location": "Beijing",
            "images": ["https://example.com/img1.png"]
        }
    }') || return
    post_id=$(echo "$post" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])" 2>/dev/null)
    [ -z "$post_id" ] && { echo "Error: failed to get post id" >&2; return; }
    echo "  [Created post: $post_id]"

    test_api "List Posts" "GET" "/api/$db/posts" ""
    test_api "Get Single Post" "GET" "/api/$db/posts/$post_id" ""

    test_api "Update Post" "PUT" "/api/$db/posts/$post_id" '{
        "title": "Updated: NoSQL Week Recap"
    }'

    c1=$(test_api "Add Root Comment" "POST" "/api/$db/posts/$post_id/comments" '{
        "content": "Great post! I learned a lot.",
        "author": "Bob"
    }') || return
    c1_id=$(echo "$c1" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])" 2>/dev/null)
    [ -z "$c1_id" ] && { echo "Error: failed to get comment id" >&2; return; }
    echo "  [Created comment: $c1_id]"

    c2=$(test_api "Reply (nested @reply)" "POST" "/api/$db/posts/$post_id/comments" "{
        \"content\": \"Thanks @Bob! Glad it helped.\",
        \"author\": \"Alice\",
        \"parent_id\": \"$c1_id\",
        \"reply_to_author\": \"Bob\",
        \"reply_to_comment_id\": \"$c1_id\"
    }") || return
    c2_id=$(echo "$c2" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])" 2>/dev/null)

    test_api "Post with Nested Comments" "GET" "/api/$db/posts/$post_id" ""

    test_api "Post with varied extensions" "POST" "/api/$db/posts" '{
        "title": "Event Planning",
        "content": "Looking for team building ideas.",
        "author": "Charlie",
        "extensions": {
            "event_date": "2026-04-01",
            "priority": "high",
            "participants": ["Alice", "Bob", "Charlie"]
        }
    }'

    [ -n "$c2_id" ] && test_api "Delete Reply Comment" "DELETE" "/api/$db/comments/$c2_id" ""
    test_api "Verify Deletion" "GET" "/api/$db/posts/$post_id" ""
}

start_server
test_db "mongo"
test_db "pg"
echo ""
echo ""
echo "All tests completed!"
