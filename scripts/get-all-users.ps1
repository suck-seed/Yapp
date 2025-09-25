# PowerShell script to print all users from database
# Usage: .\get-users.ps1 [json|table]
# Default: json

param(
    [string]$Format = "json"
)

# Database connection parameters
$DB_CONTAINER = "yapp-postgres-1"
$DB_USER = "yappUser"
$DB_NAME = "yappDev"
$QUERY = "SELECT * FROM users ORDER BY created_at DESC;"

# Check if Docker container is running
try {
    $containerStatus = docker ps --filter "name=$DB_CONTAINER" --format "{{.Status}}" 2>$null
    if (-not $containerStatus) {
        Write-Error "Container '$DB_CONTAINER' not running. Start with: docker compose up -d"
        exit 1
    }
} catch {
    Write-Error "Docker not available or container not found"
    exit 1
}

switch ($Format.ToLower()) {
    "json" {
        # Get raw data from PostgreSQL
        $rawData = docker exec $DB_CONTAINER psql -U $DB_USER -d $DB_NAME -t -c $QUERY 2>$null

        if ($LASTEXITCODE -ne 0) {
            Write-Error "Failed to execute query"
            exit 1
        }

        # Parse and format as JSON
        $users = @()
        $lines = $rawData -split "`n" | Where-Object { $_.Trim() -ne "" }

        foreach ($line in $lines) {
            $fields = $line -split "\|" | ForEach-Object { $_.Trim() }
            if ($fields.Count -ge 12) {
                $user = [PSCustomObject]@{
                    id                   = $fields[0]
                    username             = $fields[1]
                    display_name         = $fields[2]
                    email                = $fields[3]
                    password_hash        = $fields[4]
                    phone_number         = if ($fields[5] -eq "") { $null } else { $fields[5] }
                    avatar_url           = if ($fields[6] -eq "") { $null } else { $fields[6] }
                    avatar_thumbnail_url = if ($fields[7] -eq "") { $null } else { $fields[7] }
                    friend_policy        = $fields[8]
                    active               = $fields[9] -eq "t"
                    last_seen            = if ($fields[10] -eq "") { $null } else { $fields[10] }
                    created_at           = $fields[11]
                    updated_at           = if ($fields.Count -gt 12) { $fields[12] } else { $null }
                }
                $users += $user
            }
        }

        # Output JSON
        $users | ConvertTo-Json -Depth 2
    }

    "table" {
        # Output as formatted table
        docker exec $DB_CONTAINER psql -U $DB_USER -d $DB_NAME -c $QUERY 2>$null

        if ($LASTEXITCODE -ne 0) {
            Write-Error "Failed to execute query"
            exit 1
        }
    }

    default {
        Write-Error "Invalid format. Use: json (default) or table"
        Write-Host "Usage: .\get-users.ps1 [json|table]"
        exit 1
    }
}
