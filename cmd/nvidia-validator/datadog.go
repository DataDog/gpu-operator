package main

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v3"
	corev1 "k8s.io/api/core/v1"
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1validation "k8s.io/apimachinery/pkg/apis/meta/v1/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// validatorWorkloadLabelsEnvName indicates env name for extra labels (JSON object) applied to validation workload pods
const validatorWorkloadLabelsEnvName = "VALIDATOR_WORKLOAD_LABELS"

// validatorWorkloadAnnotationsEnvName indicates env name for extra annotations (JSON object) applied to validation workload pods
const validatorWorkloadAnnotationsEnvName = "VALIDATOR_WORKLOAD_ANNOTATIONS"

// extra labels and annotations applied to validation workload pods, parsed at startup
var (
	workloadLabelsFlag      map[string]string
	workloadAnnotationsFlag map[string]string
)

// datadogFlags returns the Datadog specific validator flags
func datadogFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "workload-labels",
			Usage: "extra labels (JSON object) applied to validation workload pods",
			Validator: func(value string) error {
				return parseWorkloadMetadata(value, &workloadLabelsFlag, func(m map[string]string) field.ErrorList {
					return metav1validation.ValidateLabels(m, field.NewPath("labels"))
				})
			},
			Sources: cli.EnvVars(validatorWorkloadLabelsEnvName),
		},
		&cli.StringFlag{
			Name:  "workload-annotations",
			Usage: "extra annotations (JSON object) applied to validation workload pods",
			Validator: func(value string) error {
				return parseWorkloadMetadata(value, &workloadAnnotationsFlag, func(m map[string]string) field.ErrorList {
					return apivalidation.ValidateAnnotations(m, field.NewPath("annotations"))
				})
			},
			Sources: cli.EnvVars(validatorWorkloadAnnotationsEnvName),
		},
	}
}

// parseWorkloadMetadata decodes the JSON object of string values into dest and validates it.
func parseWorkloadMetadata(value string, dest *map[string]string, validate func(map[string]string) field.ErrorList) error {
	if value == "" {
		return nil
	}

	var m map[string]string
	if err := json.Unmarshal([]byte(value), &m); err != nil {
		return fmt.Errorf("invalid JSON %q: %w", value, err)
	}
	if err := validate(m).ToAggregate(); err != nil {
		return err
	}
	*dest = m
	return nil
}

// applyWorkloadMetadata adds the extra labels and annotations to the given validation workload pod.
func applyWorkloadMetadata(pod *corev1.Pod) {
	if len(workloadLabelsFlag) > 0 {
		if pod.Labels == nil {
			pod.Labels = map[string]string{}
		}
		maps.Copy(pod.Labels, workloadLabelsFlag)
	}
	if len(workloadAnnotationsFlag) > 0 {
		if pod.Annotations == nil {
			pod.Annotations = map[string]string{}
		}
		maps.Copy(pod.Annotations, workloadAnnotationsFlag)
	}
}

// EFADriver component
type EFADriver struct {
	ctx        context.Context
	kubeClient kubernetes.Interface
}

const efaDriverStatusFile = "efa-driver-ready"
const efaEnabledNodeLabel = "nodegroups.datadoghq.com/efa-enabled"
const draNodeLabel = "managed-by-gpu-dra-plugin"

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

	// only verify MIG partitioning on Device Plugin-managed nodes
	isDRA, err := isDRANode(m.ctx, m.kubeClient)
	if err != nil {
		return false, err
	}
	if isDRA {
		log.Info("DRA node detected, skipping MIG partitioning validation")
		return false, nil
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

func isDRANode(ctx context.Context, kubeClient kubernetes.Interface) (bool, error) {
	node, err := getNode(ctx, kubeClient)
	if err != nil {
		return false, fmt.Errorf("unable to fetch node %s to check for DRA: %w", nodeNameFlag, err)
	}

	return node.GetLabels()[draNodeLabel] == "true", nil
}
