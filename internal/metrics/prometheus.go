package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	SyncFilesTotal       = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "sync_files_total", Help: "Synced files"}, []string{"direction", "status"})
	SyncDurationSeconds  = prometheus.NewHistogram(prometheus.HistogramOpts{Name: "sync_duration_seconds", Help: "Sync duration", Buckets: prometheus.DefBuckets})
	SFTPConnectionsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "sftp_connections_total", Help: "SFTP connections"}, []string{"server", "status"})
	SFTPBytesTransferred = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "sftp_bytes_transferred_total", Help: "SFTP bytes"}, []string{"direction"})
	ActiveSessions       = prometheus.NewGauge(prometheus.GaugeOpts{Name: "active_sessions", Help: "Active sessions"})
	DBQueryDurationSecs  = prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "db_query_duration_seconds", Help: "DB query duration", Buckets: prometheus.DefBuckets}, []string{"query"})
)

func Register() {
	prometheus.MustRegister(SyncFilesTotal, SyncDurationSeconds, SFTPConnectionsTotal, SFTPBytesTransferred, ActiveSessions, DBQueryDurationSecs)
}
