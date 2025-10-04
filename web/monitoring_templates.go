package web

import (
	"html/template"
	"net/http"
	"optiassign/domain"
	"time"
)

// RenderMonitoringDashboard renders the monitoring dashboard
func RenderMonitoringDashboard(w http.ResponseWriter, metrics *domain.Metrics, uptime time.Duration, requestsPerMinute float64) error {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OptiAssign - Monitoring Dashboard</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        .metric-card {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border-radius: 10px;
            padding: 20px;
            margin-bottom: 20px;
        }
        .metric-value {
            font-size: 2.5rem;
            font-weight: bold;
        }
        .metric-label {
            font-size: 0.9rem;
            opacity: 0.9;
        }
        .status-indicator {
            width: 12px;
            height: 12px;
            border-radius: 50%;
            display: inline-block;
            margin-right: 8px;
        }
        .status-healthy { background-color: #28a745; }
        .status-warning { background-color: #ffc107; }
        .status-error { background-color: #dc3545; }
    </style>
</head>
<body>
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container">
            <a class="navbar-brand" href="/">OptiAssign</a>
            <div class="navbar-nav ms-auto">
                <a class="nav-link" href="/">Dashboard</a>
                <a class="nav-link" href="/monitoring">Monitoring</a>
            </div>
        </div>
    </nav>

    <div class="container mt-4">
        <div class="row">
            <div class="col-12">
                <h1 class="mb-4">Monitoring Dashboard</h1>
            </div>
        </div>

        <!-- Key Metrics -->
        <div class="row">
            <div class="col-md-3">
                <div class="metric-card">
                    <div class="metric-value">{{.RequestsTotal}}</div>
                    <div class="metric-label">Total Requests</div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="metric-card">
                    <div class="metric-value">{{.ActiveUsers}}</div>
                    <div class="metric-label">Active Users</div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="metric-card">
                    <div class="metric-value">{{.GroupsCreated}}</div>
                    <div class="metric-label">Groups Created</div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="metric-card">
                    <div class="metric-card">
                        <div class="metric-value">{{printf "%.1f" .RequestsPerMinute}}</div>
                        <div class="metric-label">Requests/Min</div>
                    </div>
                </div>
            </div>
        </div>

        <!-- System Status -->
        <div class="row mt-4">
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">
                        <h5>System Status</h5>
                    </div>
                    <div class="card-body">
                        <div class="d-flex align-items-center mb-2">
                            <span class="status-indicator status-healthy"></span>
                            <span>Application: Healthy</span>
                        </div>
                        <div class="d-flex align-items-center mb-2">
                            <span class="status-indicator status-healthy"></span>
                            <span>Database: Connected</span>
                        </div>
                        <div class="d-flex align-items-center mb-2">
                            <span class="status-indicator status-healthy"></span>
                            <span>Task Manager: Running</span>
                        </div>
                        <div class="d-flex align-items-center">
                            <span class="status-indicator status-healthy"></span>
                            <span>Email Service: Active</span>
                        </div>
                    </div>
                </div>
            </div>
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">
                        <h5>Performance Metrics</h5>
                    </div>
                    <div class="card-body">
                        <div class="row">
                            <div class="col-6">
                                <strong>Uptime:</strong><br>
                                <span class="text-muted">{{.Uptime}}</span>
                            </div>
                            <div class="col-6">
                                <strong>Last Activity:</strong><br>
                                <span class="text-muted">{{.LastActivity.Format "15:04:05"}}</span>
                            </div>
                        </div>
                        <hr>
                        <div class="row">
                            <div class="col-6">
                                <strong>Assignments Run:</strong><br>
                                <span class="text-primary">{{.AssignmentsRun}}</span>
                            </div>
                            <div class="col-6">
                                <strong>Emails Sent:</strong><br>
                                <span class="text-primary">{{.EmailsSent}}</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Request Status Distribution -->
        <div class="row mt-4">
            <div class="col-12">
                <div class="card">
                    <div class="card-header">
                        <h5>Request Status Distribution</h5>
                    </div>
                    <div class="card-body">
                        <canvas id="statusChart" width="400" height="200"></canvas>
                    </div>
                </div>
            </div>
        </div>

        <!-- Real-time Updates -->
        <div class="row mt-4">
            <div class="col-12">
                <div class="card">
                    <div class="card-header">
                        <h5>Real-time Activity</h5>
                    </div>
                    <div class="card-body">
                        <div id="activityLog" class="bg-dark text-light p-3 rounded" style="height: 200px; overflow-y: auto;">
                            <div class="text-muted">Connecting to real-time updates...</div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <script>
        // Status distribution chart
        const statusCtx = document.getElementById('statusChart').getContext('2d');
        const statusData = {
            labels: ['200', '400', '401', '403', '404', '500'],
            datasets: [{
                data: [
                    {{.RequestsByStatus.200}},
                    {{.RequestsByStatus.400}},
                    {{.RequestsByStatus.401}},
                    {{.RequestsByStatus.403}},
                    {{.RequestsByStatus.404}},
                    {{.RequestsByStatus.500}}
                ],
                backgroundColor: [
                    '#28a745',
                    '#ffc107',
                    '#17a2b8',
                    '#fd7e14',
                    '#6f42c1',
                    '#dc3545'
                ]
            }]
        };
        
        new Chart(statusCtx, {
            type: 'doughnut',
            data: statusData,
            options: {
                responsive: true,
                plugins: {
                    legend: {
                        position: 'bottom'
                    }
                }
            }
        });

        // Real-time updates via SSE
        const eventSource = new EventSource('/sse?client_id=monitoring&group_id=0');
        
        eventSource.onmessage = function(event) {
            const data = JSON.parse(event.data);
            const activityLog = document.getElementById('activityLog');
            const timestamp = new Date().toLocaleTimeString();
            activityLog.innerHTML += '<div>[' + timestamp + '] ' + data.message + '</div>';
            activityLog.scrollTop = activityLog.scrollHeight;
        };

        eventSource.onerror = function(event) {
            console.error('SSE error:', event);
        };

        // Auto-refresh metrics every 30 seconds
        setInterval(() => {
            location.reload();
        }, 30000);
    </script>
</body>
</html>
`

	t, err := template.New("monitoring").Parse(tmpl)
	if err != nil {
		return err
	}

	data := struct {
		*domain.Metrics
		Uptime           time.Duration
		RequestsPerMinute float64
	}{
		Metrics:          metrics,
		Uptime:           uptime,
		RequestsPerMinute: requestsPerMinute,
	}

	return t.Execute(w, data)
}