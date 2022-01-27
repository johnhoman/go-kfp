/*
Copyright 2022 John Homan.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fake_test

import (
	"context"
	"encoding/json"
	"fmt"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/johnhoman/go-kfp/api/job/client/job_service"
	jobmodels "github.com/johnhoman/go-kfp/api/job/models"
	"github.com/johnhoman/go-kfp/fake"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/google/uuid"

	ps "github.com/johnhoman/go-kfp/api/pipeline/client/pipeline_service"
	"github.com/johnhoman/go-kfp/api/pipeline/models"
	up "github.com/johnhoman/go-kfp/api/pipeline_upload/client/pipeline_upload_service"
	"github.com/johnhoman/go-kfp/pipelines"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

func newCowSay(name string) runtime.NamedReadCloser {
	// Not sure if the name actually matters -- might be able to swap it for a uuid
	reader := runtime.NamedReader(name+".yaml", strings.NewReader(`
    apiVersion: argoproj.io/v1alpha1
    kind: Workflow
    metadata:
      name: whalesay
    spec:
      entrypoint: whalesay
      arguments:
        parameters:
        - {name: name, value: Jack}
      templates:
      - name: whalesay
        inputs:
          parameters:
          - {name: name}
        container:
          image: docker/whalesay
          command: [cowsay]
          args: ["Hello, {{inputs.parameters.name}}!"]
`))
	return reader
}

func newJob(pipelineId string, versionId string) *jobmodels.APIJob {
	body := &jobmodels.APIJob{
		Name:        "whalesay-1m",
		Description: "Minutely whalesay",
		PipelineSpec: &jobmodels.APIPipelineSpec{
			Parameters: nil,
			PipelineID: pipelineId,
		},
		ResourceReferences: []*jobmodels.APIResourceReference{{
			Key: &jobmodels.APIResourceKey{
				ID:   versionId,
				Type: jobmodels.NewAPIResourceType(jobmodels.APIResourceTypePIPELINEVERSION),
			},
			Relationship: jobmodels.NewAPIRelationship(jobmodels.APIRelationshipCREATOR),
		}},
		Trigger: &jobmodels.APITrigger{
			CronSchedule: &jobmodels.APICronSchedule{
				Cron:      "* * * * *",
				StartTime: strfmt.DateTime(time.Now()),
				EndTime:   strfmt.DateTime(time.Now().Add(time.Second * 10)),
			},
		},
		NoCatchup:      true,
		MaxConcurrency: "1",
		Enabled:        true,
	}
	return body

}

type UploadService = up.ClientService
type Service = ps.ClientService
type JobService = job_service.ClientService

func pipelineService() pipelines.PipelineService {
	apiServer, ok := os.LookupEnv("GO_KFP_API_SERVER_ADDRESS")
	if ok {
		if strings.HasPrefix(apiServer, "http://") {
			apiServer = strings.TrimPrefix(apiServer, "http://")
		}
		transport := httptransport.New(apiServer, "", []string{"http"})
		return PipelineService{
			UploadService: up.New(transport, strfmt.Default),
			Service:       ps.New(transport, strfmt.Default),
			JobService:    job_service.New(transport, strfmt.Default),
		}
	}
	return fake.NewPipelineService()
}

type PipelineService struct {
	UploadService
	Service
	JobService
}

var _ = Describe("PipelineService", func() {
	var service pipelines.PipelineService
	var ctx context.Context
	var cancelFunc context.CancelFunc
	var pipelineIds []string
	var reader runtime.NamedReadCloser
	var name string
	var description string
	BeforeEach(func() {
		name = "testpipeline-" + uuid.New().String()[:8]
		description = strings.Join(strings.Split(name, "-"), " ")

		pipelineIds = make([]string, 0)

		service = pipelineService()
		reader = newCowSay(name)

		ctx, cancelFunc = context.WithCancel(context.Background())
	})
	AfterEach(func() {
		if pipelineIds != nil {
			for _, id := range pipelineIds {
				out, err := service.DeletePipeline(&ps.DeletePipelineParams{
					ID:      id,
					Context: ctx,
				}, nil)
				if err != nil {
					Expect(err.(*ps.DeletePipelineDefault).Code()).To(Or(Equal(http.StatusOK), Equal(http.StatusNotFound)))
					Expect(out).Should(BeNil())
				} else {
					Expect(out).ShouldNot(BeNil())
				}
			}
		}
		pipelineIds = nil
		cancelFunc()
	})
	It("should upload a pipeline", func() {
		out, err := service.UploadPipeline(&up.UploadPipelineParams{
			Description: stringPtr(description),
			Name:        stringPtr(name),
			Uploadfile:  reader,
			Context:     ctx,
		}, nil)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out).ShouldNot(BeNil())
		pipelineIds = append(pipelineIds, out.Payload.ID)
		// TODO: check the default version

		version, err := service.GetPipelineVersion(&ps.GetPipelineVersionParams{
			VersionID: out.GetPayload().ID,
			Context:   ctx,
		}, nil)
		Expect(version.GetPayload().ID).To(Equal(out.GetPayload().ID))
		Expect(version.GetPayload().Name).To(Equal(out.GetPayload().Name))
		Expect(version.GetPayload().Description).To(Equal(""))
		Expect(version.GetPayload().CreatedAt).To(Equal(out.GetPayload().CreatedAt))
		Expect(version.GetPayload().ResourceReferences[0].Key.ID).To(Equal(out.GetPayload().ID))
	})
	It("should return 404 for a pipeline that doesn't exist", func() {
		out, err := service.GetPipeline(&ps.GetPipelineParams{
			ID:      uuid.New().String(),
			Context: ctx,
		}, nil)
		Expect(err).Should(HaveOccurred())
		Expect(out).To(BeNil())
		Expect(err.(*ps.GetPipelineDefault).Code()).To(Equal(http.StatusNotFound))
	})
	It("Should return no default pipeline version when the default version doesn't exist", func() {
		out, err := service.UploadPipeline(&up.UploadPipelineParams{
			Description: stringPtr(description),
			Name:        stringPtr(name),
			Uploadfile:  reader,
			Context:     ctx,
		}, nil)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out).ShouldNot(BeNil())
		pipelineIds = append(pipelineIds, out.GetPayload().ID)

		_, err = service.DeletePipelineVersion(&ps.DeletePipelineVersionParams{
			VersionID: out.GetPayload().ID,
			Context:   ctx,
		}, nil)
		Expect(err).ToNot(HaveOccurred())

		get, err := service.GetPipeline(&ps.GetPipelineParams{ID: out.GetPayload().ID, Context: ctx}, nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(get.GetPayload().DefaultVersion).To(BeNil())
	})
	It("Should change the default version when it's deleted", func() {
		out, err := service.UploadPipeline(&up.UploadPipelineParams{
			Description: stringPtr(description),
			Name:        stringPtr(name),
			Uploadfile:  reader,
			Context:     ctx,
		}, nil)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out).ShouldNot(BeNil())
		pipelineIds = append(pipelineIds, out.GetPayload().ID)

		reader = newCowSay(name)
		version, err := service.UploadPipelineVersion(&up.UploadPipelineVersionParams{
			Description: stringPtr(description),
			Name:        stringPtr(name + "-1"),
			Pipelineid:  stringPtr(out.GetPayload().ID),
			Uploadfile:  reader,
			Context:     ctx,
		}, nil)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out).ShouldNot(BeNil())

		get, err := service.GetPipeline(&ps.GetPipelineParams{ID: out.GetPayload().ID, Context: ctx}, nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(get.GetPayload().DefaultVersion.ID).To(Equal(version.GetPayload().ID))

		_, err = service.DeletePipelineVersion(&ps.DeletePipelineVersionParams{
			VersionID: get.GetPayload().ID,
			Context:   ctx,
		}, nil)
		Expect(err).ToNot(HaveOccurred())

		get, err = service.GetPipeline(&ps.GetPipelineParams{ID: out.GetPayload().ID, Context: ctx}, nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(get.GetPayload().DefaultVersion).ToNot(BeNil())
		Expect(get.GetPayload().ID).To(Equal(out.GetPayload().ID))
	})
	Context("PipelineVersion", func() {
		It("should upload a pipeline version", func() {
			out, err := service.UploadPipeline(&up.UploadPipelineParams{
				Description: stringPtr(description),
				Name:        stringPtr(name),
				Uploadfile:  reader,
				Context:     ctx,
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).ShouldNot(BeNil())
			// Expect(out.Payload.Parameters[0].Name).To(Equal("name"))
			// Expect(out.Payload.Parameters[0].Value).To(Equal("Jack"))
			pipelineIds = append(pipelineIds, out.Payload.ID)

			reader = newCowSay(name)
			vsOut, err := service.UploadPipelineVersion(&up.UploadPipelineVersionParams{
				Description: stringPtr(description),
				Name:        stringPtr(name + "-v1"),
				Uploadfile:  reader,
				Context:     ctx,
				Pipelineid:  stringPtr(out.Payload.ID),
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(vsOut).ShouldNot(BeNil())
			Expect(vsOut.GetPayload().ResourceReferences[0].Key.ID).To(Equal(out.GetPayload().ID))
		})
		It("Should delete a pipeline version", func() {
			out, err := service.UploadPipeline(&up.UploadPipelineParams{
				Description: stringPtr(description),
				Name:        stringPtr(name),
				Uploadfile:  reader,
				Context:     ctx,
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).ShouldNot(BeNil())
			// Expect(out.Payload.Parameters[0].Name).To(Equal("name"))
			// Expect(out.Payload.Parameters[0].Value).To(Equal("Jack"))
			pipelineIds = append(pipelineIds, out.Payload.ID)

			reader = newCowSay(name)
			vsOut, err := service.UploadPipelineVersion(&up.UploadPipelineVersionParams{
				Description: stringPtr(description),
				Name:        stringPtr(name + "-v1"),
				Uploadfile:  reader,
				Context:     ctx,
				Pipelineid:  stringPtr(out.Payload.ID),
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(vsOut).ShouldNot(BeNil())
			Expect(vsOut.GetPayload().ID).ToNot(Equal(out.GetPayload().ID))

			delOut, err := service.DeletePipelineVersion(&ps.DeletePipelineVersionParams{
				VersionID: vsOut.GetPayload().ID,
				Context:   ctx,
			}, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(delOut).ToNot(BeNil())
			Expect(*delOut).To(Equal(ps.DeletePipelineVersionOK{Payload: map[string]interface{}{}}))
		})
		It("Cannot delete a pipeline version that doesn't exist", func() {
			delOut, err := service.DeletePipelineVersion(&ps.DeletePipelineVersionParams{
				VersionID: uuid.New().String(),
				Context:   ctx,
			}, nil)
			Expect(err).To(HaveOccurred())
			Expect(delOut).To(BeNil())
			out, ok := err.(*ps.DeletePipelineVersionDefault)
			Expect(ok).To(BeTrue())
			Expect(out.Code()).To(Equal(http.StatusNotFound))
		})
		It("Should get a pipeline version that does exist", func() {
			out, err := service.UploadPipeline(&up.UploadPipelineParams{
				Description: stringPtr(description),
				Name:        stringPtr(name),
				Uploadfile:  reader,
				Context:     ctx,
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).ShouldNot(BeNil())
			// Expect(out.Payload.Parameters[0].Name).To(Equal("name"))
			// Expect(out.Payload.Parameters[0].Value).To(Equal("Jack"))
			pipelineIds = append(pipelineIds, out.Payload.ID)

			reader = newCowSay(name)
			vsOut, err := service.UploadPipelineVersion(&up.UploadPipelineVersionParams{
				Description: stringPtr(description),
				Name:        stringPtr(name + "-v1"),
				Uploadfile:  reader,
				Context:     ctx,
				Pipelineid:  stringPtr(out.Payload.ID),
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(vsOut).ShouldNot(BeNil())
			Expect(vsOut.GetPayload().ID).ToNot(Equal(out.GetPayload().ID))

			getVersion, err := service.GetPipelineVersion(&ps.GetPipelineVersionParams{
				VersionID: vsOut.GetPayload().ID,
				Context:   ctx,
			}, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(getVersion).ToNot(BeNil())
			Expect(getVersion.GetPayload().Name).To(Equal(vsOut.GetPayload().Name))
		})
		It("Cannot get a pipeline version that does not exist", func() {
			out, err := service.UploadPipeline(&up.UploadPipelineParams{
				Description: stringPtr(description),
				Name:        stringPtr(name),
				Uploadfile:  reader,
				Context:     ctx,
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).ShouldNot(BeNil())
			pipelineIds = append(pipelineIds, out.Payload.ID)

			getVersion, err := service.GetPipelineVersion(&ps.GetPipelineVersionParams{
				VersionID: uuid.New().String(),
				Context:   ctx,
			}, nil)
			Expect(err).To(HaveOccurred())
			Expect(getVersion).To(BeNil())
			_, ok := err.(*ps.GetPipelineVersionDefault)
			Expect(ok).To(BeTrue())
			Expect(err.(*ps.GetPipelineVersionDefault).Code()).To(Equal(http.StatusNotFound))
		})
	})
	Context("UpdatePipelineDefaultVersion", func() {
		It("Should update the pipeline default version", func() {
			out, err := service.UploadPipeline(&up.UploadPipelineParams{
				Description: stringPtr(description),
				Name:        stringPtr(name),
				Uploadfile:  reader,
				Context:     ctx,
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).ShouldNot(BeNil())
			pipelineIds = append(pipelineIds, out.Payload.ID)

			reader = newCowSay(name)
			v1, err := service.UploadPipelineVersion(&up.UploadPipelineVersionParams{
				Description: stringPtr(description),
				Name:        stringPtr(name + "-v1"),
				Uploadfile:  reader,
				Context:     ctx,
				Pipelineid:  stringPtr(out.Payload.ID),
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(v1).ShouldNot(BeNil())
			Expect(v1.GetPayload().ID).ToNot(Equal(out.GetPayload().ID))

			reader = newCowSay(name)
			v2, err := service.UploadPipelineVersion(&up.UploadPipelineVersionParams{
				Description: stringPtr(description),
				Name:        stringPtr(name + "-v2"),
				Uploadfile:  reader,
				Context:     ctx,
				Pipelineid:  stringPtr(out.Payload.ID),
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(v2).ShouldNot(BeNil())
			Expect(v2.GetPayload().ID).ToNot(Equal(out.GetPayload().ID))

			getPipeline, err := service.GetPipeline(&ps.GetPipelineParams{
				ID:      out.GetPayload().ID,
				Context: ctx,
			}, nil)
			// This is not guaranteed - this has to do with default server behaviour
			Expect(getPipeline.GetPayload().DefaultVersion.ID).To(Equal(v2.GetPayload().ID))

			updateOut, err := service.UpdatePipelineDefaultVersion(&ps.UpdatePipelineDefaultVersionParams{
				PipelineID: out.GetPayload().ID,
				VersionID:  v1.GetPayload().ID,
				Context:    ctx,
			}, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(updateOut).ToNot(BeNil())
			Expect(updateOut).To(Equal(&ps.UpdatePipelineDefaultVersionOK{Payload: map[string]interface{}{}}))

			getPipeline, err = service.GetPipeline(&ps.GetPipelineParams{
				ID:      out.GetPayload().ID,
				Context: ctx,
			}, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(getPipeline).ToNot(BeNil())
			Expect(getPipeline.GetPayload().DefaultVersion.ID).To(Equal(v1.GetPayload().ID))
		})
		// Kubeflow bug
		It("Should fail silently for a version that does not exist", func() {
			out, err := service.UploadPipeline(&up.UploadPipelineParams{
				Description: stringPtr(description),
				Name:        stringPtr(name),
				Uploadfile:  reader,
				Context:     ctx,
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).ShouldNot(BeNil())
			pipelineIds = append(pipelineIds, out.Payload.ID)

			getOut, err := service.GetPipeline(&ps.GetPipelineParams{
				ID:      out.GetPayload().ID,
				Context: ctx,
			}, nil)
			// Didn't update
			Expect(getOut.GetPayload().DefaultVersion.ID).To(Equal(out.GetPayload().ID))

			id := uuid.New().String()
			updateOut, err := service.UpdatePipelineDefaultVersion(&ps.UpdatePipelineDefaultVersionParams{
				PipelineID: out.GetPayload().ID,
				VersionID:  id,
				Context:    ctx,
			}, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(updateOut).ToNot(BeNil())
			Expect(updateOut).To(Equal(&ps.UpdatePipelineDefaultVersionOK{Payload: map[string]interface{}{}}))

			getOut, err = service.GetPipeline(&ps.GetPipelineParams{
				ID:      out.GetPayload().ID,
				Context: ctx,
			}, nil)
			// Didn't update -- Kubeflow Bug
			Expect(getOut.GetPayload().DefaultVersion.ID).To(Equal(""))
		})
		// Kubeflow bug
		It("Silently fails to update a default version of a pipeline that does not exist", func() {
			out, err := service.UploadPipeline(&up.UploadPipelineParams{
				Description: stringPtr(description),
				Name:        stringPtr(name),
				Uploadfile:  reader,
				Context:     ctx,
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).ShouldNot(BeNil())
			pipelineIds = append(pipelineIds, out.Payload.ID)

			reader = newCowSay(name)
			v1, err := service.UploadPipelineVersion(&up.UploadPipelineVersionParams{
				Description: stringPtr(description),
				Name:        stringPtr(name + "-v1"),
				Uploadfile:  reader,
				Context:     ctx,
				Pipelineid:  stringPtr(out.Payload.ID),
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(v1).ShouldNot(BeNil())
			Expect(v1.GetPayload().ID).ToNot(Equal(out.GetPayload().ID))

			updateOut, err := service.UpdatePipelineDefaultVersion(&ps.UpdatePipelineDefaultVersionParams{
				PipelineID: uuid.New().String(),
				VersionID:  v1.GetPayload().ID,
				Context:    ctx,
			}, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(updateOut).ToNot(BeNil())

			getOut, err := service.GetPipeline(&ps.GetPipelineParams{
				ID:      out.GetPayload().ID,
				Context: ctx,
			}, nil)
			// Didn't update
			Expect(getOut.GetPayload().DefaultVersion.ID).To(Equal(v1.GetPayload().ID))
		})
	})
	Context("CreatePipeline", func() {
		It("Should return useless error when pipeline exists", func() {
			out, err := service.UploadPipeline(&up.UploadPipelineParams{
				Description: stringPtr(description),
				Name:        stringPtr(name),
				Uploadfile:  reader,
				Context:     ctx,
			}, nil)
			pipelineIds = append(pipelineIds, out.Payload.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).ShouldNot(BeNil())

			reader = newCowSay(name)
			out, err = service.UploadPipeline(&up.UploadPipelineParams{
				Description: stringPtr(description),
				Name:        stringPtr(name),
				Uploadfile:  reader,
				Context:     ctx,
			}, nil)
			Expect(err).Should(HaveOccurred())
			Expect(out).Should(BeNil())
			Expect(err.Error()).To(ContainSubstring("is not supported by the TextConsumer"))
		})
	})
	// Internal server errors are common errors returned by Kubeflow
	Context("DeletePipeline", func() {
		It("should delete a pipeline", func() {
			out, err := service.UploadPipeline(&up.UploadPipelineParams{
				Description: stringPtr(description),
				Name:        stringPtr(name),
				Uploadfile:  reader,
				Context:     ctx,
			}, nil)
			pipelineIds = append(pipelineIds, out.Payload.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).ShouldNot(BeNil())

			delOut, err := service.DeletePipeline(&ps.DeletePipelineParams{
				Context: ctx,
				ID:      out.Payload.ID,
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*delOut).To(Equal(ps.DeletePipelineOK{Payload: map[string]interface{}{}}))
		})
		It("Should return 404 when deleting a pipeline that doesn't exist", func() {
			_, err := service.DeletePipeline(&ps.DeletePipelineParams{
				Context: ctx,
				ID:      uuid.New().String(),
			}, nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.(*ps.DeletePipelineDefault).Code()).To(Equal(http.StatusNotFound))
		})

	})
	It("Should list pipelines", func() {
		Expect(reader.Close()).To(Succeed())
		for k := 0; k < 5; k++ {
			reader = newCowSay(name)
			out, err := service.UploadPipeline(&up.UploadPipelineParams{
				Description: stringPtr(description),
				Name:        stringPtr(fmt.Sprintf("%s-%d", name, k)),
				Uploadfile:  reader,
				Context:     ctx,
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).ShouldNot(BeNil())
			pipelineIds = append(pipelineIds, out.Payload.ID)
		}
		filter := map[string]interface{}{
			"predicates": []interface{}{
				map[string]interface{}{"op": "IS_SUBSTRING", "key": "name", "string_value": name},
			},
		}
		raw, err := json.Marshal(filter)
		Expect(err).To(Succeed())

		pipelines := make([]*models.APIPipeline, 0, 5)

		listOut, err := service.ListPipelines(&ps.ListPipelinesParams{
			Filter:                   stringPtr(string(raw)),
			PageSize:                 int32Ptr(3),
			PageToken:                nil,
			Context:                  ctx,
			ResourceReferenceKeyType: stringPtr(string(models.APIResourceTypeNAMESPACE)),
		}, nil)
		// eyJTb3J0QnlGaWVsZE5hbWUiOiJDcmVhdGVkQXRJblNlYyIsIlNvcnRCeUZpZWxkVmFsdWUiOjE2NDI0NzU2MzQsIlNvcnRCeUZpZWxkUHJlZml4IjoicGlwZWxpbmVzLiIsIktleUZpZWxkTmFtZSI6IlVVSUQiLCJLZXlGaWVsZFZhbHVlIjoiYjUxMDU0YjctYzc4OS00YTQ5LWJiNGQtODZkOTEwMzQwMDZhIiwiS2V5RmllbGRQcmVmaXgiOiJwaXBlbGluZXMuIiwiSXNEZXNjIjpmYWxzZSwiTW9kZWxOYW1lIjoicGlwZWxpbmVzIiwiRmlsdGVyIjp7IkZpbHRlclByb3RvIjoie1wicHJlZGljYXRlc1wiOlt7XCJvcFwiOlwiSVNfU1VCU1RSSU5HXCIsXCJrZXlcIjpcInBpcGVsaW5lcy5OYW1lXCIsXCJzdHJpbmdWYWx1ZVwiOlwidGVzdHBpcGVsaW5lLTA2MmYzY2FjXCJ9XX0iLCJFUSI6e30sIk5FUSI6e30sIkdUIjp7fSwiR1RFIjp7fSwiTFQiOnt9LCJMVEUiOnt9LCJJTiI6e30sIlNVQlNUUklORyI6eyJwaXBlbGluZXMuTmFtZSI6WyJ0ZXN0cGlwZWxpbmUtMDYyZjNjYWMiXX19fQ==
		Expect(err).ToNot(HaveOccurred())
		Expect(listOut.Payload.Pipelines).To(HaveLen(3))
		Expect(listOut.Payload.TotalSize).To(Equal(int32(5)))
		pipelines = append(pipelines, listOut.Payload.Pipelines...)
		listOut, err = service.ListPipelines(&ps.ListPipelinesParams{
			Filter:                   stringPtr(string(raw)),
			PageSize:                 int32Ptr(2),
			PageToken:                &listOut.Payload.NextPageToken,
			ResourceReferenceKeyType: stringPtr(string(models.APIResourceTypeNAMESPACE)),
			Context:                  ctx,
		}, nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(listOut.Payload.Pipelines).To(HaveLen(2))
		Expect(listOut.Payload.TotalSize).To(Equal(int32(5)))
		pipelines = append(pipelines, listOut.Payload.Pipelines...)

		for _, k := range []string{"-0", "-1", "-2", "-3", "-4"} {
			found := false
			for _, pl := range pipelines {
				if strings.HasSuffix(pl.Name, k) {
					found = true
				}
			}
			Expect(found).To(BeTrue())
		}
	})
	It("Should list pipeline versions", func() {
		out, err := service.UploadPipeline(&up.UploadPipelineParams{
			Description: stringPtr(description),
			Name:        stringPtr(name),
			Uploadfile:  reader,
			Context:     ctx,
		}, nil)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out).ShouldNot(BeNil())
		pipelineIds = append(pipelineIds, out.Payload.ID)
		for k := 0; k < 5; k++ {
			reader = newCowSay(name)
			out, err := service.UploadPipelineVersion(&up.UploadPipelineVersionParams{
				Description: stringPtr(description),
				Name:        stringPtr(fmt.Sprintf("%s-%d", name, k)),
				Pipelineid:  stringPtr(out.GetPayload().ID),
				Uploadfile:  reader,
				Context:     ctx,
			}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(out).ShouldNot(BeNil())
		}
		filter := map[string]interface{}{
			"predicates": []interface{}{
				map[string]interface{}{"op": "IS_SUBSTRING", "key": "name", "string_value": name},
			},
		}
		raw, err := json.Marshal(filter)
		Expect(err).To(Succeed())

		versions := make([]*models.APIPipelineVersion, 0, 5)

		listOut, err := service.ListPipelineVersions(&ps.ListPipelineVersionsParams{
			Filter:          stringPtr(string(raw)),
			PageSize:        int32Ptr(3),
			PageToken:       nil,
			Context:         ctx,
			ResourceKeyType: stringPtr(string(models.APIResourceTypePIPELINE)),
			ResourceKeyID:   stringPtr(out.GetPayload().ID),
		}, nil)
		// eyJTb3J0QnlGaWVsZE5hbWUiOiJDcmVhdGVkQXRJblNlYyIsIlNvcnRCeUZpZWxkVmFsdWUiOjE2NDI0NzU2MzQsIlNvcnRCeUZpZWxkUHJlZml4IjoicGlwZWxpbmVzLiIsIktleUZpZWxkTmFtZSI6IlVVSUQiLCJLZXlGaWVsZFZhbHVlIjoiYjUxMDU0YjctYzc4OS00YTQ5LWJiNGQtODZkOTEwMzQwMDZhIiwiS2V5RmllbGRQcmVmaXgiOiJwaXBlbGluZXMuIiwiSXNEZXNjIjpmYWxzZSwiTW9kZWxOYW1lIjoicGlwZWxpbmVzIiwiRmlsdGVyIjp7IkZpbHRlclByb3RvIjoie1wicHJlZGljYXRlc1wiOlt7XCJvcFwiOlwiSVNfU1VCU1RSSU5HXCIsXCJrZXlcIjpcInBpcGVsaW5lcy5OYW1lXCIsXCJzdHJpbmdWYWx1ZVwiOlwidGVzdHBpcGVsaW5lLTA2MmYzY2FjXCJ9XX0iLCJFUSI6e30sIk5FUSI6e30sIkdUIjp7fSwiR1RFIjp7fSwiTFQiOnt9LCJMVEUiOnt9LCJJTiI6e30sIlNVQlNUUklORyI6eyJwaXBlbGluZXMuTmFtZSI6WyJ0ZXN0cGlwZWxpbmUtMDYyZjNjYWMiXX19fQ==
		Expect(err).ToNot(HaveOccurred())
		Expect(listOut.Payload.Versions).To(HaveLen(3))
		Expect(listOut.Payload.TotalSize).To(Equal(int32(6)))
		versions = append(versions, listOut.Payload.Versions...)
		listOut, err = service.ListPipelineVersions(&ps.ListPipelineVersionsParams{
			Filter:          stringPtr(string(raw)),
			PageSize:        int32Ptr(3),
			PageToken:       &listOut.Payload.NextPageToken,
			ResourceKeyType: stringPtr(string(models.APIResourceTypePIPELINE)),
			ResourceKeyID:   stringPtr(out.GetPayload().ID),
			Context:         ctx,
		}, nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(listOut.Payload.Versions).To(HaveLen(3))
		Expect(listOut.Payload.TotalSize).To(Equal(int32(6)))
		versions = append(versions, listOut.Payload.Versions...)

		for _, k := range []string{"-0", "-1", "-2", "-3", "-4"} {
			found := false
			for _, pl := range versions {
				if strings.HasSuffix(pl.Name, k) {
					found = true
				}
			}
			Expect(found).To(BeTrue())
		}
	})
	Describe("JobServiceApi", func() {
		var jobId string
		var versionPipelineId string
		AfterEach(func() {
			out, err := service.DeleteJob(&job_service.DeleteJobParams{
				ID: jobId,
				Context: ctx,
			}, nil)
			if err != nil {
				Expect(err.(*job_service.DeleteJobDefault).Code()).Should(Equal(http.StatusNotFound))
				Expect(out).Should(BeNil())
			} else {
				Expect(out).Should(Equal(&job_service.DeleteJobOK{Payload: map[string]interface{}{}}))
			}
		})
		When("A pipeline exists", func() {
			BeforeEach(func() {
				out, err := service.UploadPipeline(&up.UploadPipelineParams{
					Description: stringPtr(description),
					Name:        stringPtr(name),
					Uploadfile:  reader,
					Context:     ctx,
				}, nil)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(out).ShouldNot(BeNil())
				pipelineIds = append(pipelineIds, out.Payload.ID)
				versionPipelineId = out.Payload.ID
			})
			When("A job does not exist", func() {
				It("Creates the job", func() {
					var buf []byte
					_, err := reader.Read(buf)
					Expect(err).ToNot(BeNil())
					out, err := service.CreateJob(&job_service.CreateJobParams{
						Context: ctx,
						Body:    newJob(versionPipelineId, versionPipelineId),
					}, nil)
					jobId = out.GetPayload().ID
					Expect(err).ToNot(HaveOccurred())
					Expect(out).ToNot(BeNil())
					Expect(out.GetPayload().PipelineSpec.PipelineName).To(Equal(name))
					Expect(out.GetPayload().PipelineSpec.PipelineID).To(Equal(versionPipelineId))
					Expect(out.GetPayload().PipelineSpec.PipelineManifest).To(Equal(string(buf)))
					Expect(out.GetPayload().Mode).To(BeNil())
					Expect(out.GetPayload().Enabled).To(BeTrue())
					Expect(out.GetPayload().ResourceReferences).To(HaveLen(2))
					Expect(out.GetPayload().NoCatchup).To(BeTrue())
				})
			})
			When("A job exists", func() {
				BeforeEach(func() {
					out, err := service.CreateJob(&job_service.CreateJobParams{
						Context: ctx,
						Body:    newJob(versionPipelineId, versionPipelineId),
					}, nil)
					jobId = out.GetPayload().ID
					Expect(err).ToNot(HaveOccurred())
					Expect(out).ToNot(BeNil())
				})
				It("Can list a single job", func() {

				})
			})
		})
		When("A pipeline does not exist", func() {
			BeforeEach(func() {
				out, err := service.UploadPipeline(&up.UploadPipelineParams{
					Description: stringPtr(description),
					Name:        stringPtr(name),
					Uploadfile:  reader,
					Context:     ctx,
				}, nil)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(out).ShouldNot(BeNil())
				pipelineIds = append(pipelineIds, out.Payload.ID)
				versionPipelineId = out.Payload.ID
				It("Should return 404", func() {
					out, err := service.CreateJob(&job_service.CreateJobParams{
						Context: ctx,
						Body:    newJob(uuid.New().String(), out.GetPayload().ID),
					}, nil)
					jobId = out.GetPayload().ID
					Expect(err).Should(HaveOccurred())
					Expect(err.(*job_service.CreateJobDefault).Code()).Should(Equal(http.StatusNotFound))
					Expect(out).To(BeNil())
				})
			})
		})
		When("A pipeline version does not exist", func() {
			BeforeEach(func() {
				out, err := service.UploadPipeline(&up.UploadPipelineParams{
					Description: stringPtr(description),
					Name:        stringPtr(name),
					Uploadfile:  reader,
					Context:     ctx,
				}, nil)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(out).ShouldNot(BeNil())
				pipelineIds = append(pipelineIds, out.Payload.ID)
				versionPipelineId = out.Payload.ID
			})
			It("Should return 404", func() {
				out, err := service.CreateJob(&job_service.CreateJobParams{
					Context: ctx,
					Body:    newJob(versionPipelineId, uuid.New().String()),
				}, nil)
				Expect(err).Should(HaveOccurred())
				Expect(err.(*job_service.CreateJobDefault).Code()).Should(Equal(http.StatusNotFound))
				Expect(out).To(BeNil())
			})
		})
	})
})