
<<<<<<< HEAD
		P2PPreheatVendorType:            lib.GetEnvInt64("P2P_PREHEAT_EXECUTION_RETENTION_COUNT", 50),
		PurgeAuditVendorType:            lib.GetEnvInt64("PURGE_AUDIT_EXECUTION_RETENTION_COUNT", 10),
		ReplicationVendorType:           lib.GetEnvInt64("REPLICATION_EXECUTION_RETENTION_COUNT", 50),
		RetentionVendorType:             lib.GetEnvInt64("RETENTION_EXECUTION_RETENTION_COUNT", 50),
		SBOMJobVendorType:               lib.GetEnvInt64("SBOM_EXECUTION_RETENTION_COUNT", 1),
		ScanAllVendorType:               lib.GetEnvInt64("SCAN_ALL_EXECUTION_RETENTION_COUNT", 1),
=======
package job

import "github.com/goharbor/harbor/src/lib"

// Define the register name constants of known jobs

const (
	// SampleJob is name of demo job
	SampleJob = "DEMO"

	// ImageScanJobVendorType is name of scan job it will be used as key to register to job service.
	ImageScanJobVendorType = "IMAGE_SCAN"
	// SBOMJobVendorType key to create sbom generate execution.
	SBOMJobVendorType = "SBOM"
	// GarbageCollectionVendorType job name
	GarbageCollectionVendorType = "GARBAGE_COLLECTION"
	// ReplicationVendorType : the name of the replication job in job service
	ReplicationVendorType = "REPLICATION"
	// WebhookJobVendorType : the name of the webhook job in job service
	WebhookJobVendorType = "WEBHOOK"
	// SlackJobVendorType : the name of the slack job in job service
	SlackJobVendorType = "SLACK"
	// MatrixJobVendorType : the name of the matrix job in job service
	MatrixJobVendorType = "MATRIX"
	// RetentionVendorType : the name of the retention job
	RetentionVendorType = "RETENTION"
	// P2PPreheatVendorType : the name of the P2P preheat job
	P2PPreheatVendorType = "P2P_PREHEAT"
	// PurgeAuditVendorType : the name of purge audit job
	PurgeAuditVendorType = "PURGE_AUDIT_LOG"
	// SystemArtifactCleanupVendorType : the name of the SystemArtifact cleanup job
	SystemArtifactCleanupVendorType = "SYSTEM_ARTIFACT_CLEANUP"
	// ScanDataExportVendorType : the name of the scan data export job
	ScanDataExportVendorType = "SCAN_DATA_EXPORT"
	// ExecSweepVendorType: the name of the execution sweep job
	ExecSweepVendorType = "EXECUTION_SWEEP"
	// ScanAllVendorType: the name of the scan all job
	ScanAllVendorType = "SCAN_ALL"
	// AuditLogsGDPRCompliantVendorType : the name of the job which makes audit logs table GDPR-compliant
	AuditLogsGDPRCompliantVendorType = "AUDIT_LOGS_GDPR_COMPLIANT"
)

