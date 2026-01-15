# Grafana Dashboard for Venio

This directory contains Grafana dashboard JSON files that are automatically provisioned when Grafana starts.

## Available Dashboards

### venio-overview.json
Main dashboard showing:
- Request rate and latency
- Error rates
- Database metrics
- Redis metrics
- Authentication metrics
- Rate limiting

## Importing Dashboards

Dashboards are automatically loaded from this directory. To add a new dashboard:

1. Create/export a dashboard in Grafana UI
2. Save the JSON file to this directory
3. Restart Grafana or wait for auto-reload

## Accessing Grafana

- URL: http://localhost:3001
- Default username: admin
- Default password: admin (change in .env with GRAFANA_PASSWORD)

## Creating Custom Dashboards

You can create custom dashboards in the Grafana UI. They will persist in the Grafana volume.

To export a dashboard:
1. Open the dashboard
2. Click the "Share" icon
3. Select "Export"
4. Save as JSON
5. Copy to this directory for version control
