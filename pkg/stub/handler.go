package stub

import (
	"context"
	"fmt"

	"github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1alpha1"
	"github.com/openshift/cluster-logging-operator/pkg/k8shandler"
	"github.com/openshift/cluster-logging-operator/pkg/utils"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {

	// Ignore the delete event since the garbage collector will clean up all secondary resources for the CR
	// All secondary resources must have the CR set as their OwnerReference for this to be the case
	if event.Deleted {
		return nil
	}

	switch o := event.Object.(type) {
	case *v1alpha1.ClusterLogging:
		return Reconcile(o)
	}

	return nil
}

func Reconcile(logging *v1alpha1.ClusterLogging) (err error) {
	exists := true

	// Reconcile certs
	if exists, logging = utils.DoesClusterLoggingExist(logging); exists {
		if err = k8shandler.CreateOrUpdateCertificates(logging); err != nil {
			return fmt.Errorf("Unable to create or update certificates: %v", err)
		}
	}

	// Reconcile Log Store
	if exists, logging = utils.DoesClusterLoggingExist(logging); exists {
		if err = k8shandler.CreateOrUpdateLogStore(logging); err != nil {
			return fmt.Errorf("Unable to create or update logstore: %v", err)
		}
	}

	// Reconcile Visualization
	if exists, logging = utils.DoesClusterLoggingExist(logging); exists {
		if err = k8shandler.CreateOrUpdateVisualization(logging); err != nil {
			return fmt.Errorf("Unable to create or update visualization: %v", err)
		}
	}

	// Reconcile Curation
	if exists, logging = utils.DoesClusterLoggingExist(logging); exists {
		if err = k8shandler.CreateOrUpdateCuration(logging); err != nil {
			return fmt.Errorf("Unable to create or update curation: %v", err)
		}
	}

	// Reconcile Collection
	if exists, logging = utils.DoesClusterLoggingExist(logging); exists {
		if err = k8shandler.CreateOrUpdateCollection(logging); err != nil {
			return fmt.Errorf("Unable to create or update collection: %v", err)
		}
	}

	return nil
}
