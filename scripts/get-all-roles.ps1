# PowerShell script to get all roles from the Yapp database using Docker
# Default output format: JSON
# Optional format: table

param(
    [string]$Format = "json",   # json (default), table
    [string]$OutputFile = "",   # Optional output file
    [string]$HallId = "",       # Filter by specific hall ID
    [switch]$Help = $false
)

# Show help if requested
if ($Help) {
    Write-Host "Usage: .\get-all-roles.ps1 [OPTIONS]" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Options:" -ForegroundColor Yellow
    Write-Host "  -Format FORMAT           Output format: json (default), table" -ForegroundColor Cyan
    Write-Host "  -OutputFile FILE         Output file (default: stdout)" -ForegroundColor Cyan
    Write-Host "  -HallId HALL_ID          Filter roles by specific hall ID" -ForegroundColor Cyan
    Write-Host "  -Help                    Show this help message" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Yellow
    Write-Host "  .\get-all-roles.ps1" -ForegroundColor Green
    Write-Host "  .\get-all-roles.ps1 -Format table" -ForegroundColor Green
    Write-Host "  .\get-all-roles.ps1 -OutputFile roles.json" -ForegroundColor Green
    Write-Host "  .\get-all-roles.ps1 -HallId 'a1b2c3d4-e5f6-7890-abcd-ef1234567890'" -ForegroundColor Green
    exit 0
}

# Database connection parameters
$DB_CONTAINER = "yapp-postgres-1"
$DB_USER = "yappUser"
$DB_NAME = "yappDev"

# SQL query - all role fields
$SELECT_FIELDS = "id, hall_id, name, color, icon_url, is_default, is_admin, created_at, updated_at"

# Build query with optional hall filter
if ($HallId) {
    $QUERY = "SELECT $SELECT_FIELDS FROM roles WHERE hall_id = '$HallId' ORDER BY created_at DESC;"
    Write-Host "Filtering roles for hall ID: $HallId" -ForegroundColor Cyan
} else {
    $QUERY = "SELECT $SELECT_FIELDS FROM roles ORDER BY created_at DESC;"
}

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

            $roles = @()
            $lines = $rawData -split "`n" | Where-Object { $_.Trim() -ne "" }

            foreach ($line in $lines) {
                $fields = $line -split "\|" | ForEach-Object { $_.Trim() }
                if ($fields.Count -ge 9) {
                    $role = [PSCustomObject]@{
                        role_id      = $fields[0]
                        hall_id      = $fields[1]
                        name         = $fields[2]
                        color        = if ($fields[3] -eq "") { $null } else { $fields[3] }
                        icon_url     = if ($fields[4] -eq "") { $null } else { $fields[4] }
                        is_default   = $fields[5] -eq "t"
                        is_admin     = $fields[6] -eq "t"
                        created_at   = $fields[7]
                        updated_at   = $fields[8]
                    }

                    $roles += $role
                }
            }

            if ($roles.Count -eq 0) {
                if ($HallId) {
                    Write-Host "No roles found for hall ID: $HallId" -ForegroundColor Yellow
                } else {
                    Write-Host "No roles found in the database" -ForegroundColor Yellow
                }
                return
            }

            $output = $roles | ConvertTo-Json -Depth 3

            if ($OutputFile) {
                $output | Out-File -FilePath $OutputFile -Encoding UTF8
                Write-Host "Results saved to: $OutputFile" -ForegroundColor Yellow
                Write-Host "Found $($roles.Count) role(s)" -ForegroundColor Green
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
    Write-Host "4. If using -HallId, verify the UUID format is correct" -ForegroundColor Cyan
    exit 1
}

# Usage examples:
# .\get-all-roles.ps1                           -> All roles in JSON
# .\get-all-roles.ps1 -Format table             -> All roles in table format
# .\get-all-roles.ps1 -OutputFile "roles.json"  -> Save all roles to file
# .\get-all-roles.ps1 -HallId "uuid-here"       -> Filter by specific hall
# .\get-all-roles.ps1 -Help                     -> Show help
