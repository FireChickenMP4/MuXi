#!/bin/bash
BASE="http://localhost:8080"

test_api() {
    local name="$1" method="$2" url="$3" body="$4"

    echo ""
    echo "=== $name ==="
    echo "> $method $url"
    if [ -n "$body" ]; then
        echo "> Body: $body"
        response=$(curl -s -X "$method" "$BASE$url" -H "Content-Type: application/json" -d "$body")
    else
        response=$(curl -s -X "$method" "$BASE$url" -H "Content-Type: application/json")
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

    # 1. Create a post with custom extensions
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
    if [ -z "$post_id" ]; then
        echo "Error: failed to get post id" >&2
        return
    fi
    echo "  [Created post: $post_id]"

    # 2. List all posts
    test_api "List Posts" "GET" "/api/$db/posts" ""

    # 3. Get single post
    test_api "Get Single Post" "GET" "/api/$db/posts/$post_id" ""

    # 4. Update post title
    test_api "Update Post" "PUT" "/api/$db/posts/$post_id" '{
        "title": "Updated: NoSQL Week Recap"
    }'

    # 5. Add root comment
    c1=$(test_api "Add Root Comment" "POST" "/api/$db/posts/$post_id/comments" '{
        "content": "Great post! I learned a lot.",
        "author": "Bob"
    }') || return
    c1_id=$(echo "$c1" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])" 2>/dev/null)
    if [ -z "$c1_id" ]; then
        echo "Error: failed to get comment id" >&2
        return
    fi
    echo "  [Created comment: $c1_id]"

    # 6. Reply to comment (nested reply with @mention)
    c2=$(test_api "Reply (nested @reply)" "POST" "/api/$db/posts/$post_id/comments" "{
        \"content\": \"Thanks @Bob! Glad it helped.\",
        \"author\": \"Alice\",
        \"parent_id\": \"$c1_id\",
        \"reply_to_author\": \"Bob\",
        \"reply_to_comment_id\": \"$c1_id\"
    }") || return

    # 7. Get post with nested comments (tree structure)
    test_api "Post with Nested Comments" "GET" "/api/$db/posts/$post_id" ""

    # 8. Create another post with different extensions
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

    # 9. Delete the reply comment
    c2_id=$(echo "$c2" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])" 2>/dev/null)
    if [ -n "$c2_id" ]; then
        test_api "Delete Reply Comment" "DELETE" "/api/$db/comments/$c2_id" ""
    fi

    # 10. Verify: only root comment remains
    test_api "Verify Deletion" "GET" "/api/$db/posts/$post_id" ""
}

echo "========================================="
echo "  Make sure the server is running first!"
echo "  Start with: go run main.go"
echo "========================================="

test_db "mongo"
test_db "pg"

echo ""
echo ""
echo "All tests completed!"
