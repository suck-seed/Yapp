# PowerShell script to get all halls from the Yapp database using Docker
# Default output: JSON
# Optional output: table

param(
    [string]$Format = "json",  # json (default), table
    [string]$OutputFile = "",  # Optional output file
    [switch]$Help = $false
)

# Show help if requested
if ($Help) {
    Write-Host "Usage: .\get-halls.ps1 [OPTIONS]" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Options:" -ForegroundColor Yellow
    Write-Host "  -Format FORMAT           Output format: json (default) or table" -ForegroundColor Cyan
    Write-Host "  -OutputFile FILE         Output file (default: stdout)" -ForegroundColor Cyan
    Write-Host "  -Help                    Show this help message" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Yellow
    Write-Host "  .\get-halls.ps1" -ForegroundColor Green
    Write-Host "  .\get-halls.ps1 -Format table" -ForegroundColor Green
    Write-Host "  .\get-halls.ps1 -OutputFile halls.json" -ForegroundColor Green
    exit 0
}

# Database connection parameters
$DB_CONTAINER = "yapp-postgres-1"
$DB_USER = "yappUser"
$DB_NAME = "yappDev"

# Build the SQL query for halls
$SELECT_FIELDS = "id, name, icon_url, banner_color, description, created_at, updated_at, created_by_id"
$QUERY = "SELECT $SELECT_FIELDS FROM halls ORDER BY created_at DESC;"

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

    switch ($Format.ToLower()) {
        "json" {
            Write-Host "Executing query in JSON format..." -ForegroundColor Yellow

            $rawData = docker exec $DB_CONTAINER psql -U $DB_USER -d $DB_NAME -t -c $QUERY

            if ($LASTEXITCODE -ne 0) {
                throw "Failed to execute query"
            }

            $halls = @()
            $lines = $rawData -split "`n" | Where-Object { $_.Trim() -ne "" }

            foreach ($line in $lines) {
                $fields = $line -split "\|" | ForEach-Object { $_.Trim() }
                if ($fields.Count -ge 8) {
                    $hall = [PSCustomObject]@{
                        id            = $fields[0]
                        name          = $fields[1]
                        icon_url      = if ($fields[2] -eq "") { $null } else { $fields[2] }
                        banner_color  = if ($fields[3] -eq "") { $null } else { $fields[3] }
                        description   = if ($fields[4] -eq "") { $null } else { $fields[4] }
                        created_at    = $fields[5]
                        updated_at    = $fields[6]
                        created_by_id = $fields[7]
                    }
                    $halls += $hall
                }
            }

            $output = $halls | ConvertTo-Json -Depth 3

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

            if ($LASTEXITCODE -ne 0) {
                throw "Failed to execute query"
            }

            if ($OutputFile) {
                $output | Out-File -FilePath $OutputFile -Encoding UTF8
                Write-Host "Results saved to: $OutputFile" -ForegroundColor Yellow
            } else {
                Write-Host $output
            }
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
