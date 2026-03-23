package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// EFADriver component
type EFADriver struct {
	ctx        context.Context
	kubeClient kubernetes.Interface
}

const efaDriverStatusFile = "efa-driver-ready"
const efaEnabledNodeLabel = "nodegroups.datadoghq.com/efa-enabled"

func (e *EFADriver) validate() error {
	if err := deleteStatusFile(outputDirFlag + "/" + efaDriverStatusFile); err != nil {
		return err
	}

	if err := e.configure(); err != nil {
		return err
	}

	isEFANode, err := e.isEFANode()
	if err != nil {
		return err
	}

	if isEFANode {
		if err := e.runValidation(false); err != nil {
			log.Info("EFA driver is not ready")
			return err
		}
	}

	return createStatusFile(outputDirFlag + "/" + efaDriverStatusFile)
}

func (e *EFADriver) configure() error {
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

	e.kubeClient = kubeClient
	return nil
}

func (e *EFADriver) isEFANode() (bool, error) {
	if e.ctx == nil {
		e.ctx = context.Background()
	}

	node, err := getNode(e.ctx, e.kubeClient)
	if err != nil {
		return false, fmt.Errorf("unable to fetch node %s to check for EFA: %w", nodeNameFlag, err)
	}

	return node.GetLabels()[efaEnabledNodeLabel] == "true", nil
}

func (e *EFADriver) runValidation(silent bool) error {
	command := shell
	args := []string{"-c", "ls /dev/efa/uverbs* > /dev/null 2>&1"}

	if withWaitFlag {
		return runCommandWithWait(command, args, sleepIntervalSecondsFlag, silent)
	}
	return runCommand(command, args, silent)
}

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