var (
	// executionSweeperCount stores the count for execution retained
	executionSweeperCount = map[string]int64{
		ImageScanJobVendorType:          lib.GetEnvInt64("IMAGE_SCAN_EXECUTION_RETENTION_COUNT", 1),
		SBOMJobVendorType:               lib.GetEnvInt64("SBOM_EXECUTION_RETENTION_COUNT", 1),
		ScanAllVendorType:               lib.GetEnvInt64("SCAN_ALL_EXECUTION_RETENTION_COUNT", 1),
		PurgeAuditVendorType:            lib.GetEnvInt64("PURGE_AUDIT_EXECUTION_RETENTION_COUNT", 10),
		ExecSweepVendorType:             lib.GetEnvInt64("EXECUTION_SWEEP_EXECUTION_RETENTION_COUNT", 10),
		GarbageCollectionVendorType:     lib.GetEnvInt64("GARBAGE_COLLECTION_EXECUTION_RETENTION_COUNT", 50),
		SlackJobVendorType:              lib.GetEnvInt64("SLACK_EXECUTION_RETENTION_COUNT", 50),
		MatrixJobVendorType:             lib.GetEnvInt64("MATRIX_EXECUTION_RETENTION_COUNT", 50),
		WebhookJobVendorType:            lib.GetEnvInt64("WEBHOOK_EXECUTION_RETENTION_COUNT", 50),
		ReplicationVendorType:           lib.GetEnvInt64("REPLICATION_EXECUTION_RETENTION_COUNT", 50),
>>>>>>> feature/matrix-handler
		ScanDataExportVendorType:        lib.GetEnvInt64("SCAN_DATA_EXPORT_EXECUTION_RETENTION_COUNT", 50),
		SlackJobVendorType:              lib.GetEnvInt64("SLACK_EXECUTION_RETENTION_COUNT", 50),
		SystemArtifactCleanupVendorType: lib.GetEnvInt64("SYSTEM_ARTIFACT_CLEANUP_EXECUTION_RETENTION_COUNT", 50),
		WebhookJobVendorType:            lib.GetEnvInt64("WEBHOOK_EXECUTION_RETENTION_COUNT", 50),
	// AMQPJobVendorType : the name of the amqp job in job service
	// AuditLogsGDPRCompliantVendorType : the name of the job which makes audit logs table GDPR-compliant
	// DiscordJobVendorType : the name of the discord job in job service
	// ExecSweepVendorType: the name of the execution sweep job
	// GarbageCollectionVendorType job name
	// ImageScanJobVendorType is name of scan job it will be used as key to register to job service.
	// P2PPreheatVendorType : the name of the P2P preheat job
	// PurgeAuditVendorType : the name of purge audit job
	// ReplicationVendorType : the name of the replication job in job service
	// RetentionVendorType : the name of the retention job
	// SBOMJobVendorType key to create sbom generate execution.
	// SampleJob is name of demo job
	// ScanAllVendorType: the name of the scan all job
	// ScanDataExportVendorType : the name of the scan data export job
	// SlackJobVendorType : the name of the slack job in job service
	// SystemArtifactCleanupVendorType : the name of the SystemArtifact cleanup job
	// WebhookJobVendorType : the name of the webhook job in job service
	// executionSweeperCount stores the count for execution retained
	AMQPJobVendorType = "AMQP"
	AuditLogsGDPRCompliantVendorType = "AUDIT_LOGS_GDPR_COMPLIANT"
	DiscordJobVendorType = "DISCORD"
	ExecSweepVendorType = "EXECUTION_SWEEP"
	GarbageCollectionVendorType = "GARBAGE_COLLECTION"
	ImageScanJobVendorType = "IMAGE_SCAN"
	P2PPreheatVendorType = "P2P_PREHEAT"
	PurgeAuditVendorType = "PURGE_AUDIT_LOG"
	ReplicationVendorType = "REPLICATION"
	RetentionVendorType = "RETENTION"
	SBOMJobVendorType = "SBOM"
	SampleJob = "DEMO"
	ScanAllVendorType = "SCAN_ALL"
	ScanDataExportVendorType = "SCAN_DATA_EXPORT"
	SlackJobVendorType = "SLACK"
	SystemArtifactCleanupVendorType = "SYSTEM_ARTIFACT_CLEANUP"
	WebhookJobVendorType = "WEBHOOK"
	executionSweeperCount = map[string]int64{
	return executionSweeperCount
	}
)
//
//    http://www.apache.org/licenses/LICENSE-2.0
// Copyright Project Harbor Authors
// Define the register name constants of known jobs
// GetExecutionSweeperCount gets the count of execution records retained by the sweeper
// Licensed under the Apache License, Version 2.0 (the "License");
// See the License for the specific language governing permissions and
// Unless required by applicable law or agreed to in writing, software
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// You may obtain a copy of the License at
// distributed under the License is distributed on an "AS IS" BASIS,
// limitations under the License.
// you may not use this file except in compliance with the License.
const (
func GetExecutionSweeperCount() map[string]int64 {
import "github.com/goharbor/harbor/src/lib"
package job
var (
}
