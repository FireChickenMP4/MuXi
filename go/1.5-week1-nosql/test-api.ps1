$BASE = "http://localhost:8080"

function Test-Api {
    param($Name, $Method, $Url, $Body)

    Write-Host "`n=== $Name ===" -ForegroundColor Cyan
    Write-Host "> $Method $Url"
    if ($Body) { Write-Host "> Body: $Body" }

    try {
        $params = @{
            Method = $Method
            Uri = "$BASE$Url"
            ContentType = "application/json"
        }
        if ($Body) { $params.Body = $Body }
        $response = Invoke-RestMethod @params -ErrorAction Stop
        Write-Host ($response | ConvertTo-Json -Depth 5) -ForegroundColor Green
        return $response
    } catch {
        Write-Host "Error: $_" -ForegroundColor Red
        return $null
    }
}

function Test-Db {
    param($Db)

    Write-Host "`n`n=========================================" -ForegroundColor Yellow
    Write-Host "  Testing $Db API" -ForegroundColor Yellow
    Write-Host "=========================================" -ForegroundColor Yellow

    # 1. Create a post with custom extensions
    $post = Test-Api -Name "Create Post" -Method Post -Url "/api/$Db/posts" -Body '{
        "title": "Welcome to NoSQL Week",
        "content": "This is a post about MongoDB and PostgreSQL comparison.",
        "author": "Alice",
        "extensions": {
            "tags": ["nosql", "database", "go"],
            "location": "Beijing",
            "images": ["https://example.com/img1.png"]
        }
    }'
    if (-not $post) { return }
    $postId = $post.id
    Write-Host "  [Created post: $postId]" -ForegroundColor Magenta

    # 2. List all posts
    Test-Api -Name "List Posts" -Method Get -Url "/api/$Db/posts"

    # 3. Get single post
    Test-Api -Name "Get Single Post" -Method Get -Url "/api/$Db/posts/$postId"

    # 4. Update post title
    Test-Api -Name "Update Post" -Method Put -Url "/api/$Db/posts/$postId" -Body '{
        "title": "Updated: NoSQL Week Recap"
    }'

    # 5. Add root comment
    $c1 = Test-Api -Name "Add Root Comment" -Method Post -Url "/api/$Db/posts/$postId/comments" -Body '{
        "content": "Great post! I learned a lot.",
        "author": "Bob"
    }'
    if (-not $c1) { return }
    $c1Id = $c1.id
    Write-Host "  [Created comment: $c1Id]" -ForegroundColor Magenta

    # 6. Reply to comment (nested reply with @mention)
    $c2 = Test-Api -Name "Reply (nested @reply)" -Method Post -Url "/api/$Db/posts/$postId/comments" -Body "{
        `"content`": `"Thanks @Bob! Glad it helped.`",
        `"author`": `"Alice`",
        `"parent_id`": `"$c1Id`",
        `"reply_to_author`": `"Bob`",
        `"reply_to_comment_id`": `"$c1Id`"
    }"
    if (-not $c2) { return }
    $c2Id = $c2.id

    # 7. Get post with nested comments (tree structure)
    Test-Api -Name "Post with Nested Comments" -Method Get -Url "/api/$Db/posts/$postId"

    # 8. Create another post with different extensions
    Test-Api -Name "Post with varied extensions" -Method Post -Url "/api/$Db/posts" -Body '{
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
    Test-Api -Name "Delete Reply Comment" -Method Delete -Url "/api/$Db/comments/$c2Id"

    # 10. Verify: only root comment remains
    Test-Api -Name "Verify Deletion" -Method Get -Url "/api/$Db/posts/$postId"
}

Write-Host "=========================================" -ForegroundColor Yellow
Write-Host "  Make sure the server is running first!" -ForegroundColor Yellow
Write-Host "  Start with: go run main.go" -ForegroundColor Yellow
Write-Host "=========================================" -ForegroundColor Yellow

Test-Db "mongo"
Test-Db "pg"

Write-Host "`n`nAll tests completed!" -ForegroundColor Green
