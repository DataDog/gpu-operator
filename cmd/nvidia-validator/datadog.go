package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// MIG partition component
type MIGPartition struct {
	ctx        context.Context
	kubeClient kubernetes.Interface
}

const migPartitionStatusFile = "mig-ready"
const migPartitionNodeLabel = "nvidia.com/mig.partitioning"

func (m *MIGPartition) validate() error {
	// delete status file if already present
	err := deleteStatusFile(outputDirFlag + "/" + migPartitionStatusFile)
	if err != nil {
		return err
	}

	if err = m.configure(); err != nil {
		return err
	}

	shouldHaveMIGPartitioning, err := m.shouldHaveMIGPartitioning()
	if err != nil {
		return err
	}

	if shouldHaveMIGPartitioning {
		err = m.runValidation(false)
		if err != nil {
			log.Info("mig partioning is not ready")
			return err
		}
	}

	// create status file
	err = createStatusFile(outputDirFlag + "/" + migPartitionStatusFile)
	if err != nil {
		return err
	}
	return nil
}

func (m *MIGPartition) configure() error {
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Errorf("Error getting config cluster - %s\n", err.Error())
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Errorf("Error getting k8s client - %s\n", err.Error())
		return err
	}

	m.kubeClient = kubeClient
	return nil
}

func (m *MIGPartition) shouldHaveMIGPartitioning() (bool, error) {
	if m.ctx == nil {
		m.ctx = context.Background()
	}
	// get node info to check if MIG partioning should be present on node
	node, err := getNode(m.ctx, m.kubeClient)
	if err != nil {
		return false, fmt.Errorf("unable to fetch node by name %s to check for MIG Partitioning: %s", nodeNameFlag, err)
	}

	nodeLabels := node.GetLabels()
	partitioning, present := nodeLabels[migPartitionNodeLabel]
	return present && partitioning != "false", nil

}

func (m *MIGPartition) runValidation(silent bool) error {
	// check for MIG devices
	command := shell
	args := []string{"-c", "nvidia-smi -L | grep MIG"}

	if withWaitFlag {
		return runCommandWithWait(command, args, sleepIntervalSecondsFlag, silent)
	}
	return runCommand(command, args, silent)
}
