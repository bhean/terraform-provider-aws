// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ds

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/service/directoryservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func waitRegionCreated(ctx context.Context, conn *directoryservice.DirectoryService, directoryID, regionName string, timeout time.Duration) (*directoryservice.RegionDescription, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{directoryservice.DirectoryStageRequested, directoryservice.DirectoryStageCreating, directoryservice.DirectoryStageCreated},
		Target:  []string{directoryservice.DirectoryStageActive},
		Refresh: statusRegion(ctx, conn, directoryID, regionName),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*directoryservice.RegionDescription); ok {
		return output, err
	}

	return nil, err
}

func waitRegionDeleted(ctx context.Context, conn *directoryservice.DirectoryService, directoryID, regionName string, timeout time.Duration) (*directoryservice.RegionDescription, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{directoryservice.DirectoryStageActive, directoryservice.DirectoryStageDeleting},
		Target:  []string{},
		Refresh: statusRegion(ctx, conn, directoryID, regionName),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*directoryservice.RegionDescription); ok {
		return output, err
	}

	return nil, err
}

func waitSharedDirectoryDeleted(ctx context.Context, conn *directoryservice.DirectoryService, ownerDirectoryID, sharedDirectoryID string, timeout time.Duration) (*directoryservice.SharedDirectory, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			directoryservice.ShareStatusDeleting,
			directoryservice.ShareStatusShared,
			directoryservice.ShareStatusPendingAcceptance,
			directoryservice.ShareStatusRejectFailed,
			directoryservice.ShareStatusRejected,
			directoryservice.ShareStatusRejecting,
		},
		Target:                    []string{},
		Refresh:                   statusSharedDirectory(ctx, conn, ownerDirectoryID, sharedDirectoryID),
		Timeout:                   timeout,
		MinTimeout:                30 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*directoryservice.SharedDirectory); ok {
		return output, err
	}

	return nil, err
}

func waitDirectoryShared(ctx context.Context, conn *directoryservice.DirectoryService, id string, timeout time.Duration) (*directoryservice.SharedDirectory, error) {
	stateConf := &retry.StateChangeConf{
		Pending:                   []string{directoryservice.ShareStatusPendingAcceptance, directoryservice.ShareStatusSharing},
		Target:                    []string{directoryservice.ShareStatusShared},
		Refresh:                   statusDirectoryShareStatus(ctx, conn, id),
		Timeout:                   timeout,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*directoryservice.SharedDirectory); ok {
		return output, err
	}

	return nil, err
}
