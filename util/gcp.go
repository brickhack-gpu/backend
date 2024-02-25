package util

import (
	"context"
	"fmt"
    "log"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

func GetInstanceIP(projectID, zone, instanceName string) (string, error) {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return "", fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.GetInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	instance, err := instancesClient.Get(ctx, req)
	if err != nil {
		return "", fmt.Errorf("unable to get instance: %w", err)
	}
    log.Println(instance)

    return *instance.NetworkInterfaces[0].AccessConfigs[0].NatIP, nil
}

func DeleteInstance(projectID, zone, instanceName string) error {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.DeleteInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	op, err := instancesClient.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete instance: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Printf("Instance destroyed\n")

	return nil
}

func CreateInstance(projectID, zone, instanceName, machineType, sourceImage, region, script, gpuType string, gpuCount int32, disk int64) error {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.InsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Scheduling: &computepb.Scheduling{
				AutomaticRestart:  proto.Bool(true),
				OnHostMaintenance: proto.String("TERMINATE"),
				ProvisioningModel: proto.String("STANDARD"),
			},
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb: proto.Int64(disk),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
					Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
				},
			},
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					AccessConfigs: []*computepb.AccessConfig{
						{
							Name:        proto.String("External NAT"),
							NetworkTier: proto.String("PREMIUM"),
						},
					},
					StackType:  proto.String("IPV4_ONLY"),
					Subnetwork: proto.String(fmt.Sprintf("projects/siggpu/regions/%s/subnetworks/default", region)),
				},
			},
			Metadata: &computepb.Metadata{
				Items: []*computepb.Items{
					{
						Key:   proto.String("startup-script"),
						Value: proto.String(script),
					},
				},
			},
			GuestAccelerators: []*computepb.AcceleratorConfig{
				{
					AcceleratorCount: proto.Int32(gpuCount),
					AcceleratorType:  proto.String(fmt.Sprintf("projects/siggpu/zones/%s/acceleratorTypes/%s", zone, gpuType)),
				},
			},
			SourceMachineImage: proto.String(sourceImage),
		},
	}

	op, err := instancesClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create instance: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Printf("Instance created\n")

	return nil
}
