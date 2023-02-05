// Copyright 2022 Camunda Services GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web_modeler

import (
	corev1 "k8s.io/api/core/v1"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	appsv1 "k8s.io/api/apps/v1"
)

type restapiDeploymentTemplateTest struct {
	suite.Suite
	chartPath string
	release   string
	namespace string
	templates []string
}

func TestRestapiDeploymentTemplate(t *testing.T) {
	t.Parallel()

	chartPath, err := filepath.Abs("../../")
	require.NoError(t, err)

	suite.Run(t, &restapiDeploymentTemplateTest{
		chartPath: chartPath,
		release:   "camunda-platform-test",
		namespace: "camunda-platform-" + strings.ToLower(random.UniqueId()),
		templates: []string{"charts/web-modeler/templates/deployment-restapi.yaml"},
	})
}

func (s *restapiDeploymentTemplateTest) TestContainerShouldSetCorrectKeycloakServiceUrl() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"web-modeler.enabled":                   "true",
			"global.identity.keycloak.url.protocol": "http",
			"global.identity.keycloak.url.host":     "keycloak",
			"global.identity.keycloak.url.port":     "80",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	// then
	env := deployment.Spec.Template.Spec.Containers[0].Env
	s.Require().Contains(env,
		corev1.EnvVar{
			Name:  "RESTAPI_OAUTH2_TOKEN_ISSUER_BACKEND_URL",
			Value: "http://keycloak/auth/realms/camunda-platform",
		})
}

func (s *restapiDeploymentTemplateTest) TestContainerShouldSetCorrectKeycloakServiceUrlWithCustomPort() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"web-modeler.enabled":                   "true",
			"global.identity.keycloak.url.protocol": "http",
			"global.identity.keycloak.url.host":     "keycloak",
			"global.identity.keycloak.url.port":     "8888",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	// then
	env := deployment.Spec.Template.Spec.Containers[0].Env
	s.Require().Contains(env,
		corev1.EnvVar{
			Name:  "RESTAPI_OAUTH2_TOKEN_ISSUER_BACKEND_URL",
			Value: "http://keycloak:8888/auth/realms/camunda-platform",
		})
}

func (s *restapiDeploymentTemplateTest) TestContainerShouldSetCorrectIdentityServiceUrlWithFullnameOverride() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"web-modeler.enabled":              "true",
			"global.identity.fullnameOverride": "custom-identity-fullname",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	// then
	env := deployment.Spec.Template.Spec.Containers[0].Env
	s.Require().Contains(env,
		corev1.EnvVar{Name: "RESTAPI_IDENTITY_BASE_URL", Value: "http://custom-identity-fullname:80"})
}

func (s *restapiDeploymentTemplateTest) TestContainerShouldSetCorrectIdentityServiceUrlWithNameOverride() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"web-modeler.enabled":          "true",
			"global.identity.nameOverride": "custom-identity",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	// then
	env := deployment.Spec.Template.Spec.Containers[0].Env
	s.Require().Contains(env,
		corev1.EnvVar{Name: "RESTAPI_IDENTITY_BASE_URL", Value: "http://camunda-platform-test-custom-identity:80"})
}

func (s *restapiDeploymentTemplateTest) TestContainerShouldSetExternalDatabaseConfiguration() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"web-modeler.enabled":                           "true",
			"web-modeler.postgresql.enabled":                "false",
			"web-modeler.restapi.externalDatabase.host":     "postgres.example.com",
			"web-modeler.restapi.externalDatabase.port":     "65432",
			"web-modeler.restapi.externalDatabase.database": "modeler-database",
			"web-modeler.restapi.externalDatabase.user":     "modeler-user",
			"web-modeler.restapi.externalDatabase.password": "modeler-password",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	// then
	env := deployment.Spec.Template.Spec.Containers[0].Env
	s.Require().Contains(env, corev1.EnvVar{Name: "RESTAPI_DB_HOST", Value: "postgres.example.com"})
	s.Require().Contains(env, corev1.EnvVar{Name: "RESTAPI_DB_PORT", Value: "65432"})
	s.Require().Contains(env, corev1.EnvVar{Name: "RESTAPI_DB_NAME", Value: "modeler-database"})
	s.Require().Contains(env, corev1.EnvVar{Name: "RESTAPI_DB_USER", Value: "modeler-user"})
	s.Require().Contains(env,
		corev1.EnvVar{
			Name: "RESTAPI_DB_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: "camunda-platform-test-web-modeler-restapi"},
					Key:                  "database-password",
				},
			},
		})
}

func (s *restapiDeploymentTemplateTest) TestContainerShouldSetSmtpCredentials() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"web-modeler.enabled":                   "true",
			"web-modeler.restapi.mail.smtpUser":     "modeler-user",
			"web-modeler.restapi.mail.smtpPassword": "modeler-password",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	// then
	env := deployment.Spec.Template.Spec.Containers[0].Env
	s.Require().Contains(env, corev1.EnvVar{Name: "RESTAPI_MAIL_USER", Value: "modeler-user"})
	s.Require().Contains(env,
		corev1.EnvVar{
			Name: "RESTAPI_MAIL_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: "camunda-platform-test-web-modeler-restapi"},
					Key:                  "smtp-password",
				},
			},
		})
}

// TODO readinessProbe is disabled
// readinessProbe is enabled by default, so it's tested by golden files.

func (s *restapiDeploymentTemplateTest) TestContainerStartupProbe() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"web-modeler.enabled":                        "true",
			"web-modeler.restapi.startupProbe.enabled":   "true",
			"web-modeler.restapi.startupProbe.probePath": "/healthz",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	// then
	probe := deployment.Spec.Template.Spec.Containers[0].StartupProbe

	s.Require().Equal("/healthz", probe.HTTPGet.Path)
	s.Require().Equal("http-management", probe.HTTPGet.Port.StrVal)
}

func (s *restapiDeploymentTemplateTest) TestContainerLivenessProbe() {
	// given
	options := &helm.Options{
		SetValues: map[string]string{
			"web-modeler.enabled":                         "true",
			"web-modeler.restapi.livenessProbe.enabled":   "true",
			"web-modeler.restapi.livenessProbe.probePath": "/healthz",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", s.namespace),
	}

	// when
	output := helm.RenderTemplate(s.T(), options, s.chartPath, s.release, s.templates)
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(s.T(), output, &deployment)

	// then
	probe := deployment.Spec.Template.Spec.Containers[0].LivenessProbe

	s.Require().Equal("/healthz", probe.HTTPGet.Path)
	s.Require().Equal("http-management", probe.HTTPGet.Port.StrVal)
}
