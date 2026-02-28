# Demo API Quick Reference

## Endpoints Overview

| Method | Endpoint                 | Purpose            | Auth Required     |
| ------ | ------------------------ | ------------------ | ----------------- |
| POST   | `/v1/demos/upload`       | Upload demo file   | ✅ Secret password |
| GET    | `/v1/demos`              | Get demo info      | ❌ Public          |
| GET    | `/v1/demos/download/:id` | Download file      | ❌ Public          |
| GET    | `/v1/demos/list`         | List demos         | ❌ Public          |
| HEAD   | `/v1/demos`              | Check availability | ❌ Public          |

## Key Features

- ⚡ **Automatic Compression**: Files compressed with Zstd algorithm for storage
- 📤 **Smart Downloads**: Files automatically decompressed when downloaded
- 🎮 **Tournament Support**: Organized by tournament/match/round structure  
- 📁 **Smart Storage**: MinIO object storage with PostgreSQL metadata
- 🔍 **Flexible Queries**: Search by match, round, or tournament context
- 📄 **Pagination**: Efficient listing with page/limit support

## Quick Examples

### Upload Demo
```bash
curl -X POST "https://api.udl.tf/v1/demos/upload?secret_password=secret" \
  -F "demo_file=@demo.dem" \
  -F "is_tournament=true" \
  -F "tournament_id=1" \
  -F "match_id=1" \
  -F "round_id=1" \
  -F "raw_demo_name=demo.dem" \
  -F "blue_team_name=Team A" \
  -F "red_team_name=Team B"
```

### Get Demo Info
```bash
curl "https://api.udl.tf/v1/demos?match_id=1&round_id=1"
```

### Download Demo
```bash
curl -O "https://api.udl.tf/v1/demos/download/1"
# Downloads original demo file (automatically decompressed)
```

For detailed documentation, see [demos-api.md](./demos-api.md).