package controllers

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// efaEnabledLabelKey is the entry signal set by the Datadog nodegroup controller
	// on nodes that have EFA hardware. Equivalent to nvidia.com/gpu.present for GPU nodes.
	efaEnabledLabelKey   = "nodegroups.datadoghq.com/efa-enabled"
	efaEnabledLabelValue = "true"

	// efaDriverLabelKey is the deploy label set by the operator on EFA-enabled nodes.
	// The EFA driver DaemonSet uses this as its nodeSelector.
	efaDriverLabelKey   = "nvidia.com/gpu.deploy.efa-driver"
	efaDriverLabelValue = "true"
)

// DatadogLabelEFANodes labels EFA-enabled nodes with the operator's deploy label.
// It sets nvidia.com/gpu.deploy.efa-driver=true on nodes labeled with
// nodegroups.datadoghq.com/efa-enabled=true, and removes it when that label is gone.
func (n *ClusterPolicyController) labelEFANodes() error {
	list := &corev1.NodeList{}
	if err := n.client.List(n.ctx, list); err != nil {
		return fmt.Errorf("unable to list nodes for EFA labeling: %w", err)
	}

	for _, node := range list.Items {
		node := node
		original := node.DeepCopy()
		labels := node.GetLabels()

		hasEFAEnabled := labels[efaEnabledLabelKey] == efaEnabledLabelValue
		_, hasDeployLabel := labels[efaDriverLabelKey]

		modified := false
		if hasEFAEnabled && !hasDeployLabel {
			n.logger.Info("Setting EFA driver deploy label", "NodeName", node.Name,
				"Label", efaDriverLabelKey, "Value", efaDriverLabelValue)
			labels[efaDriverLabelKey] = efaDriverLabelValue
			modified = true
		} else if !hasEFAEnabled && hasDeployLabel {
			n.logger.Info("Removing EFA driver deploy label", "NodeName", node.Name,
				"Label", efaDriverLabelKey)
			delete(labels, efaDriverLabelKey)
			modified = true
		}

		if modified {
			node.SetLabels(labels)
			if err := n.client.Patch(n.ctx, &node, client.MergeFrom(original)); err != nil {
				return fmt.Errorf("unable to patch EFA labels on node %s: %w", node.Name, err)
			}
		}
	}
	return nil
}
