$BASE = "http://localhost:8080"
$SERVER_PID = $null

function Cleanup {
    Write-Host "`nCleaning up..." -ForegroundColor Yellow
    if ($SERVER_PID) {
        Stop-Process -Id $SERVER_PID -Force -ErrorAction SilentlyContinue
        Write-Host "  [Go server stopped]" -ForegroundColor Gray
    }
    $composeDir = if ($PSScriptRoot) { $PSScriptRoot } else { Get-Location }
    docker compose -f "$composeDir\docker-compose.yml" down 2>$null
    Write-Host "  [Docker containers stopped]" -ForegroundColor Gray
}

function Start-Server {
    $composeDir = if ($PSScriptRoot) { $PSScriptRoot } else { Get-Location }

    Write-Host "Starting Docker containers..." -ForegroundColor Cyan
    docker compose -f "$composeDir\docker-compose.yml" up -d 2>&1 | Out-Null

    Write-Host "Waiting for databases..." -NoNewline -ForegroundColor Cyan
    $mongoOk = $false
    $pgOk = $false
    for ($i = 0; $i -lt 30; $i++) {
        if (-not $mongoOk) {
            try { $null = docker exec week1-mongo mongosh --quiet --eval "db.runCommand({ping:1})" 2>$null; $mongoOk = $true } catch {}
        }
        if (-not $pgOk) {
            try { $r = docker exec week1-pg pg_isready -U postgres 2>$null; if ($LASTEXITCODE -eq 0) { $pgOk = $true } } catch {}
        }
        if ($mongoOk -and $pgOk) { break }
        Write-Host "." -NoNewline
        Start-Sleep -Seconds 1
    }
    Write-Host " ready!" -ForegroundColor Green

    Write-Host "Starting Go server..." -ForegroundColor Cyan
    $proc = Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run","." -WorkingDirectory $composeDir -PassThru -RedirectStandardError "$env:TEMP\go-server-err.log" -RedirectStandardOutput "$env:TEMP\go-server.log"
    $SERVER_PID = $proc.Id
    Start-Sleep -Seconds 3

    try {
        $null = Invoke-RestMethod -Uri "$BASE/api/mongo/posts" -TimeoutSec 3 -ErrorAction Stop
    } catch {
        Write-Host "Server failed to start, check logs:" -ForegroundColor Red
        Get-Content "$env:TEMP\go-server-err.log" -Tail 5
        Cleanup
        exit 1
    }
    Write-Host "Server is ready!" -ForegroundColor Green
}

function Test-Api {
    param($Name, $Method, $Url, $Body)

    Write-Host "`n=== $Name ===" -ForegroundColor Cyan
    Write-Host "> $Method $Url"
    if ($Body) { Write-Host "> Body: $Body" }

    try {
        $params = @{ Method = $Method; Uri = "$BASE$Url"; ContentType = "application/json" }
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

    Write-Host "`n=========================================" -ForegroundColor Yellow
    Write-Host "  Testing $Db API" -ForegroundColor Yellow
    Write-Host "=========================================" -ForegroundColor Yellow

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

    Test-Api -Name "List Posts" -Method Get -Url "/api/$Db/posts"
    Test-Api -Name "Get Single Post" -Method Get -Url "/api/$Db/posts/$postId"

    Test-Api -Name "Update Post" -Method Put -Url "/api/$Db/posts/$postId" -Body '{
        "title": "Updated: NoSQL Week Recap"
    }'

    $c1 = Test-Api -Name "Add Root Comment" -Method Post -Url "/api/$Db/posts/$postId/comments" -Body '{
        "content": "Great post! I learned a lot.",
        "author": "Bob"
    }'
    if (-not $c1) { return }
    $c1Id = $c1.id
    Write-Host "  [Created comment: $c1Id]" -ForegroundColor Magenta

    $c2 = Test-Api -Name "Reply (nested @reply)" -Method Post -Url "/api/$Db/posts/$postId/comments" -Body "{
        `"content`": `"Thanks @Bob! Glad it helped.`",
        `"author`": `"Alice`",
        `"parent_id`": `"$c1Id`",
        `"reply_to_author`": `"Bob`",
        `"reply_to_comment_id`": `"$c1Id`"
    }"
    if (-not $c2) { return }
    $c2Id = $c2.id

    Test-Api -Name "Post with Nested Comments" -Method Get -Url "/api/$Db/posts/$postId"

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

    Test-Api -Name "Delete Reply Comment" -Method Delete -Url "/api/$Db/comments/$c2Id"
    Test-Api -Name "Verify Deletion" -Method Get -Url "/api/$Db/posts/$postId"
}

try {
    Start-Server
    Test-Db "mongo"
    Test-Db "pg"
    Write-Host "`n`nAll tests completed!" -ForegroundColor Green
} finally {
    Cleanup
}
