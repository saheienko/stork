package storkctl

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	storkv1 "github.com/libopenstorage/stork/pkg/apis/stork/v1alpha1"
	storkops "github.com/portworx/sched-ops/k8s/stork"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/validation"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubernetes/pkg/printers"
)

const (
	clusterPairSubcommand = "clusterpair"
	cmdPathKey            = "cmd-path"
	gcloudPath            = "./google-cloud-sdk/bin/gcloud"
	gcloudBinaryName      = "gcloud"
)

var clusterPairColumns = []string{"NAME", "STORAGE-STATUS", "SCHEDULER-STATUS", "CREATED"}

func newGetClusterPairCommand(cmdFactory Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	getClusterPairCommand := &cobra.Command{
		Use:     clusterPairSubcommand,
		Aliases: []string{"cp"},
		Short:   "Get cluster pair resources",
		Run: func(c *cobra.Command, args []string) {
			var clusterPairs *storkv1.ClusterPairList

			namespaces, err := cmdFactory.GetAllNamespaces()
			if err != nil {
				util.CheckErr(err)
				return
			}

			if len(args) > 0 {
				clusterPairs = new(storkv1.ClusterPairList)
				for _, pairName := range args {
					for _, ns := range namespaces {
						pair, err := storkops.Instance().GetClusterPair(pairName, ns)
						if err != nil {
							util.CheckErr(err)
							return
						}
						clusterPairs.Items = append(clusterPairs.Items, *pair)
					}
				}
			} else {
				var tempClusterPairs storkv1.ClusterPairList
				for _, ns := range namespaces {
					clusterPairs, err = storkops.Instance().ListClusterPairs(ns)
					if err != nil {
						util.CheckErr(err)
						return
					}
					tempClusterPairs.Items = append(tempClusterPairs.Items, clusterPairs.Items...)
				}
				clusterPairs = &tempClusterPairs
			}

			if len(clusterPairs.Items) == 0 {
				handleEmptyList(ioStreams.Out)
				return
			}

			if err := printObjects(c, clusterPairs, cmdFactory, clusterPairColumns, clusterPairPrinter, ioStreams.Out); err != nil {
				util.CheckErr(err)
				return
			}
		},
	}

	cmdFactory.BindGetFlags(getClusterPairCommand.Flags())

	return getClusterPairCommand
}

func clusterPairPrinter(
	clusterPairList *storkv1.ClusterPairList,
	options printers.GenerateOptions,
) ([]metav1beta1.TableRow, error) {
	if clusterPairList == nil {
		return nil, nil
	}
	rows := make([]metav1beta1.TableRow, 0)
	for _, clusterPair := range clusterPairList.Items {
		creationTime := toTimeString(clusterPair.CreationTimestamp.Time)
		row := getRow(&clusterPair,
			[]interface{}{clusterPair.Name,
				clusterPair.Status.StorageStatus,
				clusterPair.Status.SchedulerStatus,
				creationTime},
		)
		rows = append(rows, row)
	}
	return rows, nil
}

func getStringData(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("error opening file %v: %v", fileName, err)
	}
	data, err := ioutil.ReadAll(bufio.NewReader(file))
	if err != nil {
		return "", fmt.Errorf("error reading file %v: %v", fileName, err)
	}

	return string(data), nil
}

func getByteData(fileName string) ([]byte, error) {
	data, err := getStringData(fileName)
	if err != nil {
		return nil, err
	}
	return []byte(data), nil
}

