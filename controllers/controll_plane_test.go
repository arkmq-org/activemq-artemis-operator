/*
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
// +kubebuilder:docs-gen:collapse=Apache License

/*
As usual, we start with the necessary imports. We also define some utility variables.
*/
package controllers

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	brokerv1beta1 "github.com/arkmq-org/activemq-artemis-operator/api/v1beta1"
	"github.com/arkmq-org/activemq-artemis-operator/pkg/utils/common"
	rtclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// installOperatorCert installs the operator certificate
func installOperatorCert(namespace string) {
	By("installing operator cert")
	InstallCert(common.DefaultOperatorCertSecretName, namespace, func(candidate *cmv1.Certificate) {
		candidate.Spec.SecretName = common.DefaultOperatorCertSecretName
		candidate.Spec.CommonName = "activemq-artemis-operator"
		candidate.Spec.IssuerRef = cmmetav1.ObjectReference{
			Name: caIssuer.Name,
			Kind: "ClusterIssuer",
		}
	})
}

// installOperandCert installs the operand (broker) certificate and returns the cert name
func installOperandCert(crdName, namespace string) string {
	sharedOperandCertName := common.DefaultOperandCertSecretName
	By("installing restricted mtls broker cert")
	InstallCert(sharedOperandCertName, namespace, func(candidate *cmv1.Certificate) {
		candidate.Spec.SecretName = sharedOperandCertName
		candidate.Spec.CommonName = "activemq-artemis-operand"
		candidate.Spec.DNSNames = []string{common.OrdinalFQDNS(crdName, namespace, 0)}
		candidate.Spec.IssuerRef = cmmetav1.ObjectReference{
			Name: caIssuer.Name,
			Kind: "ClusterIssuer",
		}
	})
	return sharedOperandCertName
}

// createBaseCRD creates a basic ActiveMQArtemis CRD with minimal configuration
func createBaseCRD(namespace string) brokerv1beta1.ActiveMQArtemis {
	return brokerv1beta1.ActiveMQArtemis{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ActiveMQArtemis",
			APIVersion: brokerv1beta1.GroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      NextSpecResourceName(),
			Namespace: namespace,
		},
	}
}

// waitForBrokerReady waits for the broker to be in Ready and ConfigApplied state
func waitForBrokerReady(ctx context.Context, k8sClient rtclient.Client, brokerKey types.NamespacedName) {
	createdCrd := &brokerv1beta1.ActiveMQArtemis{}
	By("Checking ready, operator can access broker status via jmx")
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(ctx, brokerKey, createdCrd)).Should(Succeed())

		if verbose {
			fmt.Printf("STATUS: %v\n\n", createdCrd.Status.Conditions)
		}
		g.Expect(meta.IsStatusConditionTrue(createdCrd.Status.Conditions, brokerv1beta1.ReadyConditionType)).Should(BeTrue())
		g.Expect(meta.IsStatusConditionTrue(createdCrd.Status.Conditions, brokerv1beta1.ConfigAppliedConditionType)).Should(BeTrue())

	}, existingClusterTimeout, existingClusterInterval).Should(Succeed())
}

