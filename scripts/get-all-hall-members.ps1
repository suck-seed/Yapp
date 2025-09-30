# PowerShell script to get all hall members from the Yapp database using Docker
# Default output format: JSON
# Optional format: table

param(
    [string]$Format = "json",   # json (default), table
    [string]$OutputFile = "",   # Optional output file
    [string]$HallId = "",       # Filter by specific hall ID
    [string]$UserId = "",       # Filter by specific user ID
    [switch]$IncludeUserInfo = $false,  # Include user details in output
    [switch]$IncludeRoleInfo = $false,  # Include role details in output
    [switch]$Help = $false
)

# Show help if requested
if ($Help) {
    Write-Host "Usage: .\get-all-hall-members.ps1 [OPTIONS]" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Options:" -ForegroundColor Yellow
    Write-Host "  -Format FORMAT           Output format: json (default), table" -ForegroundColor Cyan
    Write-Host "  -OutputFile FILE         Output file (default: stdout)" -ForegroundColor Cyan
    Write-Host "  -HallId HALL_ID          Filter members by specific hall ID" -ForegroundColor Cyan
    Write-Host "  -UserId USER_ID          Filter memberships by specific user ID" -ForegroundColor Cyan
    Write-Host "  -IncludeUserInfo         Include user details (username, display_name, email)" -ForegroundColor Cyan
    Write-Host "  -IncludeRoleInfo         Include role details (name, is_admin, is_default)" -ForegroundColor Cyan
    Write-Host "  -Help                    Show this help message" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Yellow
    Write-Host "  .\get-all-hall-members.ps1" -ForegroundColor Green
    Write-Host "  .\get-all-hall-members.ps1 -Format table" -ForegroundColor Green
    Write-Host "  .\get-all-hall-members.ps1 -HallId 'a1b2c3d4-e5f6-7890-abcd-ef1234567890'" -ForegroundColor Green
    Write-Host "  .\get-all-hall-members.ps1 -UserId 'user-uuid-here'" -ForegroundColor Green
    Write-Host "  .\get-all-hall-members.ps1 -IncludeUserInfo -IncludeRoleInfo" -ForegroundColor Green
    Write-Host "  .\get-all-hall-members.ps1 -OutputFile hall-members.json" -ForegroundColor Green
    exit 0
}

# Database connection parameters
$DB_CONTAINER = "yapp-postgres-1"
$DB_USER = "yappUser"
$DB_NAME = "yappDev"

# Build SQL query based on options
if ($IncludeUserInfo -and $IncludeRoleInfo) {
    $SELECT_FIELDS = @"
hm.id, hm.hall_id, hm.user_id, hm.joined_at, hm.role_id, hm.created_at, hm.updated_at,
u.username, u.display_name, u.email,
r.name as role_name, r.is_admin, r.is_default
"@
    $FROM_CLAUSE = @"
FROM hall_members hm
LEFT JOIN users u ON hm.user_id = u.id
LEFT JOIN roles r ON hm.role_id = r.id
"@
} elseif ($IncludeUserInfo) {
    $SELECT_FIELDS = @"
hm.id, hm.hall_id, hm.user_id, hm.joined_at, hm.role_id, hm.created_at, hm.updated_at,
u.username, u.display_name, u.email
"@
    $FROM_CLAUSE = @"
FROM hall_members hm
LEFT JOIN users u ON hm.user_id = u.id
"@
} elseif ($IncludeRoleInfo) {
    $SELECT_FIELDS = @"
hm.id, hm.hall_id, hm.user_id, hm.joined_at, hm.role_id, hm.created_at, hm.updated_at,
r.name as role_name, r.is_admin, r.is_default
"@
    $FROM_CLAUSE = @"
FROM hall_members hm
LEFT JOIN roles r ON hm.role_id = r.id
"@
} else {
    $SELECT_FIELDS = "id, hall_id, user_id, joined_at, role_id, created_at, updated_at"
    $FROM_CLAUSE = "FROM hall_members"
}