func newGenerateClusterPairCommand(cmdFactory Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	generateClusterPairCommand := &cobra.Command{
		Use:   clusterPairSubcommand,
		Short: "Generate a spec to be used for cluster pairing from a remote cluster",
		Run: func(c *cobra.Command, args []string) {
			if len(args) != 1 {
				util.CheckErr(fmt.Errorf("exactly one name needs to be provided for clusterpair name"))
				return
			}
			clusterPairName := args[0]
			config, err := cmdFactory.RawConfig()
			if err != nil {
				util.CheckErr(err)
				return
			}

			if errors := validation.NameIsDNSSubdomain(clusterPairName, false); len(errors) != 0 {
				err := fmt.Errorf("the Name \"%v\" is not valid: %v", clusterPairName, errors)
				util.CheckErr(err)
				return
			}
			if errors := validation.ValidateNamespaceName(cmdFactory.GetNamespace(), false); len(errors) != 0 {
				err := fmt.Errorf("the Namespace \"%v\" is not valid: %v", cmdFactory.GetNamespace(), errors)
				util.CheckErr(err)
				return
			}

			// Prune out all but the current-context and related
			// info
			currentContext := config.CurrentContext
			for context := range config.Contexts {
				if context != currentContext {
					delete(config.Contexts, context)
				}
			}
			if config.Contexts[currentContext] != nil {
				currentCluster := config.Contexts[currentContext].Cluster
				for cluster := range config.Clusters {
					if cluster != currentCluster {
						delete(config.Clusters, cluster)
					}
				}
				currentAuthInfo := config.Contexts[currentContext].AuthInfo
				for authInfo := range config.AuthInfos {
					if authInfo != currentAuthInfo {
						delete(config.AuthInfos, authInfo)
					}
				}

				if config.AuthInfos[currentAuthInfo] != nil {
					// Replace gcloud paths in the config
					if config.AuthInfos[currentAuthInfo].AuthProvider != nil &&
						config.AuthInfos[currentAuthInfo].AuthProvider.Config != nil {
						if cmdPath, present := config.AuthInfos[currentAuthInfo].AuthProvider.Config[cmdPathKey]; present {
							if strings.HasSuffix(cmdPath, gcloudBinaryName) {
								config.AuthInfos[currentAuthInfo].AuthProvider.Config[cmdPathKey] = gcloudPath
							}
						}
					}

					// Replace file paths with inline data
					if config.AuthInfos[currentAuthInfo].ClientCertificate != "" && len(config.AuthInfos[currentAuthInfo].ClientCertificateData) == 0 {
						config.AuthInfos[currentAuthInfo].ClientCertificateData, err = getByteData(config.AuthInfos[currentAuthInfo].ClientCertificate)
						if err != nil {
							util.CheckErr(err)
							return
						}
						config.AuthInfos[currentAuthInfo].ClientCertificate = ""
					}
					if config.AuthInfos[currentAuthInfo].ClientKey != "" && len(config.AuthInfos[currentAuthInfo].ClientKeyData) == 0 {
						config.AuthInfos[currentAuthInfo].ClientKeyData, err = getByteData(config.AuthInfos[currentAuthInfo].ClientKey)
						if err != nil {
							util.CheckErr(err)
							return
						}
						config.AuthInfos[currentAuthInfo].ClientKey = ""
					}
					if config.AuthInfos[currentAuthInfo].TokenFile != "" && len(config.AuthInfos[currentAuthInfo].Token) == 0 {
						config.AuthInfos[currentAuthInfo].Token, err = getStringData(config.AuthInfos[currentAuthInfo].TokenFile)
						if err != nil {
							util.CheckErr(err)
							return
						}
						config.AuthInfos[currentAuthInfo].TokenFile = ""
					}
				}
				if config.Clusters[currentCluster] != nil &&
					config.Clusters[currentCluster].CertificateAuthority != "" &&
					len(config.Clusters[currentCluster].CertificateAuthorityData) == 0 {

					config.Clusters[currentCluster].CertificateAuthorityData, err = getByteData(config.Clusters[currentCluster].CertificateAuthority)
					if err != nil {
						util.CheckErr(err)
						return
					}
					config.Clusters[currentCluster].CertificateAuthority = ""
				}

				clusterPair := &storkv1.ClusterPair{
					TypeMeta: meta.TypeMeta{
						Kind:       reflect.TypeOf(storkv1.ClusterPair{}).Name(),
						APIVersion: storkv1.SchemeGroupVersion.String(),
					},
					ObjectMeta: meta.ObjectMeta{
						Name:      clusterPairName,
						Namespace: cmdFactory.GetNamespace(),
					},

					Spec: storkv1.ClusterPairSpec{
						Config: config,
						Options: map[string]string{
							"<insert_storage_options_here>": "",
						},
					},
				}
				if err = printEncoded(c, clusterPair, "yaml", ioStreams.Out); err != nil {
					util.CheckErr(err)
					return
				}
			}
		},
	}

	return generateClusterPairCommand
}