var _ = Describe("minimal", func() {

	var installedCertManager bool = false

	BeforeEach(func() {
		BeforeEachSpec()

		if verbose {
			fmt.Println("Time with MicroSeconds: ", time.Now().Format("2006-01-02 15:04:05.000000"), " test:", CurrentSpecReport())
		}

		if os.Getenv("USE_EXISTING_CLUSTER") == "true" {
			//if cert manager/trust manager is not installed, install it
			if !CertManagerInstalled() {
				Expect(InstallCertManager()).To(Succeed())
				installedCertManager = true
			}

			rootIssuer = InstallClusteredIssuer(rootIssuerName, nil)

			rootCert = InstallCert(rootCertName, rootCertNamespce, func(candidate *cmv1.Certificate) {
				candidate.Spec.IsCA = true
				candidate.Spec.CommonName = "artemis.root.ca"
				candidate.Spec.SecretName = rootCertSecretName
				candidate.Spec.IssuerRef = cmmetav1.ObjectReference{
					Name: rootIssuer.Name,
					Kind: "ClusterIssuer",
				}
			})

			caIssuer = InstallClusteredIssuer(caIssuerName, func(candidate *cmv1.ClusterIssuer) {
				candidate.Spec.SelfSigned = nil
				candidate.Spec.CA = &cmv1.CAIssuer{
					SecretName: rootCertSecretName,
				}
			})
			InstallCaBundle(common.DefaultOperatorCASecretName, rootCertSecretName, caPemTrustStoreName)

		}

	})

	AfterEach(func() {
		if false && os.Getenv("USE_EXISTING_CLUSTER") == "true" {
			UnInstallCaBundle(common.DefaultOperatorCASecretName)
			UninstallClusteredIssuer(caIssuerName)
			UninstallCert(rootCert.Name, rootCert.Namespace)
			UninstallClusteredIssuer(rootIssuerName)

			if installedCertManager {
				Expect(UninstallCertManager()).To(Succeed())
				installedCertManager = false
			}
		}

		AfterEachSpec()
	})

	Context("restricted rbac", func() {

		It("operator role access", func() {

			if os.Getenv("USE_EXISTING_CLUSTER") != "true" {
				return
			}

			installOperatorCert(defaultNamespace)

			ctx := context.Background()

			// empty CRD, name is used for cert subject to match the headless service
			crd := createBaseCRD(defaultNamespace)

			sharedOperandCertName := installOperandCert(crd.Name, defaultNamespace)

			prometheusCertName := common.DefaultPrometheusCertSecretName
			By("installing prometheus cert")
			InstallCert(prometheusCertName, defaultNamespace, func(candidate *cmv1.Certificate) {
				candidate.Spec.SecretName = prometheusCertName
				candidate.Spec.CommonName = "prometheus"
				candidate.Spec.IssuerRef = cmmetav1.ObjectReference{
					Name: caIssuer.Name,
					Kind: "ClusterIssuer",
				}
			})

			crd.Spec.Restricted = common.NewTrue()

			// how the jdk command line can be configured or modified
			crd.Spec.Env = []corev1.EnvVar{
				{Name: "JDK_JAVA_OPTIONS", Value: "-Djavax.net.debug=ssl -Djava.security.debug=logincontext"},
				//{Name: "JAVA_ARGS_APPEND", Value: "-DordinalProp=${STATEFUL_SET_ORDINAL}"},
			}
			crd.Spec.BrokerProperties = []string{
				"messageCounterSamplePeriod=500",
			}

			By("Deploying the CRD " + crd.ObjectMeta.Name)
			Expect(k8sClient.Create(ctx, &crd)).Should(Succeed())

			brokerKey := types.NamespacedName{Name: crd.Name, Namespace: crd.Namespace}
			createdCrd := &brokerv1beta1.ActiveMQArtemis{}

			waitForBrokerReady(ctx, k8sClient, brokerKey)

			serverName := common.OrdinalFQDNS(crd.Name, defaultNamespace, 0)
			By("setting up operator identity on http client")
			httpClient := http.Client{
				Transport: http.DefaultTransport,
				// A timeout less than 3 seconds may cause connection issues when
				// the server requires to change the chiper.
				Timeout: time.Second * 3,
			}

			httpClientTransport := httpClient.Transport.(*http.Transport)
			httpClientTransport.TLSClientConfig = &tls.Config{
				ServerName:         serverName,
				InsecureSkipVerify: false,
			}
			httpClientTransport.TLSClientConfig.GetClientCertificate =
				func(cri *tls.CertificateRequestInfo) (*tls.Certificate, error) {
					return common.GetOperatorClientCertificate(k8sClient, cri)
				}

			if rootCas, err := common.GetRootCAs(k8sClient); err == nil {
				httpClientTransport.TLSClientConfig.RootCAs = rootCas
			}

			By("Checking metrics with mtls are visible")
			Eventually(func(g Gomega) {

				resp, err := httpClient.Get("https://" + serverName + ":8888/metrics")
				if verbose {
					fmt.Printf("Resp form metrics get resp: %v, error: %v\n", resp, err)
				}
				g.Expect(err).Should(Succeed())

				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				g.Expect(err).Should(Succeed())

				lines := strings.Split(string(body), "\n")

				var done = false
				for _, line := range lines {
					if verbose {
						fmt.Printf("%s\n", line)
					}
					if strings.Contains(line, "artemis_total_pending_message_count") {
						done = true
					}
				}
				g.Expect(done).To(BeTrue())

			}, existingClusterTimeout, existingClusterInterval).Should(Succeed())

			Expect(k8sClient.Delete(ctx, createdCrd)).Should(Succeed())

			UninstallCert(common.DefaultOperatorCertSecretName, defaultNamespace)
			UninstallCert(sharedOperandCertName, defaultNamespace)
		})

		// Define the inputs you want to test
		sharedOverrides := []bool{true, false}

		for _, sharedOverride := range sharedOverrides {
			Context(fmt.Sprintf("when using shared override is %t", sharedOverride), func() {
				It("control plane override - custom metrics access", func() {

					if os.Getenv("USE_EXISTING_CLUSTER") != "true" {
						return
					}

					installOperatorCert(defaultNamespace)

					newUserCert := "new-user-cert"
					By("installing new-user cert (should be allowed via override)")
					InstallCert(newUserCert, defaultNamespace, func(candidate *cmv1.Certificate) {
						candidate.Spec.SecretName = newUserCert
						candidate.Spec.CommonName = "new-user"
						candidate.Spec.IssuerRef = cmmetav1.ObjectReference{
							Name: caIssuer.Name,
							Kind: "ClusterIssuer",
						}
					})
					newUserCertName := newUserCert

					randomCert := "random-cert"
					By("installing random cert (should be denied)")
					InstallCert(randomCert, defaultNamespace, func(candidate *cmv1.Certificate) {
						candidate.Spec.SecretName = randomCert
						candidate.Spec.CommonName = "random"
						candidate.Spec.IssuerRef = cmmetav1.ObjectReference{
							Name: caIssuer.Name,
							Kind: "ClusterIssuer",
						}
					})
					randomCertName := randomCert

					ctx := context.Background()

					crd := createBaseCRD(defaultNamespace)

					sharedOperandCertName := installOperandCert(crd.Name, defaultNamespace)

					overrideSecretName := "control-plane-override"
					if sharedOverride {
						overrideSecretName = crd.Name + "-control-plane-override"
					}
					By("creating control-plane-override secret with custom cert-users and cert-roles")
					overrideSecret := &corev1.Secret{
						TypeMeta: metav1.TypeMeta{
							Kind:       "Secret",
							APIVersion: "v1",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      overrideSecretName,
							Namespace: defaultNamespace,
						},
						StringData: map[string]string{
							"_cert-users": `# generated by test
#
					hawtio=/CN = hawtio-online\.hawtio\.svc.*/
					operator=/.*activemq-artemis-operator.*/
					probe=/.*activemq-artemis-operand.*/
					new-user=/.*new-user.*/
					`,
							"_cert-roles": `# generated by test
#
					status=operator,probe
					metrics=operator,new-user
					hawtio=hawtio
					`,
						},
					}
					Expect(k8sClient.Create(ctx, overrideSecret)).Should(Succeed())
					DeferCleanup(func() {
						k8sClient.Delete(ctx, overrideSecret)
					})

					crd.Spec.Restricted = common.NewTrue()
					crd.Spec.BrokerProperties = []string{
						"messageCounterSamplePeriod=500",
					}

					By("Deploying the CRD " + crd.ObjectMeta.Name)
					Expect(k8sClient.Create(ctx, &crd)).Should(Succeed())

					brokerKey := types.NamespacedName{Name: crd.Name, Namespace: crd.Namespace}
					createdCrd := &brokerv1beta1.ActiveMQArtemis{}

					waitForBrokerReady(ctx, k8sClient, brokerKey)

					serverName := common.OrdinalFQDNS(crd.Name, defaultNamespace, 0)

					By("Test 1: Operator cert can scrape metrics (default behavior)")
					Eventually(func(g Gomega) {
						transport := http.DefaultTransport.(*http.Transport).Clone()
						httpClient := http.Client{
							Transport: transport,
							Timeout:   time.Second * 3,
						}

						httpClientTransport := httpClient.Transport.(*http.Transport)
						httpClientTransport.TLSClientConfig = &tls.Config{
							ServerName:         serverName,
							InsecureSkipVerify: false,
						}
						httpClientTransport.TLSClientConfig.GetClientCertificate =
							func(cri *tls.CertificateRequestInfo) (*tls.Certificate, error) {
								return common.GetOperatorClientCertificate(k8sClient, cri)
							}

						if rootCas, err := common.GetRootCAs(k8sClient); err == nil {
							httpClientTransport.TLSClientConfig.RootCAs = rootCas
						}

						resp, err := httpClient.Get("https://" + serverName + ":8888/metrics")
						fmt.Printf("Test 1 - Operator cert: status=%d, err=%v\n", resp.StatusCode, err)
						g.Expect(err).Should(Succeed())
						g.Expect(resp.StatusCode).Should(Equal(200))

						defer resp.Body.Close()
						body, err := io.ReadAll(resp.Body)
						g.Expect(err).Should(Succeed())

						bodyStr := string(body)
						g.Expect(bodyStr).Should(ContainSubstring("artemis_total_pending_message_count"))

					}, existingClusterTimeout, existingClusterInterval).Should(Succeed())

					By("Test 2: New-user cert can scrape metrics (added via override)")
					Eventually(func(g Gomega) {
						transport := http.DefaultTransport.(*http.Transport).Clone()
						httpClient := http.Client{
							Transport: transport,
							Timeout:   time.Second * 3,
						}

						httpClientTransport := httpClient.Transport.(*http.Transport)
						httpClientTransport.TLSClientConfig = &tls.Config{
							ServerName:         serverName,
							InsecureSkipVerify: false,
						}
						httpClientTransport.TLSClientConfig.GetClientCertificate =
							func(cri *tls.CertificateRequestInfo) (*tls.Certificate, error) {
								secretKey := types.NamespacedName{Name: newUserCertName, Namespace: defaultNamespace}
								secret := &corev1.Secret{}
								if err := k8sClient.Get(ctx, secretKey, secret); err != nil {
									return nil, err
								}

								certPEM, ok := secret.Data["tls.crt"]
								if !ok {
									return nil, fmt.Errorf("tls.crt not found in secret %s", newUserCertName)
								}

								keyPEM, ok := secret.Data["tls.key"]
								if !ok {
									return nil, fmt.Errorf("tls.key not found in secret %s", newUserCertName)
								}

								cert, err := tls.X509KeyPair(certPEM, keyPEM)
								if err != nil {
									return nil, err
								}

								return &cert, nil
							}

						if rootCas, err := common.GetRootCAs(k8sClient); err == nil {
							httpClientTransport.TLSClientConfig.RootCAs = rootCas
						}

						resp, err := httpClient.Get("https://" + serverName + ":8888/metrics")
						fmt.Printf("Test 2 - New-user cert: status=%d, err=%v\n", resp.StatusCode, err)
						g.Expect(err).Should(Succeed())
						g.Expect(resp.StatusCode).Should(Equal(200))

						defer resp.Body.Close()
						body, err := io.ReadAll(resp.Body)
						g.Expect(err).Should(Succeed())

						bodyStr := string(body)
						g.Expect(bodyStr).Should(ContainSubstring("artemis_total_pending_message_count"))

					}, existingClusterTimeout, existingClusterInterval).Should(Succeed())

					By("Test 3: Random cert cannot scrape metrics (should be denied)")
					Eventually(func(g Gomega) {
						transport := http.DefaultTransport.(*http.Transport).Clone()
						httpClient := http.Client{
							Transport: transport,
							Timeout:   time.Second * 3,
						}

						httpClientTransport := httpClient.Transport.(*http.Transport)
						httpClientTransport.TLSClientConfig = &tls.Config{
							ServerName:         serverName,
							InsecureSkipVerify: false,
						}
						httpClientTransport.TLSClientConfig.GetClientCertificate =
							func(cri *tls.CertificateRequestInfo) (*tls.Certificate, error) {
								secretKey := types.NamespacedName{Name: randomCertName, Namespace: defaultNamespace}
								secret := &corev1.Secret{}
								if err := k8sClient.Get(ctx, secretKey, secret); err != nil {
									return nil, err
								}

								certPEM, ok := secret.Data["tls.crt"]
								if !ok {
									return nil, fmt.Errorf("tls.crt not found in secret %s", randomCertName)
								}

								keyPEM, ok := secret.Data["tls.key"]
								if !ok {
									return nil, fmt.Errorf("tls.key not found in secret %s", randomCertName)
								}

								cert, err := tls.X509KeyPair(certPEM, keyPEM)
								if err != nil {
									return nil, err
								}

								return &cert, nil
							}

						if rootCas, err := common.GetRootCAs(k8sClient); err == nil {
							httpClientTransport.TLSClientConfig.RootCAs = rootCas
						}

						resp, err := httpClient.Get("https://" + serverName + ":8888/metrics")
						fmt.Printf("Test 3 - Random cert: status=%d, err=%v\n", resp.StatusCode, err)
						g.Expect(err).Should(Succeed())
						// Should get 401 Unauthorized or 403 Forbidden
						g.Expect(resp.StatusCode).Should(Or(Equal(401), Equal(403)))

					}, existingClusterTimeout, existingClusterInterval).Should(Succeed())

					By("Cleaning up")
					Expect(k8sClient.Delete(ctx, createdCrd)).Should(Succeed())
					Expect(k8sClient.Delete(ctx, overrideSecret)).Should(Succeed())

					UninstallCert(common.DefaultOperatorCertSecretName, defaultNamespace)
					UninstallCert(sharedOperandCertName, defaultNamespace)
					UninstallCert(newUserCertName, defaultNamespace)
					UninstallCert(randomCertName, defaultNamespace)
				})
			})
		}
	})
})
