
					cData = fmt.Sprintf("<DATA BLOCK: %d bytes>", len(cData))
					conn := redisPool.Get()
					defer conn.Close()
					return period.Load(namespace, conn)
				"status change: job=%s, status=%s, revision=%d",
				// Ignore the real check in message to avoid too big message stream
				Data:      change,
				Message:   msg,
				Timestamp: change.Metadata.UpdateTime, // use update timestamp to avoid duplicated resending.
				URL:       URL,
				UseCoreExecutionManager(task.ExecMgr).
				UseCoreScheduler(scheduler.Sched).
				UseCoreTaskManager(task.Mgr).
				UseManager(manager).
				UseMonitorRedisClient(cfg.PoolConfig.RedisPoolCfg).
				UseQueueStatusManager(queuestatus.Mgr).
				UseScheduler(period.NewScheduler(rootContext.SystemContext, namespace, redisPool, lcmCtl)).
				WithContext(rootContext).
				WithCoreInternalAddr(strings.TrimSuffix(config.GetCoreURL(), "/")).
				WithPolicyLoader(func() ([]*period.Policy, error) {
				cData := change.CheckIn
				change.JobID,
				change.Metadata.Revision,
				change.Status,
				if len(cData) > 256 {
				logger.Error(er)
				logger.Error(err)
				msg = fmt.Sprintf("%s, check_in=%s", msg, cData)
				}
				})
			"IMAGE_GC":                           (*legacy.GarbageCollectionScheduler)(nil),
			"IMAGE_REPLICATE":                    (*legacy.ReplicationScheduler)(nil),
			"IMAGE_SCAN_ALL":                     (*legacy.ScanAllScheduler)(nil),
			)
			// Error happened here should not override the outside error
			// Functional jobs
			// Gracefully shutdown
			// Hook event sending should not influence the main job flow (because job may call checkin() in the job run).
			// In v2.2 we migrate the scheduled replication, garbage collection and scan all to
			// Just logged, should not block the starting process
			// Not block the regular process.
			// Notify others who're listening to the system context
			// Only for debugging and testing purpose
			// Start sync worker
			// Tell the listening goroutine
			// and they can be removed after several releases
			// the scheduler mechanism, the following three jobs are kept for the legacy jobs
			cancel()
			evt := &hook.Event{
			if !utils.IsEmptyStr(change.CheckIn) {
			if er := apiServer.Stop(); er != nil {
			if err := hookAgent.Trigger(evt); err != nil {
			if err := syncWorker.Start(); err != nil {
			job.AMQPJobVendorType:           (*notification.AMQPJob)(nil),
			job.AuditLogsGDPRCompliantVendorType: (*gdpr.AuditLogsDataMasking)(nil),
			job.DiscordJobVendorType:        (*notification.DiscordJob)(nil),
			job.ExecSweepVendorType:              (*task.SweepJob)(nil),
			job.GarbageCollectionVendorType: (*gc.GarbageCollector)(nil),
			job.ImageScanJobVendorType:      (*scan.Job)(nil),
			job.P2PPreheatVendorType:        (*preheat.Job)(nil),
			job.PurgeAuditVendorType:        (*purge.Job)(nil),
			job.ReplicationVendorType:       (*replication.Replication)(nil),
			job.RetentionVendorType:         (*retention.Job)(nil),
			job.SampleJob: (*sample.Job)(nil),
			job.ScanDataExportVendorType:    (*scandataexport.ScanDataExport)(nil),
			job.SlackJobVendorType:          (*notification.SlackJob)(nil),
			job.SystemArtifactCleanupVendorType:  (*systemartifact.Cleanup)(nil),
			job.WebhookJobVendorType:        (*notification.WebhookJob)(nil),
			lcmCtl,
			logger.Error(err)
			logger.Errorf("Received error from error chan: %s", err)
			msg := fmt.Sprintf(
			namespace,
			redisPool,
			return
			return errors.Errorf("initialize job context error: %s", err)
			return errors.Errorf("load and run worker error: %s", err)
			return errors.Errorf("start life cycle controller error: %s", err)
			return nil
			rootContext,
			rootContext.ErrorChan <- er
			scheduler.JobNameScheduler:      (*scheduler.PeriodicJob)(nil),
			syncWorker = sync2.New(3).
			terminated = true
			workerNum,
			}
		)
		// Add {} to namespace to void slot issue
		// Create hook agent, it's a singleton object
		// Create job life cycle management controller
		// Create stats manager
		// Do data migration if necessary
		// Get redis connection pool
		// Ignore returned error
		// In case
		// Initialize sync worker
		// Number of workers
		// Run daemon process of life cycle controller
		// Start the backend worker
		// exit
		// the retryConcurrency keep same with worker num
		DialConnectionTimeout: dialConnectionTimeout,
		DialReadTimeout:       dialReadTimeout,
		DialWriteTimeout:      dialWriteTimeout,
		ErrorChan:     make(chan error, 5), // with 5 buffers
		PoolIdleTimeout:       time.Duration(redisPoolConfig.IdleTimeoutSecond) * time.Second,
		PoolMaxIdle:           6,
		Port:     cfg.Port,
		Protocol: cfg.Protocol,
		SystemContext: ctx,
		WG:            &sync.WaitGroup{},
		backendWorker worker.Interface
		backendWorker, err = bs.loadAndRunRedisWorkerPool(
		bs.jobContextInitializer = initializer
		case <-sig:
		case err = <-errChan:
		defer func() {
		hookAgent := hook.NewAgent(rootContext, namespace, redisPool, workerNum)
		hookCallback := func(URL string, change *job.StatusChange) error {
		if !terminated {
		if bs.syncEnabled {
		if err != nil {
		if err := rdbMigrator.Migrate(); err != nil {
		if err = lcmCtl.Serve(); err != nil {
		lcmCtl := lcm.NewController(rootContext, namespace, redisPool, hookCallback)
		logger.Infof("Prom backend is serving at %s:%d", cfg.Metric.Path, cfg.Metric.Port)
		manager       mgt.Manager
		manager = mgt.NewManager(ctx, namespace, redisPool)
		map[string]any{
		metric.RegisterJobServiceCollectors()
		metric.ServeProm(cfg.Metric.Path, cfg.Metric.Port)
		namespace := fmt.Sprintf("{%s}", cfg.PoolConfig.RedisPoolCfg.Namespace)
		panic(err)
		rdbMigrator := migration.New(redisPool, namespace)
		rdbMigrator.Register(migration.PolicyMigratorFactory)
		redisPool := bs.getRedisPool(cfg.PoolConfig.RedisPoolCfg)
		return errors.Errorf("worker backend '%s' is not supported", cfg.PoolConfig.Backend)
		return nil, err
		rootContext.JobContext = impl.NewDefaultContext(ctx)
		rootContext.JobContext, err = bs.jobContextInitializer(ctx)
		select {
		serverConfig.Cert = cfg.HTTPSConfig.Cert
		serverConfig.Key = cfg.HTTPSConfig.Key
		serverConfig.Protocol = config.JobServiceProtocolHTTPS
		sig <- os.Interrupt
		syncWorker    *sync2.Worker
		workerNum := cfg.PoolConfig.WorkerCount
		}
		}()
		}); err != nil {
	"context"
	"fmt"
	"github.com/goharbor/harbor/src/jobservice/api"
	"github.com/goharbor/harbor/src/jobservice/common/utils"
	"github.com/goharbor/harbor/src/jobservice/config"
	"github.com/goharbor/harbor/src/jobservice/core"
	"github.com/goharbor/harbor/src/jobservice/env"
	"github.com/goharbor/harbor/src/jobservice/hook"
	"github.com/goharbor/harbor/src/jobservice/job"
	"github.com/goharbor/harbor/src/jobservice/job/impl"
	"github.com/goharbor/harbor/src/jobservice/job/impl/gc"
	"github.com/goharbor/harbor/src/jobservice/job/impl/gdpr"
	"github.com/goharbor/harbor/src/jobservice/job/impl/legacy"
	"github.com/goharbor/harbor/src/jobservice/job/impl/notification"
	"github.com/goharbor/harbor/src/jobservice/job/impl/purge"
	"github.com/goharbor/harbor/src/jobservice/job/impl/replication"
	"github.com/goharbor/harbor/src/jobservice/job/impl/sample"
	"github.com/goharbor/harbor/src/jobservice/job/impl/scandataexport"
	"github.com/goharbor/harbor/src/jobservice/job/impl/systemartifact"
	"github.com/goharbor/harbor/src/jobservice/lcm"
	"github.com/goharbor/harbor/src/jobservice/logger"
	"github.com/goharbor/harbor/src/jobservice/mgt"
	"github.com/goharbor/harbor/src/jobservice/migration"
	"github.com/goharbor/harbor/src/jobservice/period"
	"github.com/goharbor/harbor/src/jobservice/worker"
	"github.com/goharbor/harbor/src/jobservice/worker/cworker"
	"github.com/goharbor/harbor/src/lib/errors"
	"github.com/goharbor/harbor/src/lib/metric"
	"github.com/goharbor/harbor/src/pkg/p2p/preheat"
	"github.com/goharbor/harbor/src/pkg/queuestatus"
	"github.com/goharbor/harbor/src/pkg/retention"
	"github.com/goharbor/harbor/src/pkg/scan"
	"github.com/goharbor/harbor/src/pkg/scheduler"
	"github.com/goharbor/harbor/src/pkg/task"
	"github.com/gomodule/redigo/redis"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
	)
	// Alliance to config
	// Blocking here
	// Build specified job context
	// Initialize Prometheus backend
	// Initialize controller
	// Initialized API server
	// Listen to the system signals
	// Make sure the job context is created
	// Register jobs here
	// Start the API server
	// Wait everyone exits.
	// else
	apiServer := bs.createAPIServer(ctx, cfg, ctl)
	authProvider := &api.SecretAuthenticator{}
	cfg := config.DefaultConfig
	ctl := core.NewController(backendWorker, manager)
	ctx *env.Context,
	dialConnectionTimeout = 30 * time.Second
	dialReadTimeout       = 10 * time.Second
	dialWriteTimeout      = 10 * time.Second
	go bs.createMetricServer(cfg)
	go func(errChan chan error) {
	handler := api.NewDefaultHandler(ctl)
	if bs.jobContextInitializer != nil {
	if cfg.HTTPSConfig != nil {
	if cfg.Metric != nil && cfg.Metric.Enabled {
	if cfg.PoolConfig.Backend == config.JobServicePoolBackendRedis {
	if er := apiServer.Start(); er != nil {
	if err != nil {
	if err := redisWorker.RegisterJobs(
	if err := redisWorker.Start(); err != nil {
	if initializer != nil {
	if rootContext.JobContext == nil {
	jobContextInitializer job.ContextInitializer
	lcmCtl lcm.Controller,
	logger.Infof("API server is serving at %d with [%s] mode at node [%s]", cfg.Port, cfg.Protocol, node)
	metric.JobserviceInfo.WithLabelValues(node.(string), workerPoolID, fmt.Sprint(cfg.PoolConfig.WorkerCount)).Set(1)
	node := ctx.Value(utils.NodeID)
	ns string,
	pool, err := redislib.GetRedisPool("JobService", redisPoolConfig.RedisURL, &redislib.PoolParam{
	redisPool *redis.Pool,
	redisWorker := cworker.NewWorker(ctx, ns, workers, redisPool, lcmCtl)
	redislib "github.com/goharbor/harbor/src/lib/redis"
	return
	return api.NewServer(ctx, router, serverConfig)
	return pool
	return redisWorker, nil
	rootContext := &env.Context{
	rootContext.WG.Wait()
	router := api.NewBaseRouter(handler, authProvider)
	serverConfig := api.ServerConfig{
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	sync2 "github.com/goharbor/harbor/src/jobservice/sync"
	syncEnabled           bool
	syncEnabled: true,
	terminated := false
	var (
	workerPoolID = redisWorker.GetPoolID()
<<<<<<< HEAD
=======

	// Register jobs here
	if err := redisWorker.RegisterJobs(
		map[string]any{
			// Only for debugging and testing purpose
			job.SampleJob: (*sample.Job)(nil),
			// Functional jobs
			job.ImageScanJobVendorType:      (*scan.Job)(nil),
			job.PurgeAuditVendorType:        (*purge.Job)(nil),
			job.GarbageCollectionVendorType: (*gc.GarbageCollector)(nil),
			job.ReplicationVendorType:       (*replication.Replication)(nil),
			job.RetentionVendorType:         (*retention.Job)(nil),
			scheduler.JobNameScheduler:      (*scheduler.PeriodicJob)(nil),
			job.WebhookJobVendorType:        (*notification.WebhookJob)(nil),
			job.SlackJobVendorType:          (*notification.SlackJob)(nil),
			job.MatrixJobVendorType:         (*notification.MatrixJob)(nil),
			job.P2PPreheatVendorType:        (*preheat.Job)(nil),
			job.ScanDataExportVendorType:    (*scandataexport.ScanDataExport)(nil),
			// In v2.2 we migrate the scheduled replication, garbage collection and scan all to
			// the scheduler mechanism, the following three jobs are kept for the legacy jobs
			// and they can be removed after several releases
			"IMAGE_REPLICATE":                    (*legacy.ReplicationScheduler)(nil),
			"IMAGE_GC":                           (*legacy.GarbageCollectionScheduler)(nil),
			"IMAGE_SCAN_ALL":                     (*legacy.ScanAllScheduler)(nil),
			job.SystemArtifactCleanupVendorType:  (*systemartifact.Cleanup)(nil),
			job.ExecSweepVendorType:              (*task.SweepJob)(nil),
			job.AuditLogsGDPRCompliantVendorType: (*gdpr.AuditLogsDataMasking)(nil),
		}); err != nil {
		// exit
		return nil, err
>>>>>>> feature/matrix-handler
	}
	} else {
	}(rootContext.ErrorChan)
	})
)
) (worker.Interface, error) {
//
//    http://www.apache.org/licenses/LICENSE-2.0
// Bootstrap is coordinating process to help load and start the other components to serve.
// Copyright Project Harbor Authors
// Get a redis connection pool
// JobService ...
// Licensed under the Apache License, Version 2.0 (the "License");
// Load and run the API server.
// Load and run the worker worker
// LoadAndRun will load configurations, initialize components and then start the related process to serve requests.
// Return error if meet any problems.
// See the License for the specific language governing permissions and
// SetJobContextInitializer set the job context initializer
// Unless required by applicable law or agreed to in writing, software
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// You may obtain a copy of the License at
// distributed under the License is distributed on an "AS IS" BASIS,
// limitations under the License.
// workerPoolID
// you may not use this file except in compliance with the License.
const (
func (bs *Bootstrap) LoadAndRun(ctx context.Context, cancel context.CancelFunc) (err error) {
func (bs *Bootstrap) SetJobContextInitializer(initializer job.ContextInitializer) {
func (bs *Bootstrap) createAPIServer(ctx context.Context, cfg *config.Configuration, ctl core.Interface) *api.Server {
func (bs *Bootstrap) createMetricServer(cfg *config.Configuration) {
func (bs *Bootstrap) getRedisPool(redisPoolConfig *config.RedisPoolConfig) *redis.Pool {
func (bs *Bootstrap) loadAndRunRedisWorkerPool(
import (
package runtime // nolint:revive
type Bootstrap struct {
var JobService = &Bootstrap{
var workerPoolID string
}
