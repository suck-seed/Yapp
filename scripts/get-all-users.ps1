# PowerShell script to get all users from the Yapp database using Docker
# Default output format: JSON
# Optional format: table

param(
    [string]$Format = "json",   # json (default), table
    [string]$OutputFile = "",   # Optional output file
    [switch]$IncludePasswordHash = $false,  # Whether to include password hashes (not recommended)
    [switch]$Help = $false
)

# Show help if requested
if ($Help) {
    Write-Host "Usage: .\get-all-users.ps1 [OPTIONS]" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Options:" -ForegroundColor Yellow
    Write-Host "  -Format FORMAT           Output format: json (default), table" -ForegroundColor Cyan
    Write-Host "  -OutputFile FILE         Output file (default: stdout)" -ForegroundColor Cyan
    Write-Host "  -IncludePasswordHash     Include password hash (not recommended)" -ForegroundColor Cyan
    Write-Host "  -Help                    Show this help message" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Yellow
    Write-Host "  .\get-all-users.ps1" -ForegroundColor Green
    Write-Host "  .\get-all-users.ps1 -Format table" -ForegroundColor Green
    Write-Host "  .\get-all-users.ps1 -OutputFile users.json" -ForegroundColor Green
    Write-Host "  .\get-all-users.ps1 -IncludePasswordHash" -ForegroundColor Green
    exit 0
}

# Database connection parameters
$DB_CONTAINER = "yapp-postgres-1"
$DB_USER = "yappUser"
$DB_NAME = "yappDev"

# Build the SQL query
if ($IncludePasswordHash) {
    $SELECT_FIELDS = "id, username, display_name, email, phone_number, avatar_url, friend_policy, active, last_seen, created_at, updated_at, password_hash"
} else {
    $SELECT_FIELDS = "id, username, display_name, email, phone_number, avatar_url, friend_policy, active, last_seen, created_at, updated_at"
}

$QUERY = "SELECT $SELECT_FIELDS FROM users ORDER BY created_at DESC;"

Write-Host "Connecting to database container: $DB_CONTAINER" -ForegroundColor Green
Write-Host "Database: $DB_NAME, User: $DB_USER" -ForegroundColor Green
Write-Host ""

try {
    # Check if Docker container is running
    $containerStatus = docker ps --filter "name=$DB_CONTAINER" --format "{{.Status}}"
    if (-not $containerStatus) {
        Write-Error "Database container '$DB_CONTAINER' is not running. Please start it with: docker compose up -d"
        exit 1
    }

    Write-Host "Container status: $containerStatus" -ForegroundColor Yellow

    # Execute the query based on format
    switch ($Format.ToLower()) {
        "json" {
            Write-Host "Executing query in JSON format..." -ForegroundColor Yellow

            $rawData = docker exec $DB_CONTAINER psql -U $DB_USER -d $DB_NAME -t -c $QUERY

            if ($LASTEXITCODE -ne 0) {
                throw "Failed to execute query"
            }

            $users = @()
            $lines = $rawData -split "`n" | Where-Object { $_.Trim() -ne "" }

            foreach ($line in $lines) {
                $fields = $line -split "\|" | ForEach-Object { $_.Trim() }
                if ($fields.Count -ge 11) {
                    $user = [PSCustomObject]@{
                        user_id      = $fields[0]
                        username     = $fields[1]
                        display_name = if ($fields[2] -eq "") { $null } else { $fields[2] }
                        email        = $fields[3]
                        phone_number = if ($fields[4] -eq "") { $null } else { $fields[4] }
                        avatar_url   = if ($fields[5] -eq "") { $null } else { $fields[5] }
                        friend_policy= $fields[6]
                        active       = $fields[7] -eq "t"
                        last_seen    = if ($fields[8] -eq "") { $null } else { $fields[8] }
                        created_at   = $fields[9]
                        updated_at   = $fields[10]
                    }

                    if ($IncludePasswordHash -and $fields.Count -ge 12) {
                        $user | Add-Member -NotePropertyName "password_hash" -NotePropertyValue $fields[11]
                    }

                    $users += $user
                }
            }

            $output = $users | ConvertTo-Json -Depth 3

            if ($OutputFile) {
                $output | Out-File -FilePath $OutputFile -Encoding UTF8
                Write-Host "Results saved to: $OutputFile" -ForegroundColor Yellow
            } else {
                Write-Host $output
            }
        }

        "table" {
            Write-Host "Executing query in table format..." -ForegroundColor Yellow
            $output = docker exec $DB_CONTAINER psql -U $DB_USER -d $DB_NAME -c $QUERY
            if ($LASTEXITCODE -ne 0) { throw "Failed to execute query" }
            if ($OutputFile) { $output | Out-File -FilePath $OutputFile -Encoding UTF8 } else { Write-Host $output }
        }

        default {
            Write-Error "Invalid format. Use: json or table"
            exit 1
        }
    }

    Write-Host ""
    Write-Host "Query completed successfully!" -ForegroundColor Green

} catch {
    Write-Error "Error: $($_.Exception.Message)"
    Write-Host ""
    Write-Host "Troubleshooting:" -ForegroundColor Yellow
    Write-Host "1. Make sure Docker is running" -ForegroundColor Cyan
    Write-Host "2. Make sure the database container is running: docker compose up -d" -ForegroundColor Cyan
    Write-Host "3. Check container name: docker ps" -ForegroundColor Cyan
    exit 1
}

# Usage examples:
# .\get-all-users.ps1              -> JSON (default)
# .\get-all-users.ps1 -Format table
# .\get-all-users.ps1 -OutputFile "users.json"
# .\get-all-users.ps1 -IncludePasswordHash
# .\get-all-users.ps1 -Help
