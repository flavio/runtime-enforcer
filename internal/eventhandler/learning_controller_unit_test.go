package eventhandler

import (
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestHandleAdmissionError(t *testing.T) {
	logger := logr.Discard()
	r := &LearningReconciler{}

	t.Run("returns nil for 422 Invalid admission error", func(t *testing.T) {
		invalidErr := apierrors.NewInvalid(
			schema.GroupKind{Group: "example.com", Kind: "Example"},
			"example-name",
			nil,
		)
		err := r.handleAdmissionError(logger, invalidErr)
		assert.NoError(t, err)
	})

	t.Run("returns nil for 410 Gone admission error", func(t *testing.T) {
		goneErr := apierrors.NewResourceExpired("resource is gone")
		err := r.handleAdmissionError(logger, goneErr)
		assert.NoError(t, err)
	})

	t.Run("returns original error for other APIStatus codes", func(t *testing.T) {
		badReqErr := apierrors.NewBadRequest("bad request")
		err := r.handleAdmissionError(logger, badReqErr)
		assert.ErrorIs(t, err, badReqErr, "expected returned error to wrap original bad request error")
	})

	t.Run("returns original error for non-APIStatus errors", func(t *testing.T) {
		plainErr := errors.New("plain error")
		err := r.handleAdmissionError(logger, plainErr)
		assert.ErrorIs(t, err, plainErr, "expected returned error to wrap original plain error")
	})
}