# Build WHERE clause
$WHERE_CONDITIONS = @()
if ($HallId) {
    $WHERE_CONDITIONS += "hall_id = '$HallId'"
}
if ($UserId) {
    $WHERE_CONDITIONS += "user_id = '$UserId'"
}

$WHERE_CLAUSE = if ($WHERE_CONDITIONS.Count -gt 0) {
    "WHERE " + ($WHERE_CONDITIONS -join " AND ")
} else {
    ""
}

$QUERY = "SELECT $SELECT_FIELDS $FROM_CLAUSE $WHERE_CLAUSE ORDER BY joined_at DESC;"

# Display filter info
if ($HallId) { Write-Host "Filtering members for hall ID: $HallId" -ForegroundColor Cyan }
if ($UserId) { Write-Host "Filtering memberships for user ID: $UserId" -ForegroundColor Cyan }
if ($IncludeUserInfo) { Write-Host "Including user information" -ForegroundColor Cyan }
if ($IncludeRoleInfo) { Write-Host "Including role information" -ForegroundColor Cyan }

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

            $members = @()
            $lines = $rawData -split "`n" | Where-Object { $_.Trim() -ne "" }

            foreach ($line in $lines) {
                $fields = $line -split "\|" | ForEach-Object { $_.Trim() }
                if ($fields.Count -ge 7) {
                    $member = [PSCustomObject]@{
                        member_id    = $fields[0]
                        hall_id      = $fields[1]
                        user_id      = $fields[2]
                        joined_at    = $fields[3]
                        role_id      = $fields[4]
                        created_at   = $fields[5]
                        updated_at   = $fields[6]
                    }

                    # Add user info if included
                    if ($IncludeUserInfo) {
                        $member | Add-Member -NotePropertyName "username" -NotePropertyValue ($fields[7] ?? $null)
                        $member | Add-Member -NotePropertyName "display_name" -NotePropertyValue (if ($fields[8] -eq "") { $null } else { $fields[8] })
                        $member | Add-Member -NotePropertyName "email" -NotePropertyValue ($fields[9] ?? $null)
                    }

                    # Add role info if included
                    if ($IncludeRoleInfo) {
                        $roleNameIndex = if ($IncludeUserInfo) { 10 } else { 7 }
                        $member | Add-Member -NotePropertyName "role_name" -NotePropertyValue ($fields[$roleNameIndex] ?? $null)
                        $member | Add-Member -NotePropertyName "is_admin" -NotePropertyValue ($fields[$roleNameIndex + 1] -eq "t")
                        $member | Add-Member -NotePropertyName "is_default" -NotePropertyValue ($fields[$roleNameIndex + 2] -eq "t")
                    }

                    $members += $member
                }
            }

            if ($members.Count -eq 0) {
                $message = "No hall members found"
                if ($HallId -and $UserId) {
                    $message += " for hall ID: $HallId and user ID: $UserId"
                } elseif ($HallId) {
                    $message += " for hall ID: $HallId"
                } elseif ($UserId) {
                    $message += " for user ID: $UserId"
                }
                Write-Host $message -ForegroundColor Yellow
                return
            }

            $output = $members | ConvertTo-Json -Depth 3

            if ($OutputFile) {
                $output | Out-File -FilePath $OutputFile -Encoding UTF8
                Write-Host "Results saved to: $OutputFile" -ForegroundColor Yellow
                Write-Host "Found $($members.Count) hall member(s)" -ForegroundColor Green
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
    Write-Host "4. If using filters, verify the UUID format is correct" -ForegroundColor Cyan
    exit 1
}

# Usage examples:
# .\get-all-hall-members.ps1                                    -> All hall members
# .\get-all-hall-members.ps1 -Format table                      -> Table format
# .\get-all-hall-members.ps1 -HallId "uuid-here"                -> Members of specific hall
# .\get-all-hall-members.ps1 -UserId "uuid-here"                -> All halls for specific user
# .\get-all-hall-members.ps1 -IncludeUserInfo -IncludeRoleInfo  -> Full details
# .\get-all-hall-members.ps1 -OutputFile "members.json"         -> Save to file
# .\get-all-hall-members.ps1 -Help                              -> Show help
