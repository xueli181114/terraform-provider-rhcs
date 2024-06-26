package exec

import (
	"context"
	"fmt"

	"github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/constants"
	"github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/helper"
)

type KubeletConfigArgs struct {
	Cluster             string `json:"cluster,omitempty"`
	PodPidsLimit        int    `json:"pod_pids_limit,omitempty"`
	NamePrefix          string `json:"name_prefix,omitempty"`
	KubeLetConfigNumber int    `json:"kubelet_config_number,omitempty"`
}

type KubeletConfig struct {
	Cluster      string `json:"cluster,omitempty"`
	PodPidsLimit int    `json:"pod_pids_limit,omitempty"`
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
}

type KubeletConfigService struct {
	CreationArgs *KubeletConfigArgs
	ManifestDir  string
	Context      context.Context
}

func (kc *KubeletConfigService) Init(manifestDirs ...string) error {
	kc.ManifestDir = constants.KubeletConfigDir
	if len(manifestDirs) != 0 {
		kc.ManifestDir = manifestDirs[0]
	}
	ctx := context.TODO()
	kc.Context = ctx
	err := runTerraformInit(ctx, kc.ManifestDir)
	if err != nil {
		return err
	}
	return nil

}

func (kc *KubeletConfigService) Apply(createArgs *KubeletConfigArgs, recordtfvars bool, extraArgs ...string) ([]*KubeletConfig, error) {
	kc.CreationArgs = createArgs
	args, tfvars := combineStructArgs(createArgs, extraArgs...)
	_, err := runTerraformApply(kc.Context, kc.ManifestDir, args...)
	if err != nil {
		return nil, err
	}
	if recordtfvars {
		recordTFvarsFile(kc.ManifestDir, tfvars)
	}
	output, err := kc.Output()
	return output, err
}
func (kc *KubeletConfigService) Plan(createArgs *KubeletConfigArgs, extraArgs ...string) (string, error) {
	kc.CreationArgs = createArgs
	args, _ := combineStructArgs(createArgs, extraArgs...)
	output, err := runTerraformPlan(kc.Context, kc.ManifestDir, args...)

	return output, err
}
func (kc *KubeletConfigService) Output() ([]*KubeletConfig, error) {
	out, err := runTerraformOutput(kc.Context, kc.ManifestDir)
	if err != nil {
		return nil, err
	}
	kubeletConfigsList := helper.DigArray(out["kubelet_configs"], "value")
	if kubeletConfigsList == nil {
		return nil, nil
	}
	kubeletConfigs := []*KubeletConfig{}
	for _, kubeletConfigsArray := range kubeletConfigsList {
		kubeletConfig := new(KubeletConfig)
		err = helper.MapStructure(kubeletConfigsArray.(map[string]interface{}), kubeletConfig)
		if err != nil {
			return kubeletConfigs, err
		}
		kubeletConfigs = append(kubeletConfigs, kubeletConfig)
	}
	return kubeletConfigs, nil
}

func (kc *KubeletConfigService) Destroy(createArgs ...*KubeletConfigArgs) (string, error) {
	if kc.CreationArgs == nil && len(createArgs) == 0 {
		return "", fmt.Errorf("got unset destroy args, set it in object or pass as a parameter")
	}
	destroyArgs := kc.CreationArgs
	if len(createArgs) != 0 {
		destroyArgs = createArgs[0]
	}
	args, _ := combineStructArgs(destroyArgs)
	output, err := runTerraformDestroy(kc.Context, kc.ManifestDir, args...)
	return output, err
}

func NewKubeletConfigService(manifestDir ...string) (*KubeletConfigService, error) {
	kc := &KubeletConfigService{}
	err := kc.Init(manifestDir...)
	return kc, err
}
