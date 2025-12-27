# MSM Cron Configuration

The cron file is now generated dynamically from your `msm.conf` settings.

## Generating and Installing

```bash
# Preview the generated cron file
msm cron generate

# Install to /etc/cron.d/msm (requires root)
sudo msm cron install
```

## Configuration Options

Add these to your `/etc/msm.conf`:

```bash
# Path to the msm binary (default: /usr/local/bin/msm)
CRON_MSM_BINARY="/usr/local/bin/msm"

# Hour for maintenance tasks: backups and archive cleanup (default: 5)
CRON_MAINTENANCE_HOUR="11"

# Days to retain archive files before deletion (default: 30)
CRON_ARCHIVE_RETENTION_DAYS="30"
```

## Scheduled Tasks

The generated cron file includes:

| Time | Task |
|------|------|
| (maintenance-1):55 | Log rolling |
| maintenance:02 | World backups |
| maintenance:10-16 | Archive cleanup |
| @hourly | Crash recovery |

Note: RAM disk syncing is handled automatically by the `msm-sync` screen session when servers are running — no cron job required.
